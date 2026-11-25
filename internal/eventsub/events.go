package eventsub

import (
	"context"
	"errors"
	"log/slog"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/twitch"
)

// Events reacts to EventSub notifications. It implements Dispatcher.
type Events struct {
	channels  *channel.Service
	repo      Channels
	ingame    *ingame.Service
	bot       *bot.Bot
	botUserID func() string
}

func NewEvents(channels *channel.Service, repo Channels, in *ingame.Service, b *bot.Bot, botUserID func() string) *Events {
	return &Events{channels: channels, repo: repo, ingame: in, bot: b, botUserID: botUserID}
}

// StreamOnline marks a channel live.
func (e *Events) StreamOnline(ctx context.Context, channelID string) error {
	slog.InfoContext(ctx, "stream online", "channel", channelID)
	return e.channels.SetLive(ctx, channelID, true)
}

// StreamOffline marks a channel offline and clears map, room and open
// guesses: plugin data only counts while the stream is live.
func (e *Events) StreamOffline(ctx context.Context, channelID string) error {
	slog.InfoContext(ctx, "stream offline", "channel", channelID)
	if err := e.channels.SetLive(ctx, channelID, false); err != nil {
		return err
	}
	return e.ingame.ClearAll(ctx, channelID)
}

// ChatMessage runs the commands of channels with the bot enabled.
func (e *Events) ChatMessage(ctx context.Context, msg bot.Message) error {
	if msg.ChatterID == e.botUserID() {
		return nil
	}
	c, err := e.channels.Get(ctx, msg.BroadcasterID)
	if errors.Is(err, domain.ErrChannelNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if !c.Bot.Enabled {
		return nil
	}
	return e.bot.Handle(ctx, c, msg)
}

// Revoked records that Twitch ended a chat subscription, for example
// because the streamer removed the bot's authorization.
func (e *Events) Revoked(ctx context.Context, sub twitch.Subscription) error {
	if sub.Type != twitch.SubChatMessage {
		return nil
	}
	return e.repo.SetChatSubscription(ctx, sub.Condition["broadcaster_user_id"], channel.ChatNeedsAuthorization)
}
