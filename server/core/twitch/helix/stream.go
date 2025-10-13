package helix

import (
	"net/url"
	"time"

	"github.com/SoWieMarkus/wheelgpt/core/http"
)

type Pagination struct {
	Cursor string `json:"cursor"`
}

type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Tags         []string  `json:"tags"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailUrl string    `json:"thumbnail_url"`
	IsMature     bool      `json:"is_mature"`
}

type StreamResponse struct {
	Data       []Stream   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

// See: https://dev.twitch.tv/docs/api/reference#get-streams
func (c *Client) GetStreams(userIds []string) (*StreamResponse, error) {
	var streams StreamResponse
	token, err := c.getAppToken()
	if err != nil {
		return nil, err
	}

	params := url.Values{
		"user_id": userIds,
	}

	headers := map[string]string{
		"Client-ID":     c.config.ClientID,
		"Authorization": "Bearer " + token,
	}

	request := http.HttpRequest{
		Endpoint: "/streams",
		Headers:  &headers,
		Params:   &params,
	}

	_, err = c.client.Get(&request, &streams)
	if err != nil {
		return nil, err
	}
	return &streams, nil
}
