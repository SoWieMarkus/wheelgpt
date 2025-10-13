package tmx

import (
	"net/url"

	"github.com/SoWieMarkus/wheelgpt/core/http"
)

type MapInfo struct {
	TrackID string `json:"TrackID"`
}

func (c *Client) GetMapInfo(mapID string) (*MapInfo, error) {
	request := http.HttpRequest{
		Endpoint: "/map/" + url.PathEscape(mapID),
	}

	var mapInfo MapInfo
	_, err := c.client.Get(&request, &mapInfo)
	if err != nil {
		return nil, err
	}
	return &mapInfo, nil
}
