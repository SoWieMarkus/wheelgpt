// Package twitch is a small client for the Twitch Helix API, OAuth and
// EventSub subscriptions. Everything the bot does on Twitch goes through it.
package twitch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Default endpoints.
const (
	DefaultHelixURL = "https://api.twitch.tv/helix"
	DefaultAuthURL  = "https://id.twitch.tv/oauth2"
)

// Client talks to Twitch with an app access token.
type Client struct {
	ClientID     string
	ClientSecret string
	HelixURL     string
	AuthURL      string
	HTTP         *http.Client
	// OnError is called for failed API calls (metrics), may be nil.
	OnError func(endpoint string)

	mu          sync.Mutex
	appToken    string
	appTokenExp time.Time
}

// New returns a client with default endpoints.
func New(clientID, clientSecret string) *Client {
	return &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HelixURL:     DefaultHelixURL,
		AuthURL:      DefaultAuthURL,
		HTTP:         &http.Client{Timeout: 10 * time.Second},
	}
}

// APIError is a non-2xx answer from Twitch.
type APIError struct {
	Endpoint string
	Status   int
	Message  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("twitch %s: %d %s", e.Endpoint, e.Status, e.Message)
}

// StatusOf returns the HTTP status of a Twitch API error, or 0.
func StatusOf(err error) int {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Status
	}
	return 0
}

// HasAppToken reports whether an app token was fetched successfully.
func (c *Client) HasAppToken() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.appToken != "" && time.Now().Before(c.appTokenExp)
}

func (c *Client) token(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.appToken != "" && time.Now().Before(c.appTokenExp) {
		return c.appToken, nil
	}
	form := url.Values{
		"client_id":     {c.ClientID},
		"client_secret": {c.ClientSecret},
		"grant_type":    {"client_credentials"},
	}
	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := c.postForm(ctx, "app_token", c.AuthURL+"/token", form, &res); err != nil {
		return "", err
	}
	c.appToken = res.AccessToken
	// Refresh a little early.
	c.appTokenExp = time.Now().Add(time.Duration(res.ExpiresIn)*time.Second - 5*time.Minute)
	return c.appToken, nil
}

func (c *Client) invalidateToken() {
	c.mu.Lock()
	c.appToken = ""
	c.mu.Unlock()
}

func (c *Client) postForm(ctx context.Context, endpoint, u string, form url.Values, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.do(req, endpoint, out)
}

// helix calls a Helix endpoint with the app token (or userToken if set).
func (c *Client) helix(ctx context.Context, method, path string, query url.Values, body, out any, userToken string) error {
	endpoint := strings.TrimPrefix(path, "/")
	for attempt := 0; ; attempt++ {
		token := userToken
		if token == "" {
			var err error
			if token, err = c.token(ctx); err != nil {
				return err
			}
		}
		u := c.HelixURL + path
		if len(query) > 0 {
			u += "?" + query.Encode()
		}
		var reader io.Reader
		if body != nil {
			b, err := json.Marshal(body)
			if err != nil {
				return err
			}
			reader = bytes.NewReader(b)
		}
		req, err := http.NewRequestWithContext(ctx, method, u, reader)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Client-Id", c.ClientID)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		err = c.do(req, endpoint, out)
		// An app token can be revoked early; fetch a new one once.
		if StatusOf(err) == http.StatusUnauthorized && userToken == "" && attempt == 0 {
			c.invalidateToken()
			continue
		}
		return err
	}
}

