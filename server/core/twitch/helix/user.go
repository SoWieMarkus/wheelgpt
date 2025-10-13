package helix

import (
	"fmt"
	"net/url"

	"github.com/SoWieMarkus/wheelgpt/core/http"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/identity"
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

	token, err := c.getAppToken()
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"Client-ID":     c.config.ClientID,
		"Authorization": "Bearer " + token,
	}

	params := url.Values{
		"id": userIds,
	}

	request := http.HttpRequest{
		Endpoint: "/users",
		Headers:  &headers,
		Params:   &params,
	}

	_, err = c.client.Get(&request, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// See: https://dev.twitch.tv/docs/api/reference#get-users
func (c *Client) GetUserByAccessToken(token *identity.UserAccessToken) (*User, error) {
	if token == nil {
		return nil, fmt.Errorf("UserAccessToken is nil")
	}

	var users []User

	headers := map[string]string{
		"Client-ID":     c.config.ClientID,
		"Authorization": "Bearer " + token.AccessToken,
	}

	request := http.HttpRequest{
		Endpoint: "/users",
		Headers:  &headers,
	}

	_, err := c.client.Get(&request, &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no user found for the provided access token")
	}
	user := &users[0]
	return user, nil
}
