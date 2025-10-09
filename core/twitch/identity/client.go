package identity

import "net/http"

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient() Client {
	return Client{
		client:  &http.Client{},
		baseURL: "https://id.twitch.tv/oauth2",
	}
}
