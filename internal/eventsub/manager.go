package eventsub

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/twitch"
)

// Channels is the part of the channel repository the sync needs.
type Channels interface {
	ListForSync(ctx context.Context) ([]channel.SyncEntry, error)
	SetChatSubscription(ctx context.Context, id, status string) error
}

// Subscriptions is the EventSub part of the Twitch client.
type Subscriptions interface {
	Subscriptions(ctx context.Context, userID string) ([]twitch.Subscription, error)
	Subscribe(ctx context.Context, typ string, condition map[string]string, callback, secret string) error
	Unsubscribe(ctx context.Context, id string) error
}

// Manager creates and removes the EventSub subscriptions of channels. It
// implements channel.ChatSubscriber.
type Manager struct {
	twitch    Subscriptions
	channels  Channels
	callback  string
	secret    string
	botUserID func() string
	// active is false in development: no subscription is touched, they
	// belong to the Twitch application of production.
	active bool
}

type ManagerConfig struct {
	Callback  string
	Secret    string
	BotUserID func() string
	Active    bool
}

func NewManager(tw Subscriptions, channels Channels, cfg ManagerConfig) *Manager {
	return &Manager{twitch: tw, channels: channels, callback: cfg.Callback, secret: cfg.Secret, botUserID: cfg.BotUserID, active: cfg.Active}
}

type subKey struct{ typ, broadcaster string }

func managed(typ string) bool {
	return typ == twitch.SubStreamOnline || typ == twitch.SubStreamOffline || typ == twitch.SubChatMessage
}

func (m *Manager) chatCondition(channelID string) map[string]string {
	return map[string]string{"broadcaster_user_id": channelID, "user_id": m.botUserID()}
}

func streamCondition(channelID string) map[string]string {
	return map[string]string{"broadcaster_user_id": channelID}
}

// Sync makes Twitch's subscriptions match the database: stream online and
// offline for every channel, chat messages for channels with the bot enabled.
func (m *Manager) Sync(ctx context.Context) error {
	if !m.active || m.botUserID() == "" {
		return nil
	}
	channels, err := m.channels.ListForSync(ctx)
	if err != nil {
		return err
	}
	existing, err := m.twitch.Subscriptions(ctx, "")
	if err != nil {
		return fmt.Errorf("list subscriptions: %w", err)
	}

	want := map[subKey]bool{}
	for _, c := range channels {
		want[subKey{twitch.SubStreamOnline, c.ID}] = true
		want[subKey{twitch.SubStreamOffline, c.ID}] = true
		if c.BotEnabled {
			want[subKey{twitch.SubChatMessage, c.ID}] = true
		}
	}

	have := map[subKey]bool{}
	for _, s := range existing {
		if !managed(s.Type) {
			continue
		}
		k := subKey{s.Type, s.Condition["broadcaster_user_id"]}
		healthy := s.Status == "enabled" || s.Status == "webhook_callback_verification_pending"
		ours := s.Transport.Callback == m.callback &&
			(s.Type != twitch.SubChatMessage || s.Condition["user_id"] == m.botUserID())
		if want[k] && healthy && ours && !have[k] {
			have[k] = true
			if s.Type == twitch.SubChatMessage {
				m.setChat(ctx, k.broadcaster, channel.ChatEnabled)
			}
			continue
		}
		if err := m.twitch.Unsubscribe(ctx, s.ID); err != nil {
			slog.WarnContext(ctx, "eventsub: remove subscription failed", "id", s.ID, "error", err)
		}
	}

	var created, failed int
	for k := range want {
		if have[k] {
			continue
		}
		if k.typ == twitch.SubChatMessage {
			if m.subscribeChat(ctx, k.broadcaster) == channel.ChatEnabled {
				created++
			} else {
				failed++
			}
			continue
		}
		if err := m.twitch.Subscribe(ctx, k.typ, streamCondition(k.broadcaster), m.callback, m.secret); err != nil {
			failed++
			slog.WarnContext(ctx, "eventsub: subscribe failed", "type", k.typ, "channel", k.broadcaster, "error", err)
			continue
		}
		created++
	}
	slog.InfoContext(ctx, "eventsub: synced", "channels", len(channels), "created", created, "failed", failed)
	return nil
}

// subscribeChat subscribes to a channel's chat and records the result.
func (m *Manager) subscribeChat(ctx context.Context, channelID string) string {
	err := m.twitch.Subscribe(ctx, twitch.SubChatMessage, m.chatCondition(channelID), m.callback, m.secret)
	status := channel.ChatEnabled
	switch {
	case err == nil:
	case twitch.StatusOf(err) == http.StatusForbidden:
		// Neither channel:bot from the broadcaster nor a moderator bot.
		status = channel.ChatNeedsAuthorization
	default:
		slog.WarnContext(ctx, "eventsub: chat subscription failed", "channel", channelID, "error", err)
		status = channel.ChatUnknown
	}
	m.setChat(ctx, channelID, status)
	return status
}

func (m *Manager) setChat(ctx context.Context, channelID, status string) {
	if err := m.channels.SetChatSubscription(ctx, channelID, status); err != nil {
		slog.WarnContext(ctx, "eventsub: store chat state failed", "channel", channelID, "error", err)
	}
}

// Connect subscribes to the streams and chat of a channel.
func (m *Manager) Connect(ctx context.Context, channelID string) (string, error) {
	if !m.active || m.botUserID() == "" {
		m.setChat(ctx, channelID, channel.ChatUnknown)
		return channel.ChatUnknown, nil
	}
	if err := m.removeTypes(ctx, channelID, twitch.SubChatMessage); err != nil {
		return channel.ChatUnknown, err
	}
	for _, typ := range []string{twitch.SubStreamOnline, twitch.SubStreamOffline} {
		if err := m.twitch.Subscribe(ctx, typ, streamCondition(channelID), m.callback, m.secret); err != nil {
			slog.WarnContext(ctx, "eventsub: subscribe failed", "type", typ, "channel", channelID, "error", err)
		}
	}
	status := m.subscribeChat(ctx, channelID)
	if status == channel.ChatUnknown {
		return status, domain.ErrTwitchFailed
	}
	return status, nil
}

// DisconnectChat removes the chat subscription of a channel.
func (m *Manager) DisconnectChat(ctx context.Context, channelID string) error {
	m.setChat(ctx, channelID, channel.ChatDisabled)
	if !m.active {
		return nil
	}
	return m.removeTypes(ctx, channelID, twitch.SubChatMessage)
}

// Remove deletes every subscription of a channel.
func (m *Manager) Remove(ctx context.Context, channelID string) error {
	if !m.active {
		return nil
	}
	return m.removeTypes(ctx, channelID, twitch.SubStreamOnline, twitch.SubStreamOffline, twitch.SubChatMessage)
}

func (m *Manager) removeTypes(ctx context.Context, channelID string, types ...string) error {
	subs, err := m.twitch.Subscriptions(ctx, channelID)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrTwitchFailed, err)
	}
	for _, s := range subs {
		if s.Condition["broadcaster_user_id"] != channelID {
			continue
		}
		for _, t := range types {
			if s.Type == t {
				if err := m.twitch.Unsubscribe(ctx, s.ID); err != nil {
					return fmt.Errorf("%w: %w", domain.ErrTwitchFailed, err)
				}
			}
		}
	}
	return nil
}
