package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch/eventsub"
	"github.com/gorilla/websocket"
)

const EVENTSUB_WEBSOCKET_URL = "wss://eventsub.wss.twitch.tv/ws"

type WebSocketMessageHandler interface {
	OnSessionWelcome(sessionID string) error
	OnSessionKeepalive()
	OnSessionReconnect() error
	OnNotification(message eventsub.EventSubMessage) error
	OnError(err error)
}

type WebSocketMessage struct {
	Metadata MessageMetadata `json:"metadata"`
	Payload  MessagePayload  `json:"payload"`
}

type MessageMetadata struct {
	MessageType         string    `json:"message_type"`
	MessageID           string    `json:"message_id"`
	MessageTimestamp    time.Time `json:"message_timestamp"`
	SubscriptionType    string    `json:"subscription_type"`
	SubscriptionVersion string    `json:"subscription_version"`
}

type MessagePayload struct {
	Session      *SessionData               `json:"session,omitempty"`
	Subscription *eventsub.SubscriptionInfo `json:"subscription,omitempty"`
	Event        json.RawMessage            `json:"event,omitempty"`
}

type SessionData struct {
	ID                      string `json:"id"`
	Status                  string `json:"status"`
	ConnectedAt             string `json:"connected_at"`
	KeepaliveTimeoutSeconds int    `json:"keepalive_timeout_seconds"`
	ReconnectURL            string `json:"reconnect_url,omitempty"`
}

type WebSocketClient struct {
	conn      *websocket.Conn
	sessionID string
	connected bool
	handler   WebSocketMessageHandler
}

func NewClient(handler WebSocketMessageHandler) *WebSocketClient {
	return &WebSocketClient{
		handler: handler,
	}
}

func (ws *WebSocketClient) Connect() error {
	retries := 10
	for {
		err := ws.connectOnce()
		if err == nil {
			return nil
		}
		retries--
		if retries == 0 {
			return fmt.Errorf("max retries reached")
		}

		slog.Error("WebSocket connection failed", "error", err)
		ws.handler.OnError(err)
		slog.Info("Reconnecting in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

func (ws *WebSocketClient) connectOnce() error {
	conn, _, err := websocket.DefaultDialer.Dial(EVENTSUB_WEBSOCKET_URL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	ws.conn = conn
	ws.connected = true
	defer func() {
		ws.connected = false
		conn.Close()
	}()

	slog.Info("Connected to Twitch EventSub WebSocket")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		var socketMsg WebSocketMessage
		if err := json.Unmarshal(message, &socketMsg); err != nil {
			slog.Error("Failed to unmarshal message", "error", err)
			continue
		}

		if err := ws.handleMessage(socketMsg); err != nil {
			return err
		}
	}
}

func (ws *WebSocketClient) handleMessage(message WebSocketMessage) error {
	switch message.Metadata.MessageType {
	case "session_welcome":
		if message.Payload.Session == nil {
			return fmt.Errorf("session_welcome without session data")
		}
		ws.sessionID = message.Payload.Session.ID
		return ws.handler.OnSessionWelcome(ws.sessionID)

	case "session_keepalive":
		ws.handler.OnSessionKeepalive()

	case "session_reconnect":
		return ws.handler.OnSessionReconnect()

	case "notification":
		if message.Payload.Subscription == nil {
			return fmt.Errorf("notification without subscription")
		}

		notification := eventsub.EventSubMessage{
			Subscription: *message.Payload.Subscription,
			Event:        message.Payload.Event,
		}

		return ws.handler.OnNotification(notification)

	default:
		slog.Info("Unknown message type", "type", message.Metadata.MessageType)
	}

	return nil
}

func (ws *WebSocketClient) GetSessionID() string {
	return ws.sessionID
}

func (ws *WebSocketClient) IsConnected() bool {
	return ws.connected
}

func (ws *WebSocketClient) Close() error {
	if ws.conn != nil {
		return ws.conn.Close()
	}
	return nil
}
