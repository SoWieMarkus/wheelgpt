package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/api"
	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/bot/commands"
	"github.com/SoWieMarkus/wheelgpt/internal/config"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
	"github.com/SoWieMarkus/wheelgpt/internal/twitch"
)

func main() {
	cfg := config.Load()

	// Database
	conn := db.Open(cfg.DatabaseURL)
	if err := db.Migrate(conn); err != nil {
		slog.Error("migration failed", "err", err)
		os.Exit(1)
	}
	store := db.NewStore(conn)

	// Twitch clientfunc
	tc := twitch.NewClient(cfg.TwitchClientID, cfg.TwitchClientSecret, cfg.TwitchRedirectURL)

	// Bot — CommandFactory breaks the bot ↔ commands import cycle.
	// main imports both; neither imports the other.
	factory := func(channelID string, s *db.Store) *bot.Registry {
		reg := bot.NewRegistry()
		reg.Register(&commands.GuessCmd{Store: s}, 0)
		reg.Register(&commands.MyGuessCmd{Store: s}, 0)
		reg.Register(&commands.ResultCommand{Store: s}, 0)
		reg.Register(&commands.MapCmd{Store: s}, 10*time.Second)
		reg.Register(&commands.RoomCmd{Store: s}, 10*time.Second)
		reg.Register(&commands.FormatCmd{}, 5*time.Second)
		reg.Register(&commands.ResetGuessesCmd{Store: s}, 0)
		reg.Register(&commands.EmotesCmd{}, 10*time.Second)
		reg.Register(&commands.LeaderboardCmd{Store: s}, 10*time.Second)
		reg.Register(&commands.WheelGPTCmd{}, 10*time.Second)
		return reg
	}

	b := bot.New(cfg.BotUsername, cfg.BotOAuthToken, store, factory)

	// Boot sequence (mirrors the TS index.ts order)
	ctx := context.Background()

	if err := updateChannelDetails(ctx, store, tc); err != nil {
		slog.Warn("updateChannelDetails failed", "err", err)
	}

	if cfg.UpdateWebHooks {
		channels, _ := store.ListChannels(ctx)
		ids := make([]string, len(channels))
		for i, c := range channels {
			ids[i] = c.ID
		}
		if err := tc.SyncWebhooks(ctx, ids, cfg.EventSubCallbackURL, cfg.EventSubSecret); err != nil {
			slog.Warn("syncWebhooks failed", "err", err)
		}
	}

	if err := updateStreamStatus(ctx, store, tc); err != nil {
		slog.Warn("updateStreamStatus failed", "err", err)
	}

	// Start bot in background
	go func() {
		if err := b.Start(ctx); err != nil {
			slog.Error("bot.Start failed", "err", err)
			os.Exit(1)
		}
	}()

	// HTTP server
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: api.NewServer(cfg, store, tc, b).Handler(),
	}

	go func() {
		slog.Info("listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down")

	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutCtx)
	b.Stop()
}

func updateChannelDetails(ctx context.Context, store *db.Store, tc *twitch.Client) error {
	channels, err := store.ListChannels(ctx)
	if err != nil {
		return err
	}
	for _, ch := range channels {
		u, err := tc.GetUserByID(ctx, ch.ID)
		if err != nil {
			slog.Warn("updateChannelDetails: user lookup failed", "id", ch.ID, "err", err)
			continue
		}
		_ = store.UpdateChannelDetails(ctx, db.UpdateChannelDetailsParams{
			DisplayName: u.DisplayName,
			Login:       u.Login,
			ID:          ch.ID,
		})
	}
	return nil
}

func updateStreamStatus(ctx context.Context, store *db.Store, tc *twitch.Client) error {
	channels, err := store.ListChannels(ctx)
	if err != nil {
		return err
	}
	for _, ch := range channels {
		stream, err := tc.GetStream(ctx, ch.ID)
		if err != nil {
			slog.Warn("updateStreamStatus: lookup failed", "id", ch.ID, "err", err)
			continue
		}
		live := int64(0)
		if stream.Type == "live" {
			live = 1
		}
		_ = store.SetChannelLive(ctx, db.SetChannelLiveParams{IsLive: live, ID: ch.ID})
	}
	return nil
}
