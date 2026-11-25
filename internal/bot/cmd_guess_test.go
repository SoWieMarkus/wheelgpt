package bot_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/channel"
)

func TestGuessRound(t *testing.T) {
	h := newHarness(t)
	if got := h.say("Alice", "!guess 1:00.000", bot.Everyone); got != "" {
		t.Errorf("guess replied %q", got)
	}
	h.say("Bob", "!g 59.900", bot.Everyone)
	if got := h.say("Carol", "!guess nope", bot.Everyone); !strings.Contains(got, "hh:mm:ss.xxx") {
		t.Errorf("bad format reply = %q", got)
	}
	if got := h.say("Alice", "!mg", bot.Everyone); got != "@Alice 1:00.000" {
		t.Errorf("myguess = %q", got)
	}
	if got := h.say("Alice", "!result 59.950", bot.Everyone); got != "" {
		t.Errorf("viewer result replied %q", got)
	}
	if got := h.say("Mod", "!result 59.900", bot.Moderator); !strings.Contains(got, "@Bob the GIGACHAD guessed it correctly!") {
		t.Errorf("result = %q", got)
	}
	if got := h.say("Alice", "!mg", bot.Everyone); !strings.Contains(got, "can't find any guess") {
		t.Errorf("after round = %q", got)
	}
}

func TestMinimumGuessAge(t *testing.T) {
	h := newHarness(t)
	ten := 30.0
	if _, err := h.ch.UpdateSettings(context.Background(), "100", channel.SettingsPatch{GuessMinAgeSeconds: &ten}); err != nil {
		t.Fatal(err)
	}
	h.say("Early", "!guess 50.000", bot.Everyone)
	h.now = h.now.Add(time.Minute)
	h.say("Late", "!guess 40.000", bot.Everyone) // placed right before the PB
	got := h.say("Mod", "!result 40.000", bot.Moderator)
	if !strings.Contains(got, "@Early") || strings.Contains(got, "@Late") {
		t.Errorf("late guess must not count: %q", got)
	}
}
