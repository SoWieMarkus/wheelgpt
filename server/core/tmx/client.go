package tmx

import "github.com/SoWieMarkus/wheelgpt/core/http"

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: http.NewClient("https://trackmania.exchange/api", nil),
	}
}
