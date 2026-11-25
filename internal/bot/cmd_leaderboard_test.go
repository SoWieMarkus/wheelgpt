package bot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
)

func TestLeaderboardCommand(t *testing.T) {
	h := newHarness(t)
	if got := h.say("A", "!lb", bot.Everyone); got != "No points this month yet. Leaderboard: https://wheelgpt.dev/c/streamer" {
		t.Errorf("empty = %q", got)
	}
	h.say("Alice", "!guess 1.000", bot.Everyone)
	h.say("Mod", "!result 1.100", bot.Moderator)
	h.now = h.now.Add(11 * time.Second)
	if got := h.say("A", "!lb", bot.Everyone); !strings.Contains(got, "Top 5 Guessers of November") || !strings.Contains(got, "1. Alice - 1") {
		t.Errorf("monthly = %q", got)
	}
	h.now = h.now.Add(11 * time.Second)
	if got := h.say("A", "!lb @alice", bot.Everyone); got != "@Alice is #1 this month (1 pt) and #1 all-time (1 pt)." {
		t.Errorf("standing = %q", got)
	}
	h.now = h.now.Add(11 * time.Second)
	if got := h.say("A", "!lb all", bot.Everyone); !strings.Contains(got, "of all time") {
		t.Errorf("all-time = %q", got)
	}
}
