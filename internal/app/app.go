// Package app is the composition root: it loads the configuration, wires
// repositories, services and handlers together, builds the router and runs
// the HTTP servers and background jobs until the process is asked to stop.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/SoWieMarkus/wheelgpt/internal/auth"
	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/eventsub"
	"github.com/SoWieMarkus/wheelgpt/internal/guess"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/leaderboard"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/metrics"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/tmx"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/twitch"
	"github.com/SoWieMarkus/wheelgpt/internal/plugin"
)

const shutdownTimeout = 20 * time.Second

// Run starts the application and blocks until it receives SIGINT or SIGTERM.
func Run(version string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	setupLogger(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := os.MkdirAll(filepath.Dir(cfg.DatabasePath()), 0o750); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}
	db, err := store.Open(ctx, cfg.DatabasePath())
	if err != nil {
		return err
	}
	defer db.Close()
	if err := store.Migrate(ctx, db); err != nil {
		return err
	}
	st := store.NewStore(db)
	m := metrics.New(version, db)

	ext := externals{tmx: tmx.New()}
	if cfg.TwitchEnabled() {
		ext.twitch = twitch.New(cfg.TwitchClientID, cfg.TwitchClientSecret)
		ext.twitch.OnError = m.TwitchFailure
	}
	c, err := build(cfg, version, st, m, ext)
	if err != nil {
		return err
	}

	apiServer := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           c.router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      35 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	metricsMux := http.NewServeMux()
	metricsMux.Handle("GET /metrics", m.Handler())
	metricsServer := &http.Server{
		Addr:              cfg.MetricsAddr,
		Handler:           metricsMux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if !cfg.SecureCookies() {
		slog.Warn("PUBLIC_URL is not https - session cookies are sent without the Secure flag. Only use this for local development.")
	}
	if !cfg.TwitchEnabled() {
		slog.Warn("Twitch is not configured - logins are disabled and chat replies are only logged.")
	}
	if cfg.DevMode {
		slog.Warn("DEV_MODE is enabled - anyone can log in as any channel at /api/v1/auth/dev and fake the live status.")
	}

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error { return serve(apiServer, "api", version) })
	g.Go(func() error { return serve(metricsServer, "metrics", version) })
	g.Go(func() error {
		c.startup(ctx)
		return nil
	})
	g.Go(func() error {
		return every(ctx, time.Hour, "purge sessions", c.auth.PurgeExpiredSessions)
	})
	g.Go(func() error {
		<-ctx.Done()
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
		defer cancel()
		err := errors.Join(apiServer.Shutdown(shutdownCtx), metricsServer.Shutdown(shutdownCtx))
		// Pending PB announcements are sent right away.
		c.bot.Stop()
		c.ingame.Wait()
		return err
	})
	return g.Wait()
}

func serve(srv *http.Server, name, version string) error {
	slog.Info("listening", "server", name, "addr", srv.Addr, "version", version)
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return fmt.Errorf("%s server: %w", name, err)
}

