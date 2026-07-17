package tmx

import "net/http"

// Client is a http client for making requests to the Trackmania Exchange API.
type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{},
	}
}
