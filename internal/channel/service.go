package channel

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"time"
)

// ChatSubscriber connects the bot to a channel's chat (EventSub).
type ChatSubscriber interface {
	// Connect subscribes to stream and chat events and returns the chat
	// subscription state.
	Connect(ctx context.Context, channelID string) (string, error)
	DisconnectChat(ctx context.Context, channelID string) error
	Remove(ctx context.Context, channelID string) error
}

type Service struct {
	repo Repository
	subs ChatSubscriber
	now  func() time.Time
}

func NewService(repo Repository, subs ChatSubscriber) *Service {
	return &Service{repo: repo, subs: subs, now: time.Now}
}

func (s *Service) Get(ctx context.Context, id string) (Channel, error) {
	return s.repo.Get(ctx, id)
}

var numeric = regexp.MustCompile(`^\d+$`)

// Resolve finds a channel by Twitch user id or login.
func (s *Service) Resolve(ctx context.Context, ref string) (Channel, error) {
	if numeric.MatchString(ref) {
		c, err := s.repo.Get(ctx, ref)
		if !errors.Is(err, ErrChannelNotFound) {
			return c, err
		}
	}
	return s.repo.GetByLogin(ctx, ref)
}

func (s *Service) ListPublic(ctx context.Context, liveOnly bool, query string, limit, offset int64) ([]Summary, int64, error) {
	return s.repo.ListPublic(ctx, liveOnly, query, limit, offset)
}

func (s *Service) Stats(ctx context.Context) (Stats, error) {
	return s.repo.Stats(ctx)
}

func (s *Service) SetLive(ctx context.Context, id string, live bool) error {
	return s.repo.SetLive(ctx, id, live)
}

func (s *Service) UpdateProfile(ctx context.Context, p Profile) error {
	return s.repo.UpdateProfile(ctx, p)
}

// Register creates or updates a channel after a Twitch login and connects
// the bot with the new grant.
func (s *Service) Register(ctx context.Context, id, login, displayName, profileImage string, scopes []string) (bool, error) {
	isNew, err := s.repo.Register(ctx, Profile{ID: id, Login: login, DisplayName: displayName, ProfileImage: profileImage}, scopes, s.now())
	if err != nil {
		return false, err
	}
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return isNew, err
	}
	if c.Bot.Enabled {
		if _, err := s.subs.Connect(ctx, id); err != nil {
			slog.WarnContext(ctx, "connect after login failed", "channel", id, "error", err)
		}
	}
	return isNew, nil
}

// SettingsPatch changes some settings; nil fields stay as they are.
type SettingsPatch struct {
	GuessDelaySeconds  *float64
	GuessMinAgeSeconds *float64
	Public             *bool
}

func (s *Service) UpdateSettings(ctx context.Context, id string, p SettingsPatch) (Settings, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		return Settings{}, err
	}
	next := c.Settings
	if p.GuessDelaySeconds != nil {
		next.GuessDelaySeconds = *p.GuessDelaySeconds
	}
	if p.GuessMinAgeSeconds != nil {
		next.GuessMinAgeSeconds = *p.GuessMinAgeSeconds
	}
	if p.Public != nil {
		next.Public = *p.Public
	}
	return next, s.repo.UpdateSettings(ctx, id, next)
}

// ConnectBot enables the bot and subscribes to the chat.
func (s *Service) ConnectBot(ctx context.Context, id string) (Bot, error) {
	if err := s.repo.SetBotEnabled(ctx, id, true); err != nil {
		return Bot{}, err
	}
	if _, err := s.subs.Connect(ctx, id); err != nil {
		return Bot{}, err
	}
	c, err := s.repo.Get(ctx, id)
	return c.Bot, err
}

// DisconnectBot disables the bot and removes the chat subscription.
func (s *Service) DisconnectBot(ctx context.Context, id string) (Bot, error) {
	if err := s.repo.SetBotEnabled(ctx, id, false); err != nil {
		return Bot{}, err
	}
	if err := s.subs.DisconnectChat(ctx, id); err != nil {
		slog.WarnContext(ctx, "remove chat subscription failed", "channel", id, "error", err)
	}
	c, err := s.repo.Get(ctx, id)
	return c.Bot, err
}

// Delete removes the bot from the channel and deletes all its data.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.subs.Remove(ctx, id); err != nil {
		slog.WarnContext(ctx, "remove subscriptions failed", "channel", id, "error", err)
	}
	return s.repo.Delete(ctx, id)
}
