package helix

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/identity"
)

type Client struct {
	client      *http.Client
	baseURL     string
	config      *config.TwitchConfig
	identity    *identity.Client
	bearerToken *string
}

func NewClient(config *config.TwitchConfig) *Client {
	return &Client{
		client:   &http.Client{},
		baseURL:  "https://api.twitch.tv/helix",
		config:   config,
		identity: identity.NewClient(config),
	}
}

func (c *Client) getBearerToken() string {
	return ""
}

func (c *Client) get(endpoint string, result any, params map[string]string) (*http.Response, error) {
	url := c.baseURL + endpoint

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bearerToken := c.getBearerToken()
	request.Header.Set("Authorization", "Bearer "+bearerToken)
	request.Header.Set("Content-Type", "application/json")

	q := request.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	request.URL.RawQuery = q.Encode()

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check status code
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", response.StatusCode)
	}

	// Decode JSON into the provided interface
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, nil
}