// every runs job once right away and then at every interval until ctx ends.
// Job errors are logged, they never stop the server.
func every(ctx context.Context, interval time.Duration, name string, job func(context.Context) error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		if err := job(ctx); err != nil && ctx.Err() == nil {
			slog.ErrorContext(ctx, "background job failed", "job", name, "error", err)
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func setupLogger(cfg Config) {
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	var handler slog.Handler = slog.NewJSONHandler(os.Stdout, opts)
	if cfg.LogFormat == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	slog.SetDefault(slog.New(handler))
}

// externals are the clients of outside services. Tests pass fakes.
type externals struct {
	twitch *twitch.Client // nil if Twitch is not configured
	tmx    ingame.TMX
	sender bot.Sender // nil: Helix with the bot account, or logging only
}

// components are the wired parts Run needs besides the router.
type components struct {
	router    http.Handler
	auth      *auth.Service
	bot       *bot.Bot
	ingame    *ingame.Service
	channels  *channel.Service
	plugin    *plugin.Service
	chRepo    *channel.SQLRepository
	subs      *eventsub.Manager
	twitch    *twitch.Client
	botUserID *atomic.Value
	cfg       Config
}

// build wires repositories, services and handlers.
func build(cfg Config, version string, st *store.Store, m *metrics.Metrics, ext externals) (components, error) {
	botUserID := &atomic.Value{}
	botUserID.Store("")
	getBotID := func() string { return botUserID.Load().(string) }

	var subsClient eventsub.Subscriptions = noSubscriptions{}
	if ext.twitch != nil {
		subsClient = ext.twitch
	}
	chRepo := channel.NewRepository(st)
	subs := eventsub.NewManager(subsClient, chRepo, eventsub.ManagerConfig{
		Callback:  cfg.PublicOrigin() + "/api/v1/twitch/eventsub",
		Secret:    cfg.EventSubSecret,
		BotUserID: getBotID,
		Active:    cfg.EventSubSync,
	})
	channels := channel.NewService(chRepo, subs)
	inGame := ingame.NewService(ingame.NewRepository(st), ext.tmx)
	guesses := guess.NewService(guess.NewRepository(st), inGame, m)
	boards := leaderboard.NewService(leaderboard.NewRepository(st))

	sender := ext.sender
	if sender == nil {
		sender = newChatSender(ext.twitch, getBotID, m)
	}
	b := bot.New(bot.Deps{
		Channels: channels, Guesses: guesses, Ingame: inGame, Leaderboard: boards,
		Sender: sender, Recorder: m, PublicURL: cfg.PublicOrigin(),
	})

	var oauth auth.OAuth
	if ext.twitch != nil {
		oauth = ext.twitch
	}
	authService := auth.NewService(auth.NewRepository(st), oauth, channels, m, auth.Config{
		SessionTTL:  cfg.SessionTTL,
		BotUsername: cfg.BotUsername,
		OnBotAuthorized: func(id string) {
			botUserID.Store(id)
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				if err := subs.Sync(ctx); err != nil {
					slog.ErrorContext(ctx, "eventsub sync after bot authorization failed", "error", err)
				}
			}()
		},
	})
	pluginService := plugin.NewService(plugin.NewRepository(st), channels, inGame, guesses, b, m,
		plugin.Config{LegacySecret: cfg.LegacyJWTSecret})

	streamEvents := eventsub.NewEvents(channels, chRepo, inGame, b, getBotID)
	events := eventsub.NewHandler(cfg.EventSubSecret, streamEvents)
	events.OnMessage = m.EventSub

	router, err := newRouter(routerConfig{
		version:           version,
		publicOrigin:      cfg.PublicOrigin(),
		trustProxyHeaders: cfg.TrustProxyHeaders,
	}, st, m, handlers{
		auth: auth.NewHandler(authService, auth.HandlerConfig{
			PublicOrigin:   cfg.PublicOrigin(),
			SecureCookies:  cfg.SecureCookies(),
			LoginRateLimit: cfg.LoginRateLimit,
			DevLogin:       cfg.DevMode,
		}),
		channel:     channel.NewHandler(channels, streamEvents, channel.HandlerConfig{DevMode: cfg.DevMode, ChatEvents: cfg.EventSubSync}),
		ingame:      ingame.NewHandler(inGame, channels),
		leaderboard: leaderboard.NewHandler(boards, channels),
		plugin:      plugin.NewHandler(pluginService, plugin.HandlerConfig{RateLimit: cfg.PluginRateLimit}),
		bot:         bot.NewHandler(b),
		eventsub:    events,
	})
	if err != nil {
		return components{}, err
	}
	return components{
		router: router, auth: authService, bot: b, ingame: inGame, channels: channels, plugin: pluginService,
		chRepo: chRepo, subs: subs, twitch: ext.twitch, botUserID: botUserID, cfg: cfg,
	}, nil
}

// startup refreshes channel profiles and live status from Twitch and syncs
// the EventSub subscriptions.
func (c components) startup(ctx context.Context) {
	if c.twitch == nil {
		return
	}
	if u, err := c.twitch.UserByLogin(ctx, c.cfg.BotUsername); err != nil {
		slog.ErrorContext(ctx, "look up bot account failed", "bot", c.cfg.BotUsername, "error", err)
	} else {
		c.botUserID.Store(u.ID)
	}

	entries, err := c.chRepo.ListForSync(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "list channels failed", "error", err)
		return
	}
	ids := make([]string, len(entries))
	for i, e := range entries {
		ids[i] = e.ID
	}
	users, err := c.twitch.Users(ctx, ids)
	if err != nil {
		slog.ErrorContext(ctx, "refresh channel profiles failed", "error", err)
	}
	for _, u := range users {
		p := channel.Profile{ID: u.ID, Login: u.Login, DisplayName: u.DisplayName, ProfileImage: u.ProfileImageURL}
		if err := c.channels.UpdateProfile(ctx, p); err != nil {
			slog.WarnContext(ctx, "update channel profile failed", "channel", u.ID, "error", err)
		}
	}
	live, err := c.twitch.LiveUserIDs(ctx, ids)
	if err != nil {
		slog.ErrorContext(ctx, "refresh live status failed", "error", err)
	} else {
		for _, id := range ids {
			if err := c.channels.SetLive(ctx, id, live[id]); err != nil {
				slog.WarnContext(ctx, "set live failed", "channel", id, "error", err)
			}
		}
	}
	slog.InfoContext(ctx, "channels refreshed", "channels", len(ids), "profiles", len(users))
	if err := c.subs.Sync(ctx); err != nil {
		slog.ErrorContext(ctx, "eventsub sync failed", "error", err)
	}
}

// noSubscriptions stands in for Twitch when it is not configured.
type noSubscriptions struct{}

func (noSubscriptions) Subscriptions(context.Context, string) ([]twitch.Subscription, error) {
	return nil, nil
}

func (noSubscriptions) Subscribe(context.Context, string, map[string]string, string, string) error {
	return nil
}

func (noSubscriptions) Unsubscribe(context.Context, string) error { return nil }

// chatSender sends chat messages as the bot account through Helix, or only
// logs them when Twitch is not configured.
type chatSender struct {
	twitch *twitch.Client
	botID  func() string
	m      *metrics.Metrics
}

func newChatSender(tw *twitch.Client, botID func() string, m *metrics.Metrics) chatSender {
	return chatSender{twitch: tw, botID: botID, m: m}
}

func (s chatSender) Send(ctx context.Context, broadcasterID, text string) error {
	if s.twitch == nil {
		slog.InfoContext(ctx, "chat message (not sent, Twitch is not configured)", "channel", broadcasterID, "text", text)
		s.m.ChatMessage("logged")
		return nil
	}
	if s.botID() == "" {
		s.m.ChatMessage("error")
		return errors.New("bot user id unknown, authorize the bot at /api/v1/auth/bot")
	}
	if err := s.twitch.SendChatMessage(ctx, broadcasterID, s.botID(), text); err != nil {
		s.m.ChatMessage("error")
		return err
	}
	s.m.ChatMessage("sent")
	return nil
}
