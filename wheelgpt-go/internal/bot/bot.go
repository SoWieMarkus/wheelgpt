package bot

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/db"
	twitchirc "github.com/gempir/go-twitch-irc/v4"
)

// CommandFactory builds a Registry for a channel. Defined by the caller (main)
// to avoid an import cycle between bot and bot/commands.
type CommandFactory func(channelID string, store *db.Store) *Registry

type Bot struct {
	client  *twitchirc.Client
	store   *db.Store
	factory CommandFactory

	mu       sync.RWMutex
	channels map[string]*Channel // keyed by login
}

func New(username, oauthToken string, store *db.Store, factory CommandFactory) *Bot {
	b := &Bot{
		store:    store,
		factory:  factory,
		channels: make(map[string]*Channel),
	}
	b.client = twitchirc.NewClient(username, oauthToken)
	b.client.OnPrivateMessage(b.onMessage)
	b.client.OnConnect(func() {
		slog.Info("bot connected to Twitch IRC")
	})
	return b
}

func (b *Bot) Start(ctx context.Context) error {
	channels, err := b.store.ListChannels(ctx)
	if err != nil {
		return err
	}
	for _, ch := range channels {
		b.registerChannel(ChannelConfig{
			ID:                   ch.ID,
			Login:                ch.Login,
			DisplayName:          ch.DisplayName,
			GuessDelayTime:       ch.GuessDelayTime,
			BotActiveWhenOffline: ch.BotActiveWhenOffline != 0,
		})
		time.Sleep(500 * time.Millisecond) // respect Twitch JOIN rate limit
	}
	return b.client.Connect()
}

func (b *Bot) Stop() {
	b.client.Disconnect()
}

func (b *Bot) registerChannel(cfg ChannelConfig) {
	reg := b.factory(cfg.ID, b.store)
	ch := &Channel{Config: cfg, registry: reg}
	b.mu.Lock()
	b.channels[cfg.Login] = ch
	b.mu.Unlock()
	b.client.Join(cfg.Login)
	slog.Info("joined channel", "login", cfg.Login)
}

func (b *Bot) Join(cfg ChannelConfig) {
	b.registerChannel(cfg)
}

func (b *Bot) Part(login string) {
	b.mu.Lock()
	delete(b.channels, login)
	b.mu.Unlock()
	b.client.Depart(login)
}

func (b *Bot) Reload(ctx context.Context, channelID string) error {
	ch, err := b.store.GetChannel(ctx, channelID)
	if err != nil {
		return err
	}
	cfg := ChannelConfig{
		ID:                   ch.ID,
		Login:                ch.Login,
		DisplayName:          ch.DisplayName,
		GuessDelayTime:       ch.GuessDelayTime,
		BotActiveWhenOffline: ch.BotActiveWhenOffline != 0,
	}
	b.mu.Lock()
	b.channels[cfg.Login] = &Channel{Config: cfg, registry: b.factory(cfg.ID, b.store)}
	b.mu.Unlock()
	return nil
}

// NotifyNewPB is called by the plugin webhook handler.
func (b *Bot) NotifyNewPB(login string, timeMS int64) {
	b.mu.RLock()
	ch, ok := b.channels[login]
	b.mu.RUnlock()
	if !ok {
		slog.Warn("NotifyNewPB: channel not found", "login", login)
		return
	}
	delay := ch.GuessDelay()
	go func() {
		time.Sleep(delay)
		t := TMTime(timeMS)
		// HandleNewPB lives in commands package; we call it via an injected func
		// to keep the bot package free of the commands import.
		ch.registry.handlePB(context.Background(), b.store, ch.Config.ID, t, func(msg string) {
			b.client.Say(login, msg)
		})
	}()
}

func (b *Bot) onMessage(msg twitchirc.PrivateMessage) {
	login := strings.TrimPrefix(msg.Channel, "#")

	b.mu.RLock()
	ch, ok := b.channels[login]
	b.mu.RUnlock()
	if !ok {
		return
	}

	user := User{
		ID:          msg.User.ID,
		Name:        msg.User.Name,
		DisplayName: msg.User.DisplayName,
		Level:       accessLevel(msg, login),
	}

	ch.registry.Dispatch(ch.Config.ID, user, msg.Message, func(reply string) {
		b.client.Say(login, reply)
	})
}

func accessLevel(msg twitchirc.PrivateMessage, channelLogin string) AccessLevel {
	if msg.User.Name == channelLogin {
		return AccessStreamer
	}
	if msg.User.Badges["moderator"] == 1 || msg.User.Badges["broadcaster"] == 1 {
		return AccessMod
	}
	if msg.User.Badges["vip"] == 1 {
		return AccessVIP
	}
	if msg.User.Badges["subscriber"] == 1 {
		return AccessSubscriber
	}
	return AccessUser
}