func (c *Client) do(req *http.Request, endpoint string, out any) error {
	res, err := c.HTTP.Do(req)
	if err != nil {
		c.failed(endpoint)
		return fmt.Errorf("twitch %s: %w", endpoint, err)
	}
	defer func() { _ = res.Body.Close() }()
	data, err := io.ReadAll(io.LimitReader(res.Body, 4<<20))
	if err != nil {
		c.failed(endpoint)
		return fmt.Errorf("twitch %s: read body: %w", endpoint, err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		c.failed(endpoint)
		var e struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &e)
		return &APIError{Endpoint: endpoint, Status: res.StatusCode, Message: e.Message}
	}
	if out == nil || len(data) == 0 {
		return nil
	}
	if err := json.Unmarshal(data, out); err != nil {
		c.failed(endpoint)
		return fmt.Errorf("twitch %s: decode: %w", endpoint, err)
	}
	return nil
}

func (c *Client) failed(endpoint string) {
	if c.OnError != nil {
		c.OnError(endpoint)
	}
}

// User is a Twitch user.
type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

// Users returns the users with the given ids (any number, batched by 100).
func (c *Client) Users(ctx context.Context, ids []string) ([]User, error) {
	var all []User
	for batch := range chunks(ids, 100) {
		q := url.Values{"id": batch}
		var res struct{ Data []User }
		if err := c.helix(ctx, http.MethodGet, "/users", q, nil, &res, ""); err != nil {
			return all, err
		}
		all = append(all, res.Data...)
	}
	return all, nil
}

// UserByLogin returns a user by login.
func (c *Client) UserByLogin(ctx context.Context, login string) (User, error) {
	var res struct{ Data []User }
	if err := c.helix(ctx, http.MethodGet, "/users", url.Values{"login": {login}}, nil, &res, ""); err != nil {
		return User{}, err
	}
	if len(res.Data) == 0 {
		return User{}, fmt.Errorf("twitch user %q not found", login)
	}
	return res.Data[0], nil
}

// Me returns the user a user access token belongs to.
func (c *Client) Me(ctx context.Context, userToken string) (User, error) {
	var res struct{ Data []User }
	if err := c.helix(ctx, http.MethodGet, "/users", nil, nil, &res, userToken); err != nil {
		return User{}, err
	}
	if len(res.Data) == 0 {
		return User{}, errors.New("twitch returned no user for the token")
	}
	return res.Data[0], nil
}

// LiveUserIDs returns the ids of the given users that are live.
func (c *Client) LiveUserIDs(ctx context.Context, ids []string) (map[string]bool, error) {
	live := map[string]bool{}
	for batch := range chunks(ids, 100) {
		q := url.Values{"user_id": batch, "first": {"100"}}
		var res struct {
			Data []struct {
				UserID string `json:"user_id"`
				Type   string `json:"type"`
			}
		}
		if err := c.helix(ctx, http.MethodGet, "/streams", q, nil, &res, ""); err != nil {
			return live, err
		}
		for _, s := range res.Data {
			if s.Type == "live" {
				live[s.UserID] = true
			}
		}
	}
	return live, nil
}

// SendChatMessage posts a message as senderID into the broadcaster's chat.
// With an app token this needs user:bot from the sender and channel:bot from
// the broadcaster (or the sender being a moderator).
func (c *Client) SendChatMessage(ctx context.Context, broadcasterID, senderID, message string) error {
	body := map[string]string{"broadcaster_id": broadcasterID, "sender_id": senderID, "message": message}
	var res struct {
		Data []struct {
			IsSent     bool `json:"is_sent"`
			DropReason *struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"drop_reason"`
		}
	}
	if err := c.helix(ctx, http.MethodPost, "/chat/messages", nil, body, &res, ""); err != nil {
		return err
	}
	if len(res.Data) > 0 && !res.Data[0].IsSent {
		reason := "unknown"
		if d := res.Data[0].DropReason; d != nil {
			reason = d.Code + ": " + d.Message
		}
		return fmt.Errorf("twitch dropped chat message: %s", reason)
	}
	return nil
}

func chunks(s []string, n int) func(yield func([]string) bool) {
	return func(yield func([]string) bool) {
		for i := 0; i < len(s); i += n {
			if !yield(s[i:min(i+n, len(s))]) {
				return
			}
		}
	}
}
