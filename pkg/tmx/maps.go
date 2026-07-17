package tmx

import (
	"context"
	"encoding/json"
	"net/http"
)

// Map represents a trackmania map retrieved from the Trackmania Exchange API.
// Most of the fields are omitted. See https://api2.mania.exchange/Method/Index/37 for more details.
type Map struct {
	TrackID string `json:"track_id"`
}

// GetMapByUID retrieves a map from the Trackmania Exchange API by its UID.
func (c *Client) GetMapByUID(ctx context.Context, mapUID string) (*Map, error) {
	url := "https://trackmania.exchange/api/maps/get_map_info/uid/" + mapUID
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var m Map
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
