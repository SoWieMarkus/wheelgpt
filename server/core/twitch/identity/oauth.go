package identity

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/SoWieMarkus/wheelgpt/core/http"
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
	body := url.Values{}
	body.Set("client_id", c.config.ClientID)
	body.Set("client_secret", c.config.ClientSecret)
	body.Set("code", code)
	body.Set("grant_type", "authorization_code")
	body.Set("redirect_uri", redirectURI)
	bodyReader := strings.NewReader(body.Encode())

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	request := http.HttpRequest{
		Endpoint: "/token",
		Headers:  &headers,
	}

	var token UserAccessToken
	_, err := c.client.Post(&request, bodyReader, &token)
	if err != nil {
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

// Check if the app access token is expired based on the provided timestamp, considering a margin for safety.
func (a *AppAccessToken) IsExpired(tokenIssuedAt *time.Time) bool {
	margin := 300 // 5 minutes margin
	duration := time.Duration(a.ExpiresIn-margin) * time.Second
	return tokenIssuedAt.Add(duration).Before(time.Now())
}

// Request an app access token using client credentials flow.
func (c *Client) RequestAppAccessToken() (*AppAccessToken, error) {
	body := url.Values{}
	body.Set("client_id", c.config.ClientID)
	body.Set("client_secret", c.config.ClientSecret)
	body.Set("grant_type", "client_credentials")
	bodyReader := strings.NewReader(body.Encode())

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	request := http.HttpRequest{
		Endpoint: "/token",
		Headers:  &headers,
	}

	var token AppAccessToken
	_, err := c.client.Post(&request, bodyReader, &token)
	if err != nil {
		return nil, fmt.Errorf("failed to request app access token: %w", err)
	}
	return &token, nil
}
