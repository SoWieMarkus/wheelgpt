package plugin

import (
	"context"
	"crypto/subtle"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SoWieMarkus/wheelgpt/internal/auth"
	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/guess"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

// TokenPrefix marks plugin tokens issued by this backend.
const TokenPrefix = "wgpt_"

// Announcer posts the result of a round into the chat after a delay.
type Announcer interface {
	Announce(channelID, text string, delay time.Duration)
}

// Recorder counts personal bests (metrics).
type Recorder interface {
	PersonalBest(outcome string)
}

type Config struct {
	// LegacySecret verifies plugin tokens of the old backend (HS256 JWTs).
	// Empty rejects them.
	LegacySecret string
}

type Service struct {
	repo      Repository
	channels  *channel.Service
	ingame    *ingame.Service
	guesses   *guess.Service
	announcer Announcer
	rec       Recorder
	cfg       Config
	now       func() time.Time
}

func NewService(repo Repository, channels *channel.Service, in *ingame.Service, guesses *guess.Service, announcer Announcer, rec Recorder, cfg Config) *Service {
	return &Service{repo: repo, channels: channels, ingame: in, guesses: guesses, announcer: announcer, rec: rec, cfg: cfg, now: time.Now}
}

// Caller is an authenticated plugin.
type Caller struct {
	ChannelID string
	// Legacy is true for a token of the old backend.
	Legacy bool
}

// Authenticate checks a plugin token: a static token of this backend or a
// JWT of the old one.
func (s *Service) Authenticate(ctx context.Context, token string) (Caller, error) {
	if token == "" || len(token) > 4096 {
		return Caller{}, domain.ErrInvalidPluginToken
	}
	if strings.HasPrefix(token, TokenPrefix) {
		t, ok, err := s.repo.TokenByHash(ctx, auth.Hash(token))
		if err != nil {
			return Caller{}, err
		}
		if !ok {
			return Caller{}, domain.ErrInvalidPluginToken
		}
		if err := s.repo.TouchToken(ctx, t.ChannelID, s.now()); err != nil {
			slog.WarnContext(ctx, "touch plugin token failed", "error", err)
		}
		return Caller{ChannelID: t.ChannelID}, nil
	}
	return s.authenticateLegacy(ctx, token)
}

type legacyClaims struct {
	ID        string `json:"id"`
	ChannelID string `json:"channelId"` // login, used by very old tokens
	Token     string `json:"token"`
	jwt.RegisteredClaims
}

// authenticateLegacy verifies a token signed by the old backend: an HS256 JWT
// with the channel id (or, older, the login) and the channel's secret.
func (s *Service) authenticateLegacy(ctx context.Context, token string) (Caller, error) {
	if s.cfg.LegacySecret == "" {
		return Caller{}, domain.ErrInvalidPluginToken
	}
	var claims legacyClaims
	_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
		return []byte(s.cfg.LegacySecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || (claims.ID == "") == (claims.ChannelID == "") || claims.Token == "" {
		return Caller{}, domain.ErrInvalidPluginToken
	}
	ref := claims.ID
	if ref == "" {
		ref = claims.ChannelID
	}
	c, err := s.channels.Resolve(ctx, ref)
	if errors.Is(err, domain.ErrChannelNotFound) {
		return Caller{}, domain.ErrInvalidPluginToken
	}
	if err != nil {
		return Caller{}, err
	}
	if c.LegacySecret == "" || subtle.ConstantTimeCompare([]byte(c.LegacySecret), []byte(claims.Token)) != 1 {
		return Caller{}, domain.ErrInvalidPluginToken
	}
	return Caller{ChannelID: c.ID, Legacy: true}, nil
}

// TokenInfo describes the plugin token of a channel.
type TokenInfo struct {
	Token             *Token
	LegacyTokenActive bool
}

func (s *Service) TokenInfo(ctx context.Context, channelID string) (TokenInfo, error) {
	t, ok, err := s.repo.Token(ctx, channelID)
	if err != nil {
		return TokenInfo{}, err
	}
	if ok {
		return TokenInfo{Token: &t}, nil
	}
	c, err := s.channels.Get(ctx, channelID)
	if err != nil {
		return TokenInfo{}, err
	}
	return TokenInfo{LegacyTokenActive: s.cfg.LegacySecret != "" && c.LegacySecret != ""}, nil
}

// CreateToken issues a new token. It is returned once, only its hash is
// stored. The previous token and old-style tokens stop working.
func (s *Service) CreateToken(ctx context.Context, channelID string) (string, time.Time, error) {
	token := TokenPrefix + auth.RandomToken(32)
	now := s.now().UTC().Truncate(time.Millisecond)
	return token, now, s.repo.ReplaceToken(ctx, channelID, auth.Hash(token), now)
}

// Status is what the plugin learns about its channel.
type Status struct {
	Channel   channel.Channel
	Accepting bool
}

func (s *Service) Status(ctx context.Context, channelID string) (Status, error) {
	c, err := s.channels.Get(ctx, channelID)
	return Status{Channel: c, Accepting: c.Accepting()}, err
}

// live loads the channel and rejects plugin data while it is not live.
func (s *Service) live(ctx context.Context, channelID string) (channel.Channel, error) {
	c, err := s.channels.Get(ctx, channelID)
	if err != nil {
		return c, err
	}
	if !c.Accepting() {
		return c, domain.ErrChannelOffline
	}
	return c, nil
}

// SetMap stores the current map.
func (s *Service) SetMap(ctx context.Context, channelID string, m trackmania.Map) (trackmania.Map, error) {
	if _, err := s.live(ctx, channelID); err != nil {
		return trackmania.Map{}, err
	}
	return s.ingame.SetMap(ctx, channelID, m)
}

func (s *Service) ClearMap(ctx context.Context, channelID string) error {
	if _, err := s.live(ctx, channelID); err != nil {
		return err
	}
	return s.ingame.ClearMap(ctx, channelID)
}

func (s *Service) SetRoom(ctx context.Context, channelID string, r ingame.Room) error {
	if _, err := s.live(ctx, channelID); err != nil {
		return err
	}
	return s.ingame.SetRoom(ctx, channelID, r)
}

func (s *Service) ClearRoom(ctx context.Context, channelID string) error {
	if _, err := s.live(ctx, channelID); err != nil {
		return err
	}
	return s.ingame.ClearRoom(ctx, channelID)
}

// PersonalBestResult is what happened to a reported personal best.
type PersonalBestResult struct {
	// Reason explains why the result is not announced, if it is not.
	Reason     string
	AnnounceAt *time.Time
	Round      guess.Round
}

// PersonalBest ends the guessing round and schedules the chat announcement
// after the channel's delay.
func (s *Service) PersonalBest(ctx context.Context, channelID string, pb trackmania.Time) (PersonalBestResult, error) {
	c, err := s.live(ctx, channelID)
	if errors.Is(err, domain.ErrChannelOffline) {
		s.rec.PersonalBest("offline")
	}
	if err != nil {
		return PersonalBestResult{}, err
	}
	round, err := s.guesses.EvaluatePB(ctx, c, pb)
	if err != nil {
		return PersonalBestResult{}, err
	}
	s.rec.PersonalBest("evaluated")
	out := PersonalBestResult{Round: round}
	if !c.Bot.Enabled {
		out.Reason = "The bot is disconnected, the result is not announced."
		return out, nil
	}
	delay := c.Settings.GuessDelay()
	s.announcer.Announce(c.ID, round.Message, delay)
	at := s.now().Add(delay).UTC()
	out.AnnounceAt = &at
	return out, nil
}
