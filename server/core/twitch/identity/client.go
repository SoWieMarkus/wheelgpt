package identity

import (
	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/http"
)

type Client struct {
	client *http.Client
	config *config.TwitchConfig
}

func NewClient(config *config.TwitchConfig) *Client {
	return &Client{
		client: http.NewClient(
			"https://id.twitch.tv/oauth2",
			nil,
		),
		config: config,
	}
}
