// Package channel manages the Twitch channels that use WheelGPT: their
// profile, settings and the bot's connection to their chat.
package channel

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store/dbtime"
)

// Chat subscription states, stored in ChannelState.chatSubscription.
const (
	ChatUnknown            = "unknown"
	ChatEnabled            = "enabled"
	ChatNeedsAuthorization = "needs_authorization"
	ChatDisabled           = "disabled"
	// ChatOff is never stored: the server doesn't manage EventSub
	// subscriptions (EVENTSUB_SYNC is off, as in local development).
	ChatOff = "off"
)

// Channel is a streamer's channel.
type Channel struct {
	ID           string
	Login        string
	DisplayName  string
	ProfileImage string
	IsLive       bool
	CreatedAt    time.Time
	Settings     Settings
	Bot          Bot
	// LegacySecret is the secret inside plugin tokens of the old backend.
	// It is empty for channels that never had one or created a new token.
	LegacySecret string
}

// Accepting reports whether plugin data is processed right now: only while
// the stream is live.
func (c Channel) Accepting() bool { return c.IsLive }

// Settings change how the guessing game works in a channel.
type Settings struct {
	GuessDelaySeconds  float64
	GuessMinAgeSeconds float64
	Public             bool
}

// GuessDelay is the wait between a personal best and its announcement.
func (s Settings) GuessDelay() time.Duration {
	return time.Duration(s.GuessDelaySeconds * float64(time.Second))
}

// GuessMinAge is how old a guess must be to count.
func (s Settings) GuessMinAge() time.Duration {
	return time.Duration(s.GuessMinAgeSeconds * float64(time.Second))
}

// Bot is the state of the bot in a channel. A channel of the old backend
// has Enabled=true, no AuthorizedAt and ChatSubscription "unknown".
type Bot struct {
	Enabled          bool
	Scopes           []string
	AuthorizedAt     *time.Time
	ChatSubscription string
}

// Authorized reports whether the streamer granted channel:bot.
func (b Bot) Authorized() bool {
	if b.AuthorizedAt == nil {
		return false
	}
	for _, s := range b.Scopes {
		if s == "channel:bot" {
			return true
		}
	}
	return false
}

// Summary is the public view of a channel.
type Summary struct {
	ID           string
	Login        string
	DisplayName  string
	ProfileImage string
	IsLive       bool
}

// Profile is what Twitch tells us about a user.
type Profile struct {
	ID           string
	Login        string
	DisplayName  string
	ProfileImage string
}

// SyncEntry is what the EventSub sync needs to know about a channel.
type SyncEntry struct {
	ID         string
	BotEnabled bool
}

// Stats are usage totals.
type Stats struct {
	Channels int64
}

