package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type subscription struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Condition struct {
		BroadcasterUserID string `json:"broadcaster_user_id"`
	} `json:"condition"`
}

type subscriptionType string

const (
	subscriptionTypeStreamOnline  subscriptionType = "stream.online"
	subscriptionTypeStreamOffline subscriptionType = "stream.offline"
)

type subscriptionCondition struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type subscriptionTransport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

type subscriptionRequest struct {
	Type      subscriptionType `json:"type"`
	Version   string           `json:"version"`
	Condition struct {
		BroadcasterUserID string `json:"broadcaster_user_id"`
	} `json:"condition"`
	Transport struct {
		Method   string `json:"method"`
		Callback string `json:"callback"`
		Secret   string `json:"secret"`
	} `json:"transport"`
}

// listSubscriptions retrieves the list of EventSub subscriptions for the app from the Twitch API.
func (c *Client) listSubscriptions(ctx context.Context) ([]subscription, error) {
	token, err := c.getAppAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.twitch.tv/helix/eventsub/subscriptions", nil)
	if err != nil {
		return nil, err
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)

	request.Header.Set("Authorization", bearerToken)
	request.Header.Set("Client-Id", c.clientID)
	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var result struct {
		Data []subscription `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// registerSubscription registers a new EventSub subscription for the specified broadcaster and subscription type.
func (c *Client) registerSubscription(ctx context.Context, subscriptionType subscriptionType, broadcasterUserID string) error {
	token, err := c.getAppAccessToken(ctx)
	if err != nil {
		return err
	}

	body := subscriptionRequest{
		Type:    subscriptionType,
		Version: "1",
		Condition: subscriptionCondition{
			BroadcasterUserID: broadcasterUserID,
		},
		Transport: subscriptionTransport{
			Method:   "webhook",
			Callback: c.eventSubCallbackURL,
			Secret:   c.eventSubSecret,
		},
	}

	bearerToken := fmt.Sprintf("Bearer %s", token)
	requestBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.twitch.tv/helix/eventsub/subscriptions", strings.NewReader(string(requestBody)))
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", bearerToken)
	request.Header.Set("Client-Id", c.clientID)
	request.Header.Set("Content-Type", "application/json")
	response, err := c.http.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		log.Printf("Failed to register subscription, status code: %d", response.StatusCode)
		return fmt.Errorf("failed to register subscription, status code: %d", response.StatusCode)
	}

	return nil
}
