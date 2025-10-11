package wheelgpt

import (
	"encoding/json"
	"log/slog"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch/eventsub/events"
)

type WheelGPT struct {
	bot *twitch.Bot
	db  *db.Database
}

func NewInstance(db *db.Database, config *config.TwitchConfig) *WheelGPT {
	wheelgpt := &WheelGPT{
		db: db,
	}
	wheelgpt.bot = twitch.NewBot(config, wheelgpt)
	return wheelgpt
}

func (w *WheelGPT) JoinChannel(channel string) {
	w.bot.Join(channel)
}

func (w *WheelGPT) LeaveChannel(channel string) {
	w.bot.Leave(channel)
}

func (w *WheelGPT) Start() {
	slog.Info("Starting WheelGPT...")
	w.bot.Start()
}

func (w *WheelGPT) Stop() error {
	slog.Info("Stopping WheelGPT...")
	return w.bot.Stop()
}

func (w *WheelGPT) GetChannels() []string {
	return w.bot.GetChannels()
}

func (w *WheelGPT) OnChatMessage(event events.ChatMessageEvent) error {
	slog.Info("WheelGPT received chat message",
		"channel", event.BroadcasterUserLogin,
		"user", event.ChatterUserLogin,
		"message", event.Message.Text)
	return nil
}

func (w *WheelGPT) OnFollow(event events.ChannelFollowEvent) error {
	slog.Info("New follower",
		"channel", event.BroadcasterUserLogin,
		"follower", event.UserLogin)
	return nil
}

func (w *WheelGPT) OnSubscription(event events.ChannelSubscribeEvent) error {
	slog.Info("New subscription",
		"channel", event.BroadcasterUserLogin,
		"subscriber", event.UserLogin,
		"tier", event.Tier)
	return nil
}

func (w *WheelGPT) OnStreamOnline(event events.StreamOnlineEvent) error {
	slog.Info("Stream went online",
		"channel", event.BroadcasterUserLogin,
		"type", event.Type)
	return nil
}

func (w *WheelGPT) OnStreamOffline(event events.StreamOfflineEvent) error {
	slog.Info("Stream went offline",
		"channel", event.BroadcasterUserLogin)
	return nil
}

func (w *WheelGPT) OnUnknownEvent(eventType string, rawData json.RawMessage) error {
	slog.Warn("Unknown event type received", "type", eventType)
	return nil
}
