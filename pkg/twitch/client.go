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

// Client is a Twitch API client that handles authentication, user information retrieval, and EventSub subscription management.
type Client struct {
	clientID            string
	clientSecret        string
	redirectURL         string
	http                *http.Client
	eventSubCallbackURL string
	eventSubSecret      string

	mu        sync.Mutex
	appToken  string
	expiresAt time.Time
}

// ClientParams holds the parameters required to create a new Twitch API client.
type ClientParams struct {
	ClientID            string
	ClientSecret        string
	RedirectURL         string
	EventSubCallbackURL string
	EventSubSecret      string
}

// NewClient creates a new Twitch API client with the provided credentials and redirect URL.
func NewClient(params ClientParams) *Client {
	return &Client{
		clientID:            params.ClientID,
		clientSecret:        params.ClientSecret,
		redirectURL:         params.RedirectURL,
		http:                &http.Client{},
		eventSubCallbackURL: params.EventSubCallbackURL,
		eventSubSecret:      params.EventSubSecret,
	}
}

// getAppAccessToken retrieves the app access token, refreshing it if necessary.
func (c *Client) getAppAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Safeguard against using a token that is about to expire within the next 30 seconds
	minDelta := 30 * time.Second

	tokenExpiresAt := c.expiresAt.Add(-minDelta)
	now := time.Now()

	// Check if the cached token is expired or will expire within the next 30 seconds
	if c.appToken != "" && now.Before(tokenExpiresAt) {
		return c.appToken, nil
	}
	return c.refreshAppAccessToken(ctx)
}

// refreshAppAccessToken requests a new app access token from Twitch and updates the cached token and expiration time.
func (c *Client) refreshAppAccessToken(ctx context.Context) (string, error) {
	values := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"grant_type":    {"client_credentials"},
	}

	body := strings.NewReader(values.Encode())

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://id.twitch.tv/oauth2/token", body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.http.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}

	tokenExpiresIn := time.Duration(result.ExpiresIn) * time.Second
	tokenExpiresAt := time.Now().Add(tokenExpiresIn)

	c.appToken = result.AccessToken
	c.expiresAt = tokenExpiresAt

	return c.appToken, nil
}

// ExchangeCode exchanges an authorization code from the Twitch OAuth2 flow for an access token using the Twitch OAuth2 API.
func (c *Client) ExchangeCode(ctx context.Context, code string) (string, error) {
	values := url.Values{
		"client_id":     {c.clientID},
		"client_secret": {c.clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {c.redirectURL},
	}

	body := strings.NewReader(values.Encode())

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://id.twitch.tv/oauth2/token", body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.http.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// GetUser retrieves user information from the Twitch API using the provided user access token.
func (c *Client) GetUser(ctx context.Context, userAccessToken string) (*User, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		return nil, err
	}

	bearerToken := fmt.Sprintf("Bearer %s", userAccessToken)
	request.Header.Set("Authorization", bearerToken)
	request.Header.Set("Client-Id", c.clientID)

	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}

	// Even though we are querying for one user, the Twitch API returns an array of users in the "data" field. We will decode this array and return the first user.
	var result struct {
		Data []User `json:"data"`
	}

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no user data found")
	}
	return &result.Data[0], nil
}

func (c *Client) GetUsers(ctx context.Context, channelIds []string) ([]User, error) {
	users := []User{}
	if len(channelIds) == 0 {
		return users, nil
	}

	// Split channelIds into chunks of 100, as the Twitch API has a limit of 100 IDs per request
	chunkSize := 100
	chunks := make([][]string, 0, (len(channelIds)+chunkSize-1)/chunkSize)
	for i := 0; i < len(channelIds); i += chunkSize {
		end := i + chunkSize
		if end > len(channelIds) {
			end = len(channelIds)
		}
		chunks = append(chunks, channelIds[i:end])
	}

	for _, chunk := range chunks {
		// Build the request URL with the chunk of channel IDs
		params := url.Values{
			"id": chunk,
		}
		requestURL := "https://api.twitch.tv/helix/users?" + params.Encode()

		request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, err
		}
		request.Header.Set("Client-Id", c.clientID)

		response, err := c.http.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		var result struct {
			Data []User `json:"data"`
		}
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			return nil, err
		}
		users = append(users, result.Data...)
	}

	return users, nil
}
