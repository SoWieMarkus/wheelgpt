package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

const eventsubBase = helixBase + "/eventsub/subscriptions"

type subscription struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Condition struct {
		BroadcasterUserID string `json:"broadcaster_user_id"`
	} `json:"condition"`
}

func (c *Client) SyncWebhooks(ctx context.Context, channelIDs []string, callbackURL, secret string) error {
	existing, err := c.listSubscriptions(ctx)
	if err != nil {
		return err
	}

	// Delete non-enabled subscriptions
	for _, sub := range existing {
		if sub.Status != "enabled" {
			_ = c.deleteSubscription(ctx, sub.ID)
		}
	}

	// Build required set
	type key struct{ typ, id string }
	required := make(map[key]bool)
	for _, id := range channelIDs {
		required[key{"stream.online", id}] = true
		required[key{"stream.offline", id}] = true
	}

	// Mark existing enabled ones
	for _, sub := range existing {
		if sub.Status == "enabled" {
			delete(required, key{sub.Type, sub.Condition.BroadcasterUserID})
		}
	}

	// Register missing
	for k := range required {
		if err := c.createSubscription(ctx, k.typ, k.id, callbackURL, secret); err != nil {
			slog.Error("eventsub: create failed", "type", k.typ, "channel", k.id, "err", err)
		}
	}
	return nil
}

func (c *Client) AddWebhooks(ctx context.Context, channelID, callbackURL, secret string) error {
	// Remove stale ones for this channel first
	existing, err := c.listSubscriptions(ctx)
	if err != nil {
		return err
	}
	for _, sub := range existing {
		if sub.Condition.BroadcasterUserID == channelID {
			_ = c.deleteSubscription(ctx, sub.ID)
		}
	}
	if err := c.createSubscription(ctx, "stream.online", channelID, callbackURL, secret); err != nil {
		return err
	}
	return c.createSubscription(ctx, "stream.offline", channelID, callbackURL, secret)
}

func (c *Client) RemoveWebhooks(ctx context.Context, channelID string) error {
	existing, err := c.listSubscriptions(ctx)
	if err != nil {
		return err
	}
	for _, sub := range existing {
		if sub.Condition.BroadcasterUserID == channelID {
			if err := c.deleteSubscription(ctx, sub.ID); err != nil {
				slog.Error("eventsub: delete failed", "id", sub.ID, "err", err)
			}
		}
	}
	return nil
}

func (c *Client) listSubscriptions(ctx context.Context) ([]subscription, error) {
	token, err := c.appAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, eventsubBase, nil)
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Data []subscription `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out.Data, nil
}

func (c *Client) createSubscription(ctx context.Context, typ, channelID, callbackURL, secret string) error {
	token, err := c.appAccessToken(ctx)
	if err != nil {
		return err
	}
	body := fmt.Sprintf(`{
		"type": %q,
		"version": "1",
		"condition": {"broadcaster_user_id": %q},
		"transport": {"method": "webhook", "callback": %q, "secret": %q}
	}`, typ, channelID, callbackURL, secret)

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, eventsubBase,
		strings.NewReader(body))
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("eventsub create: status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) deleteSubscription(ctx context.Context, id string) error {
	token, err := c.appAccessToken(ctx)
	if err != nil {
		return err
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete,
		eventsubBase+"?id="+id, nil)
	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