// Repository is the persistence the channel service needs.
type Repository interface {
	Get(ctx context.Context, id string) (Channel, error)
	GetByLogin(ctx context.Context, login string) (Channel, error)
	ListPublic(ctx context.Context, liveOnly bool, query string, limit, offset int64) ([]Summary, int64, error)
	ListForSync(ctx context.Context) ([]SyncEntry, error)
	Register(ctx context.Context, p Profile, scopes []string, now time.Time) (isNew bool, err error)
	UpdateProfile(ctx context.Context, p Profile) error
	SetLive(ctx context.Context, id string, live bool) error
	UpdateSettings(ctx context.Context, id string, s Settings) error
	SetBotEnabled(ctx context.Context, id string, enabled bool) error
	SetChatSubscription(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
	Stats(ctx context.Context) (Stats, error)
}

type SQLRepository struct {
	st *store.Store
}

func NewRepository(st *store.Store) *SQLRepository {
	return &SQLRepository{st: st}
}

func fromRow(r store.GetChannelRow) Channel {
	c := Channel{
		ID: r.ID, Login: r.Login, DisplayName: r.DisplayName, ProfileImage: r.ProfileImage,
		IsLive: r.IsLive, CreatedAt: r.CreatedAt.Time, LegacySecret: r.Token,
		Settings: Settings{
			GuessDelaySeconds:  r.GuessDelayTime,
			GuessMinAgeSeconds: r.GuessMinRequiredAgeTime,
			Public:             r.UsagePublic,
		},
		Bot: Bot{Enabled: r.BotEnabled, Scopes: strings.Fields(r.Scopes), ChatSubscription: r.ChatSubscription},
	}
	if c.Bot.Scopes == nil {
		c.Bot.Scopes = []string{}
	}
	if r.AuthorizedAt != nil {
		t := time.UnixMilli(*r.AuthorizedAt).UTC()
		c.Bot.AuthorizedAt = &t
	}
	if !c.Bot.Enabled {
		c.Bot.ChatSubscription = ChatDisabled
	}
	return c
}

func (r *SQLRepository) Get(ctx context.Context, id string) (Channel, error) {
	row, err := r.st.GetChannel(ctx, id)
	if store.IsNotFound(err) {
		return Channel{}, ErrChannelNotFound
	}
	if err != nil {
		return Channel{}, fmt.Errorf("get channel: %w", err)
	}
	return fromRow(row), nil
}

func (r *SQLRepository) GetByLogin(ctx context.Context, login string) (Channel, error) {
	row, err := r.st.GetChannelByLogin(ctx, login)
	if store.IsNotFound(err) {
		return Channel{}, ErrChannelNotFound
	}
	if err != nil {
		return Channel{}, fmt.Errorf("get channel by login: %w", err)
	}
	return fromRow(store.GetChannelRow(row)), nil
}

func (r *SQLRepository) ListPublic(ctx context.Context, liveOnly bool, query string, limit, offset int64) ([]Summary, int64, error) {
	var live int64
	if liveOnly {
		live = 1
	}
	rows, err := r.st.ListPublicChannels(ctx, store.ListPublicChannelsParams{LiveOnly: live, Query: query, Lim: limit, Off: offset})
	if err != nil {
		return nil, 0, fmt.Errorf("list channels: %w", err)
	}
	total, err := r.st.CountPublicChannels(ctx, store.CountPublicChannelsParams{LiveOnly: live, Query: query})
	if err != nil {
		return nil, 0, fmt.Errorf("count channels: %w", err)
	}
	out := make([]Summary, len(rows))
	for i, row := range rows {
		out[i] = Summary{ID: row.ID, Login: row.Login, DisplayName: row.DisplayName, ProfileImage: row.ProfileImage, IsLive: row.IsLive}
	}
	return out, total, nil
}

func (r *SQLRepository) ListForSync(ctx context.Context) ([]SyncEntry, error) {
	rows, err := r.st.ListChannels(ctx)
	if err != nil {
		return nil, fmt.Errorf("list channels: %w", err)
	}
	out := make([]SyncEntry, len(rows))
	for i, row := range rows {
		out[i] = SyncEntry{ID: row.ID, BotEnabled: row.BotEnabled}
	}
	return out, nil
}

func (r *SQLRepository) Register(ctx context.Context, p Profile, scopes []string, now time.Time) (bool, error) {
	_, err := r.st.GetChannel(ctx, p.ID)
	isNew := store.IsNotFound(err)
	if err != nil && !isNew {
		return false, fmt.Errorf("get channel: %w", err)
	}
	authorizedAt := now.UnixMilli()
	err = r.st.InTx(ctx, func(q *store.Queries) error {
		// Token is the secret of old-style plugin tokens; a new channel never
		// had one. On conflict the existing secret is kept.
		if err := q.UpsertChannel(ctx, store.UpsertChannelParams{
			ID: p.ID, Login: strings.ToLower(p.Login), Token: "", DisplayName: p.DisplayName,
			ProfileImage: p.ProfileImage, CreatedAt: dbtime.From(now),
		}); err != nil {
			return err
		}
		return q.UpsertChannelAuthorization(ctx, store.UpsertChannelAuthorizationParams{
			ChannelID: p.ID, Scopes: strings.Join(scopes, " "), AuthorizedAt: &authorizedAt,
		})
	})
	if err != nil {
		return false, fmt.Errorf("register channel: %w", err)
	}
	return isNew, nil
}

func (r *SQLRepository) UpdateProfile(ctx context.Context, p Profile) error {
	err := r.st.UpdateChannelProfile(ctx, store.UpdateChannelProfileParams{
		ID: p.ID, Login: strings.ToLower(p.Login), DisplayName: p.DisplayName, ProfileImage: p.ProfileImage,
	})
	if err != nil {
		return fmt.Errorf("update channel profile: %w", err)
	}
	return nil
}

func (r *SQLRepository) SetLive(ctx context.Context, id string, live bool) error {
	if err := r.st.SetChannelLive(ctx, store.SetChannelLiveParams{ID: id, IsLive: live}); err != nil {
		return fmt.Errorf("set live: %w", err)
	}
	return nil
}

func (r *SQLRepository) UpdateSettings(ctx context.Context, id string, s Settings) error {
	err := r.st.UpdateChannelSettings(ctx, store.UpdateChannelSettingsParams{
		ID: id, GuessDelayTime: s.GuessDelaySeconds, GuessMinRequiredAgeTime: s.GuessMinAgeSeconds,
		UsagePublic: s.Public,
	})
	if err != nil {
		return fmt.Errorf("update settings: %w", err)
	}
	return nil
}

func (r *SQLRepository) SetBotEnabled(ctx context.Context, id string, enabled bool) error {
	if err := r.st.SetBotEnabled(ctx, store.SetBotEnabledParams{ChannelID: id, BotEnabled: enabled}); err != nil {
		return fmt.Errorf("set bot enabled: %w", err)
	}
	return nil
}

func (r *SQLRepository) SetChatSubscription(ctx context.Context, id, status string) error {
	if err := r.st.SetChatSubscription(ctx, store.SetChatSubscriptionParams{ChannelID: id, Status: status}); err != nil {
		return fmt.Errorf("set chat subscription: %w", err)
	}
	return nil
}

// Delete removes a channel and everything that belongs to it.
func (r *SQLRepository) Delete(ctx context.Context, id string) error {
	return r.st.InTx(ctx, func(q *store.Queries) error {
		if _, err := q.DeleteGuesses(ctx, id); err != nil {
			return err
		}
		for _, step := range []func(context.Context, string) error{
			q.DeleteMap, q.DeleteRoom, q.DeleteChannelResults, q.DeleteChannelRounds,
			q.DeleteChannelAllTime, q.DeletePluginToken, q.DeleteChannelSessions,
			q.DeleteChannelState, q.DeleteChannel,
		} {
			if err := step(ctx, id); err != nil {
				return fmt.Errorf("delete channel: %w", err)
			}
		}
		return nil
	})
}

func (r *SQLRepository) Stats(ctx context.Context) (Stats, error) {
	var s Stats
	var err error
	if s.Channels, err = r.st.CountChannels(ctx); err != nil {
		return s, fmt.Errorf("count channels: %w", err)
	}
	return s, nil
}
