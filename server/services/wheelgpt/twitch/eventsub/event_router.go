package eventsub

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch/eventsub/events"
)

type EventSubMessage struct {
	Subscription SubscriptionInfo `json:"subscription"`
	Event        json.RawMessage  `json:"event"`
}

type SubscriptionInfo struct {
	ID        string         `json:"id"`
	Status    string         `json:"status"`
	Type      string         `json:"type"`
	Version   string         `json:"version"`
	Condition map[string]any `json:"condition"`
	Transport TransportInfo  `json:"transport"`
	CreatedAt time.Time      `json:"created_at"`
	Cost      int            `json:"cost"`
}

type TransportInfo struct {
	Method    string `json:"method"`
	SessionID string `json:"session_id"`
}

type EventHandler interface {
	// Callback for chat message event
	OnChatMessage(event events.ChatMessageEvent) error
	// Callback for follow event
	OnFollow(event events.ChannelFollowEvent) error
	// Callback for subscription event
	OnSubscription(event events.ChannelSubscribeEvent) error
	// Callback for stream online event
	OnStreamOnline(event events.StreamOnlineEvent) error
	// Callback for stream offline event
	OnStreamOffline(event events.StreamOfflineEvent) error
	// Handle unknown event types
	OnUnknownEvent(eventType string, rawData json.RawMessage) error
}

// EventRouter routes incoming EventSub events to the appropriate handler methods
type EventRouter struct {
	handler EventHandler
}

func NewEventRouter(handler EventHandler) *EventRouter {
	return &EventRouter{
		handler: handler,
	}
}

func (er *EventRouter) RouteEvent(notification EventSubMessage) error {
	eventType := notification.Subscription.Type

	slog.Debug("Routing event",
		"type", eventType,
		"subscription_id", notification.Subscription.ID)

	switch eventType {
	case "channel.chat.message":
		return er.handleChatMessage(notification.Event)
	case "channel.follow":
		return er.handleFollow(notification.Event)
	case "channel.subscribe":
		return er.handleSubscription(notification.Event)
	case "stream.online":
		return er.handleStreamOnline(notification.Event)
	case "stream.offline":
		return er.handleStreamOffline(notification.Event)
	default:
		slog.Warn("Unknown event type", "type", eventType)
		return er.handler.OnUnknownEvent(eventType, notification.Event)
	}
}

func (er *EventRouter) handleChatMessage(rawEvent json.RawMessage) error {
	var event events.ChatMessageEvent
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		return fmt.Errorf("failed to unmarshal chat message: %w", err)
	}
	return er.handler.OnChatMessage(event)
}

func (er *EventRouter) handleFollow(rawEvent json.RawMessage) error {
	var event events.ChannelFollowEvent
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		return fmt.Errorf("failed to unmarshal follow event: %w", err)
	}
	return er.handler.OnFollow(event)
}

func (er *EventRouter) handleSubscription(rawEvent json.RawMessage) error {
	var event events.ChannelSubscribeEvent
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		return fmt.Errorf("failed to unmarshal subscription event: %w", err)
	}
	return er.handler.OnSubscription(event)
}

func (er *EventRouter) handleStreamOnline(rawEvent json.RawMessage) error {
	var event events.StreamOnlineEvent
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		return fmt.Errorf("failed to unmarshal stream online event: %w", err)
	}
	return er.handler.OnStreamOnline(event)
}

func (er *EventRouter) handleStreamOffline(rawEvent json.RawMessage) error {
	var event events.StreamOfflineEvent
	if err := json.Unmarshal(rawEvent, &event); err != nil {
		return fmt.Errorf("failed to unmarshal stream offline event: %w", err)
	}
	return er.handler.OnStreamOffline(event)
}
