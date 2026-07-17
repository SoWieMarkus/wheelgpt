package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	helixBase = "https://api.twitch.tv/helix"
	tokenURL  = "https://id.twitch.tv/oauth2/token"
)

type Client struct {
	clientID     string
	clientSecret string
	redirectURL  string
	http         *http.Client

	mu        sync.Mutex
	appToken  string
	expiresAt time.Time
}

func NewClient(clientID, clientSecret, redirectURL string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		http:         &http.Client{Timeout: 10 * time.Second},
	}
}

// --- App token (client credentials) ---

func (c *Client) appAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.appToken != "" && time.Now().Before(c.expiresAt.Add(-30*time.Second)) {
		return c.appToken, nil
	}
	return c.refreshAppToken(ctx)
}

func (c *Client) refreshAppToken(ctx context.Context) (string, error) {
	vals := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"grant_type":    {"client_credentials"},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	c.appToken = out.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(out.ExpiresIn) * time.Second)
	return c.appToken, nil
}

// --- OAuth code exchange (user login) ---

type TwitchUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	ProfileImageURL string `json:"profile_image_url"`
}

func (c *Client) ExchangeCode(ctx context.Context, code string) (string, error) {
	vals := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {c.redirectURL},
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.AccessToken, nil
}

func (c *Client) GetUser(ctx context.Context, userToken string) (*TwitchUser, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, helixBase+"/users", nil)
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+userToken)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Data []TwitchUser `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Data) == 0 {
		return nil, fmt.Errorf("no user returned")
	}
	return &out.Data[0], nil
}

// --- Helix helpers (app token) ---

type StreamInfo struct {
	UserID string `json:"user_id"`
	Type   string `json:"type"` // "live" or ""
}

func (c *Client) GetStream(ctx context.Context, userID string) (*StreamInfo, error) {
	token, err := c.appAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		helixBase+"/streams?user_id="+userID, nil)
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Data []StreamInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Data) == 0 {
		return &StreamInfo{UserID: userID}, nil
	}
	return &out.Data[0], nil
}

func (c *Client) GetUserByID(ctx context.Context, userID string) (*TwitchUser, error) {
	token, err := c.appAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		helixBase+"/users?id="+userID, nil)
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Data []TwitchUser `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Data) == 0 {
		return nil, fmt.Errorf("user not found: %s", userID)
	}
	return &out.Data[0], nil
}
