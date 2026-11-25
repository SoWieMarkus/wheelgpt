package twitch

import (
	"context"
	"net/http"
	"net/url"
)

// EventSub subscription types the bot uses.
const (
	SubStreamOnline  = "stream.online"
	SubStreamOffline = "stream.offline"
	SubChatMessage   = "channel.chat.message"
)

// Subscription is an EventSub subscription.
type Subscription struct {
	ID        string            `json:"id"`
	Status    string            `json:"status"`
	Type      string            `json:"type"`
	Version   string            `json:"version"`
	Condition map[string]string `json:"condition"`
	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
	} `json:"transport"`
}

// Subscriptions lists the subscriptions of the app (follows pagination). A
// non-empty userID only returns subscriptions whose condition names that user.
func (c *Client) Subscriptions(ctx context.Context, userID string) ([]Subscription, error) {
	var all []Subscription
	cursor := ""
	for {
		q := url.Values{}
		if userID != "" {
			q.Set("user_id", userID)
		}
		if cursor != "" {
			q.Set("after", cursor)
		}
		var res struct {
			Data       []Subscription
			Pagination struct {
				Cursor string `json:"cursor"`
			}
		}
		if err := c.helix(ctx, http.MethodGet, "/eventsub/subscriptions", q, nil, &res, ""); err != nil {
			return all, err
		}
		all = append(all, res.Data...)
		if res.Pagination.Cursor == "" || len(res.Data) == 0 {
			return all, nil
		}
		cursor = res.Pagination.Cursor
	}
}

// Subscribe creates a webhook subscription. A 409 (already exists) is not an
// error. A 403 means the required authorization is missing.
func (c *Client) Subscribe(ctx context.Context, typ string, condition map[string]string, callback, secret string) error {
	body := map[string]any{
		"type":      typ,
		"version":   "1",
		"condition": condition,
		"transport": map[string]string{"method": "webhook", "callback": callback, "secret": secret},
	}
	err := c.helix(ctx, http.MethodPost, "/eventsub/subscriptions", nil, body, nil, "")
	if StatusOf(err) == http.StatusConflict {
		return nil
	}
	return err
}

// Unsubscribe deletes a subscription.
func (c *Client) Unsubscribe(ctx context.Context, id string) error {
	err := c.helix(ctx, http.MethodDelete, "/eventsub/subscriptions", url.Values{"id": {id}}, nil, nil, "")
	if StatusOf(err) == http.StatusNotFound {
		return nil
	}
	return err
}
