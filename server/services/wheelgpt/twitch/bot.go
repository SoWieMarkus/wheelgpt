package twitch

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/helix"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/identity"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch/eventsub"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch/websocket"
)

type Bot struct {
	channels      []string
	channelsMutex sync.RWMutex
	helix         *helix.Client
	identity      *identity.Client
	wsClient      *websocket.WebSocketClient
	eventHandler  *eventsub.EventRouter
	running       bool
}

func NewBot(config *config.TwitchConfig, handler eventsub.EventHandler) *Bot {
	bot := &Bot{
		channels:     make([]string, 0),
		helix:        helix.NewClient(config),
		identity:     identity.NewClient(config),
		eventHandler: eventsub.NewEventRouter(handler),
	}

	bot.wsClient = websocket.NewClient(bot)
	return bot
}

func (b *Bot) Join(channel string) {
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	// Check if already joined
	if slices.Contains(b.channels, channel) {
		slog.Warn("Attempted to join channel that the bot is already connected to", "channel", channel)
		return
	}

	b.channels = append(b.channels, channel)
	slog.Info("Added channel to join list", "channel", channel)

	// If WebSocket is connected, register EventSub immediately
	if b.wsClient.IsConnected() && b.wsClient.GetSessionID() != "" {
		go func() {
			// TODO register event sub for channel
		}()
	}
}

func (b *Bot) Leave(channel string) {
	b.channelsMutex.Lock()
	defer b.channelsMutex.Unlock()

	// Remove from channels list
	for i, ch := range b.channels {
		if ch == channel {
			b.channels = append(b.channels[:i], b.channels[i+1:]...)
			slog.Info("Removed channel from join list", "channel", channel)
			break
		}
	}

	// TODO: Unregister EventSub subscription
}

func (b *Bot) Start() {
	b.running = true
	slog.Info("Starting bot...")

	// Start WebSocket connection
	go func() {
		if err := b.wsClient.Connect(); err != nil {
			slog.Error("WebSocket client failed", "error", err)
		}
	}()
}

func (b *Bot) Stop() error {
	b.running = false
	slog.Info("Stopping bot...")
	return b.wsClient.Close()
}

func (b *Bot) GetChannels() []string {
	b.channelsMutex.RLock()
	defer b.channelsMutex.RUnlock()

	channels := make([]string, len(b.channels))
	copy(channels, b.channels)
	return channels
}

// WebSocketMessageHandler Interface Implementation
func (b *Bot) OnSessionWelcome(sessionID string) error {
	slog.Info("Bot received session welcome", "session_id", sessionID)

	// Register EventSub for all channels
	channels := b.GetChannels()
	if len(channels) > 0 {
		go func() {
			for _, channel := range channels {
				slog.Info("Registering EventSub for channel", "channel", channel)
				// TODO register event sub for channel
			}
		}()
	}

	return nil
}

func (b *Bot) OnSessionKeepalive() {
	// TODO add health monitoring here
	slog.Debug("Received keepalive")
}

func (b *Bot) OnSessionReconnect() error {
	slog.Info("Bot handling session reconnect")
	return fmt.Errorf("session reconnect requested")
}

func (b *Bot) OnError(err error) {
	slog.Error("WebSocket error in bot", "error", err)
}

func (b *Bot) OnNotification(message eventsub.EventSubMessage) error {
	b.eventHandler.RouteEvent(message)
	return nil
}
