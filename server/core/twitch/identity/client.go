package identity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/SoWieMarkus/wheelgpt/core/config"
)

type Client struct {
	client  *http.Client
	BaseURL string
	config  *config.TwitchConfig
}

func NewClient(config *config.TwitchConfig) *Client {
	return &Client{
		client:  &http.Client{},
		BaseURL: "https://id.twitch.tv/oauth2",
		config:  config,
	}
}

func (c *Client) post(endpoint string, data url.Values, result any) error {
	url := c.BaseURL + endpoint

	// Create POST request
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response into result
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
