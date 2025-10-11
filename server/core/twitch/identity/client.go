package identity

import (
	"net/http"

	"github.com/SoWieMarkus/wheelgpt/core/config"
)

type Client struct {
	client  *http.Client
	baseURL string
	config  *config.TwitchConfig
}

func NewClient(config *config.TwitchConfig) *Client {
	return &Client{
		client:  &http.Client{},
		baseURL: "https://id.twitch.tv/oauth2",
		config:  config,
	}
}
