package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"log/slog"
	"strings"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/twitch"
)

// Scopes requested in the OAuth flows.
var (
	// StreamerScopes is all a streamer grants: the bot may read and write
	// chat in their channel. Identity needs no scope.
	StreamerScopes = []string{"channel:bot"}
	// BotScopes is what the bot account grants once.
	BotScopes = []string{"user:read:chat", "user:write:chat", "user:bot"}
)

// OAuth is the part of the Twitch client the login needs.
type OAuth interface {
	AuthorizeURL(redirectURI string, scopes []string, state string, forceVerify bool) string
	ExchangeCode(ctx context.Context, code, redirectURI string) (twitch.UserToken, error)
	Me(ctx context.Context, userToken string) (twitch.User, error)
}

// Registrar creates or updates a channel after a login.
type Registrar interface {
	Register(ctx context.Context, id, login, displayName, profileImage string, scopes []string) (isNew bool, err error)
}

// Recorder counts logins (metrics).
type Recorder interface {
	Login(outcome string)
}

type Config struct {
	SessionTTL  time.Duration
	BotUsername string
	// OnBotAuthorized is called with the bot's user id after the operator
	// authorized the bot account.
	OnBotAuthorized func(userID string)
}

// Identity is the logged in streamer.
type Identity struct {
	ChannelID string
	SessionID string
}

// Issued is a new session.
type Issued struct {
	Token     string
	ExpiresAt time.Time
	ChannelID string
}

type Service struct {
	repo      Repository
	oauth     OAuth // nil if Twitch is not configured
	registrar Registrar
	rec       Recorder
	cfg       Config
	now       func() time.Time
}

func NewService(repo Repository, oauth OAuth, registrar Registrar, rec Recorder, cfg Config) *Service {
	return &Service{repo: repo, oauth: oauth, registrar: registrar, rec: rec, cfg: cfg, now: time.Now}
}

// AuthorizeURL returns the Twitch URL a browser is sent to.
func (s *Service) AuthorizeURL(redirectURI, state string, bot bool) (string, error) {
	if s.oauth == nil {
		return "", domain.ErrTwitchNotConfigured
	}
	if bot {
		return s.oauth.AuthorizeURL(redirectURI, BotScopes, state, true), nil
	}
	return s.oauth.AuthorizeURL(redirectURI, StreamerScopes, state, false), nil
}

// Login finishes the Twitch login of a streamer. The user token only
// identifies the streamer and is dropped afterwards.
func (s *Service) Login(ctx context.Context, code, redirectURI string) (Issued, error) {
	if s.oauth == nil {
		return Issued{}, domain.ErrTwitchNotConfigured
	}
	token, err := s.oauth.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		s.rec.Login("failure")
		return Issued{}, fmt.Errorf("%w: %w", domain.ErrTwitchFailed, err)
	}
	user, err := s.oauth.Me(ctx, token.AccessToken)
	if err != nil {
		s.rec.Login("failure")
		return Issued{}, fmt.Errorf("%w: %w", domain.ErrTwitchFailed, err)
	}
	isNew, err := s.registrar.Register(ctx, user.ID, user.Login, user.DisplayName, user.ProfileImageURL, token.Scopes)
	if err != nil {
		return Issued{}, err
	}
	outcome := "returning"
	if isNew {
		outcome = "new"
	}
	s.rec.Login(outcome)
	slog.InfoContext(ctx, "streamer logged in", "channel", user.ID, "login", user.Login, "new", isNew, "scopes", token.Scopes)
	return s.issue(ctx, user.ID)
}

// AuthorizeBot finishes the one-time authorization of the bot account.
func (s *Service) AuthorizeBot(ctx context.Context, code, redirectURI string) error {
	if s.oauth == nil {
		return domain.ErrTwitchNotConfigured
	}
	token, err := s.oauth.ExchangeCode(ctx, code, redirectURI)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrTwitchFailed, err)
	}
	user, err := s.oauth.Me(ctx, token.AccessToken)
	if err != nil {
		return fmt.Errorf("%w: %w", domain.ErrTwitchFailed, err)
	}
	if !strings.EqualFold(user.Login, s.cfg.BotUsername) {
		return domain.ErrWrongBotAccount
	}
	slog.InfoContext(ctx, "bot account authorized", "bot", user.Login, "scopes", token.Scopes)
	if s.cfg.OnBotAuthorized != nil {
		s.cfg.OnBotAuthorized(user.ID)
	}
	return nil
}

// DevLogin registers a fake channel and logs it in. Only for local testing.
func (s *Service) DevLogin(ctx context.Context, login string) (Issued, error) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(login))
	id := fmt.Sprintf("9%09d", h.Sum32()%1_000_000_000)
	if _, err := s.registrar.Register(ctx, id, login, login, "", StreamerScopes); err != nil {
		return Issued{}, err
	}
	return s.issue(ctx, id)
}

func (s *Service) issue(ctx context.Context, channelID string) (Issued, error) {
	token := RandomToken(32)
	now := s.now()
	expires := now.Add(s.cfg.SessionTTL).UTC()
	if err := s.repo.CreateSession(ctx, Session{ID: Hash(token), ChannelID: channelID, ExpiresAt: expires}, now); err != nil {
		return Issued{}, err
	}
	return Issued{Token: token, ExpiresAt: expires, ChannelID: channelID}, nil
}

// Authenticate resolves a session cookie.
func (s *Service) Authenticate(ctx context.Context, token string) (Identity, error) {
	if token == "" || len(token) > 200 {
		return Identity{}, domain.ErrUnauthenticated
	}
	sess, err := s.repo.GetSession(ctx, Hash(token), s.now())
	if err != nil {
		return Identity{}, err
	}
	return Identity{ChannelID: sess.ChannelID, SessionID: sess.ID}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteSession(ctx, Hash(token))
}

func (s *Service) PurgeExpiredSessions(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx, s.now())
}

// Hash returns the hex SHA-256 of a secret. Tokens are high entropy random
// values, so a plain hash is enough.
func Hash(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(sum[:])
}

// RandomToken returns n random bytes, base64url encoded.
func RandomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
