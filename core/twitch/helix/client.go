package helix

import (
	"net/http"
)

type Client struct {
	client    *http.Client
	baseURL   string
	authToken string
}

func NewClient(authToken string) *Client {
	return &Client{
		client:    &http.Client{},
		baseURL:   "https://api.twitch.tv/helix",
		authToken: authToken,
	}
}
