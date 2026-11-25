// Package tmx looks maps up on Trackmania Exchange.
package tmx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client queries trackmania.exchange.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// New returns a client for trackmania.exchange.
func New() *Client {
	return &Client{BaseURL: "https://trackmania.exchange", HTTP: &http.Client{Timeout: 10 * time.Second}}
}

// MapID returns the TMX id of a map uid. ok is false if the map is not on TMX.
func (c *Client) MapID(ctx context.Context, uid string) (id int64, ok bool, err error) {
	u := c.BaseURL + "/api/maps/get_map_info/uid/" + url.PathEscape(uid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, http.NoBody)
	if err != nil {
		return 0, false, err
	}
	req.Header.Set("User-Agent", "WheelGPT (+https://wheelgpt.dev)")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return 0, false, err
	}
	defer func() { _ = res.Body.Close() }()
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return 0, false, err
	}
	if res.StatusCode == http.StatusNotFound || len(body) == 0 {
		return 0, false, nil
	}
	if res.StatusCode != http.StatusOK {
		return 0, false, fmt.Errorf("tmx: status %d", res.StatusCode)
	}
	var info struct {
		TrackID int64 `json:"TrackID"`
	}
	if err := json.Unmarshal(body, &info); err != nil || info.TrackID == 0 {
		return 0, false, nil
	}
	return info.TrackID, true, nil
}
