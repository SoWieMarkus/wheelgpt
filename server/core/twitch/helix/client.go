package helix

import (
	"fmt"
	"time"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/http"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/identity"
)

type Client struct {
	client           *http.Client
	config           *config.TwitchConfig
	identity         identity.TwitchIdentityAPI
	appToken         *identity.AppAccessToken
	appTokenIssuedAt *time.Time
}

func NewClient(config *config.TwitchConfig) *Client {
	defaultHeaders := map[string]string{
		"Client-ID": config.ClientID,
	}

	httpClient := http.NewClient(
		"https://api.twitch.tv/helix",
		&defaultHeaders,
	)

	identityClient := identity.NewClient(config)

	return &Client{
		client:   httpClient,
		config:   config,
		identity: identityClient,
	}
}

func (c *Client) getAppToken() (string, error) {
	now := time.Now()

	if c.appToken == nil || c.appTokenIssuedAt == nil || c.appToken.IsExpired(c.appTokenIssuedAt) {
		token, err := c.identity.RequestAppAccessToken()
		if err != nil {
			return "", fmt.Errorf("failed to get app access token: %w", err)
		}
		c.appToken = token
		c.appTokenIssuedAt = &now
	}

	return c.appToken.AccessToken, nil
}
