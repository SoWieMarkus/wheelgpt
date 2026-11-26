package leaderboard_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/guess"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/leaderboard"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store/storetest"
	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

type noRecorder struct{}

func (noRecorder) Guess() {}

func TestParsePeriod(t *testing.T) {
	p, err := leaderboard.ParsePeriod("2026-02")
	if err != nil || !p.From.Equal(time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)) || !p.To.Equal(time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("got %+v, %v", p, err)
	}
	if p, err := leaderboard.ParsePeriod("all-time"); err != nil || !p.IsAllTime() {
		t.Fatalf("all-time: %+v, %v", p, err)
	}
	for _, bad := range []string{"2026-13", "2026", "everything", ""} {
		if _, err := leaderboard.ParsePeriod(bad); !errors.Is(err, leaderboard.ErrInvalidPeriod) {
			t.Errorf("%q: %v", bad, err)
		}
	}
}

// Rounds in different calendar months land on different monthly boards,
// the all-time board sums everything.
func TestMonthlyPartition(t *testing.T) {
	ctx := context.Background()
	st := storetest.New(t)
	storetest.Channel(t, st, "100", "streamer")
	now := time.Date(2026, 10, 31, 23, 59, 0, 0, time.UTC)
	clock := func() time.Time { return now }

	guesses := guess.NewService(guess.NewRepository(st), ingame.NewService(ingame.NewRepository(st), nil), noRecorder{})
	guesses.SetClock(clock)
	boards := leaderboard.NewService(leaderboard.NewRepository(st))
	boards.SetClock(clock)
	ch, err := channel.NewRepository(st).Get(ctx, "100")
	if err != nil {
		t.Fatal(err)
	}

	round := func(placed map[string]tm.Time, pb tm.Time) {
		t.Helper()
		for user, g := range placed {
			if err := guesses.Place(ctx, "100", user, user, g); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := guesses.EvaluatePB(ctx, ch, pb); err != nil {
			t.Fatal(err)
		}
	}
	round(map[string]tm.Time{"1": 1000, "2": 2000}, 1000) // October: 1 gets 5
	now = now.Add(2 * time.Minute)                        // November
	round(map[string]tm.Time{"1": 1000, "2": 2000}, 2100) // November: 2 gets 1
	round(map[string]tm.Time{"2": 3000}, 3000)            // November: 2 gets 5

	october, _ := leaderboard.ParsePeriod("2026-10")
	november, _ := leaderboard.ParsePeriod("2026-11")
	check := func(p leaderboard.Period, wantTotal int64, wantFirst string, wantPoints int64) {
		t.Helper()
		board, total, err := boards.Board(ctx, "100", p, "", 10, 0)
		if err != nil || total != wantTotal || board[0].UserID != wantFirst || board[0].Points != wantPoints {
			t.Fatalf("%s = %+v (%d), %v", p.Key, board, total, err)
		}
	}
	check(october, 1, "1", 5)
	check(november, 1, "2", 6)
	check(leaderboard.AllTimePeriod(), 2, "2", 6)

	current, periods, err := boards.Periods(ctx, "100")
	if err != nil || current != "2026-11" || len(periods) != 3 || periods[1] != "2026-10" || periods[2] != "all-time" {
		t.Fatalf("periods = %s %v, %v", current, periods, err)
	}
}
