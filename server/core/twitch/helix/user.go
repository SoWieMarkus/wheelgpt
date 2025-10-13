package helix

import (
	"net/url"

	"github.com/SoWieMarkus/wheelgpt/core/http"
)

type User struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageUrl string `json:"profile_image_url"`
	OfflineImageUrl string `json:"offline_image_url"`
	CreatedAt       string `json:"created_at"`
}

// See: https://dev.twitch.tv/docs/api/reference#get-users
func (c *Client) GetUsers(userIds []string) ([]User, error) {
	var users []User

	headers := map[string]string{
		"Client-ID":     c.config.ClientID,
		"Authorization": "Bearer " + func() string { token, _ := c.getAppToken(); return token }(),
	}

	params := url.Values{
		"id": userIds,
	}

	request := http.HttpRequest{
		Endpoint: "/users",
		Headers:  &headers,
		Params:   &params,
	}

	_, err := c.client.Get(&request, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
