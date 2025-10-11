package identity

import (
	"fmt"
	"net/url"
	"time"
)

// Access token from the Twitch OAuth authorization code flow.
type UserAccessToken struct {
	AccessToken  string   `json:"access_token"`
	ExpiresIn    int      `json:"expires_in"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	Scope        []string `json:"scope"`
}

// Request a user access token using the authorization code flow.
func (c *Client) RequestUserAccessToken(code, redirectURI string) (*UserAccessToken, error) {
	data := url.Values{}
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", redirectURI)

	var token UserAccessToken
	if err := c.post("/token", data, &token); err != nil {
		return nil, fmt.Errorf("failed to request user access token: %w", err)
	}
	return &token, nil
}

// Access token from the Twitch OAuth client credentials flow.
type AppAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (a *AppAccessToken) IsExpired(tokenIssuedAt *time.Time) bool {
	margin := 300 // 5 minutes margin
	duration := time.Duration(a.ExpiresIn-margin) * time.Second
	return tokenIssuedAt.Add(duration).Before(time.Now())
}

// Request an app access token using client credentials flow.
func (c *Client) RequestAppAccessToken() (*AppAccessToken, error) {
	data := url.Values{}
	data.Set("client_id", c.config.ClientID)
	data.Set("client_secret", c.config.ClientSecret)
	data.Set("grant_type", "client_credentials")

	var token AppAccessToken
	if err := c.post("/token", data, &token); err != nil {
		return nil, fmt.Errorf("failed to request app access token: %w", err)
	}

	return &token, nil
}
