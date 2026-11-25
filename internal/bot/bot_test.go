package bot_test

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/guess"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/leaderboard"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store/storetest"
)

type sender struct {
	mu   sync.Mutex
	sent []string
}

func (s *sender) Send(_ context.Context, _, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sent = append(s.sent, text)
	return nil
}

func (s *sender) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.sent)
}

func (s *sender) last() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sent[len(s.sent)-1]
}

type noSubs struct{}

func (noSubs) Connect(context.Context, string) (string, error) { return channel.ChatUnknown, nil }
func (noSubs) DisconnectChat(context.Context, string) error    { return nil }
func (noSubs) Remove(context.Context, string) error            { return nil }

type noRecorder struct{}

func (noRecorder) Command(string) {}
func (noRecorder) Guess()         {}

type harness struct {
	t      *testing.T
	st     *store.Store
	bot    *bot.Bot
	sender *sender
	now    time.Time
	ch     *channel.Service
}

func newHarness(t *testing.T) *harness {
	t.Helper()
	st := storetest.New(t)
	storetest.Channel(t, st, "100", "streamer")
	h := &harness{t: t, st: st, sender: &sender{}, now: time.Date(2026, 11, 20, 18, 0, 0, 0, time.UTC)}
	clock := func() time.Time { return h.now }

	h.ch = channel.NewService(channel.NewRepository(st), noSubs{})
	in := ingame.NewService(ingame.NewRepository(st), nil)
	guesses := guess.NewService(guess.NewRepository(st), in, noRecorder{})
	guesses.SetClock(clock)
	boards := leaderboard.NewService(leaderboard.NewRepository(st))
	boards.SetClock(clock)
	h.bot = bot.New(bot.Deps{
		Channels: h.ch, Guesses: guesses, Ingame: in, Leaderboard: boards,
		Sender: h.sender, Recorder: noRecorder{}, PublicURL: "https://wheelgpt.dev",
	})
	return h
}

func (h *harness) say(user, text string, access bot.Access) string {
	h.t.Helper()
	c, err := h.ch.Get(context.Background(), "100")
	if err != nil {
		h.t.Fatal(err)
	}
	before := h.sender.count()
	if err := h.bot.Handle(context.Background(), c, bot.Message{
		BroadcasterID: "100", ChatterID: "id-" + user, ChatterName: user, Text: text, Access: access,
	}); err != nil {
		h.t.Fatalf("%s: %v", text, err)
	}
	if h.sender.count() == before {
		return ""
	}
	return h.sender.last()
}

func TestParse(t *testing.T) {
	key, args, ok := bot.Parse("!Guess 1:02.345 \U000E0000")
	if !ok || key != "guess" || len(args) != 1 || args[0] != "1:02.345" {
		t.Fatalf("got %q %q %v", key, args, ok)
	}
	if _, _, ok := bot.Parse("hello !guess"); ok {
		t.Error("not a command")
	}
	if _, _, ok := bot.Parse("!"); ok {
		t.Error("empty command")
	}
}

func TestAccessFromBadges(t *testing.T) {
	if got := bot.AccessFromBadges([]string{"subscriber", "moderator"}); got != bot.Moderator {
		t.Errorf("got %v", got)
	}
	if got := bot.AccessFromBadges(nil); got != bot.Everyone {
		t.Errorf("got %v", got)
	}
}

func TestCooldown(t *testing.T) {
	h := newHarness(t)
	if got := h.say("A", "!map", bot.Everyone); !strings.Contains(got, "No map") {
		t.Fatalf("map = %q", got)
	}
	if got := h.say("B", "!map", bot.Everyone); got != "" {
		t.Errorf("cooldown ignored: %q", got)
	}
	h.now = h.now.Add(11 * time.Second)
	if got := h.say("B", "!map", bot.Everyone); got == "" {
		t.Error("cooldown did not expire")
	}
}

// Stopping the bot sends pending announcements right away, so no result is
// lost on a restart.
func TestAnnounceFlushOnStop(t *testing.T) {
	h := newHarness(t)
	h.bot.Announce("100", "hello", time.Hour)
	h.bot.Stop()
	if h.sender.last() != "hello" {
		t.Errorf("announcement not flushed: %q", h.sender.last())
	}
}

func TestTruncate(t *testing.T) {
	if got := []rune(bot.Truncate(strings.Repeat("ä", 600))); len(got) != bot.MaxMessageLength {
		t.Errorf("len = %d", len(got))
	}
}
