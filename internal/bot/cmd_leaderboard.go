package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/leaderboard"
	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

func leaderboardCommand() *Command {
	return &Command{
		Name: "wgpt-leaderboard", Aliases: []string{"lb"}, Access: Everyone, Cooldown: 10 * time.Second,
		Usage:       "!lb [all | name]",
		Description: "Top guessers of this month or of all time, or one viewer's position.",
		Run:         runLeaderboard,
	}
}

// LeaderboardSize is the number of entries !lb shows.
const LeaderboardSize = 5

// allTimeArgs select the all-time board in !lb.
var allTimeArgs = map[string]bool{"all": true, "all-time": true, "alltime": true} //nolint:misspell // chat keyword

func runLeaderboard(ctx context.Context, b *Bot, inv Invocation) (string, error) {
	url := b.PublicURL + "/c/" + inv.Channel.Login
	month := b.Leaderboard.CurrentMonth()

	if len(inv.Args) > 0 && !allTimeArgs[strings.ToLower(inv.Args[0])] {
		name := strings.TrimPrefix(inv.Args[0], "@")
		monthly, inMonth, err := b.Leaderboard.Find(ctx, inv.Channel.ID, month, "", name)
		if err != nil {
			return "", err
		}
		allTime, inAllTime, err := b.Leaderboard.Find(ctx, inv.Channel.ID, leaderboard.AllTimePeriod(), "", name)
		if err != nil {
			return "", err
		}
		if !inMonth && !inAllTime {
			return fmt.Sprintf("%s No leaderboard entry found for %s.", inv.mention(), name), nil
		}
		var parts []string
		display := name
		if inMonth {
			display = monthly.DisplayName
			parts = append(parts, fmt.Sprintf("#%d this month (%s)", monthly.Position, pts(monthly.Points)))
		} else {
			parts = append(parts, "no points this month")
		}
		if inAllTime {
			display = allTime.DisplayName
			parts = append(parts, fmt.Sprintf("#%d all-time (%s)", allTime.Position, pts(allTime.Points)))
		}
		return fmt.Sprintf("%s is %s.", tm.Mention(display), strings.Join(parts, " and ")), nil
	}

	period, title, empty := month, "of "+month.From.Format("January"), "No points this month yet."
	if len(inv.Args) > 0 {
		period, title, empty = leaderboard.AllTimePeriod(), "of all time", "No guesses have been made yet."
	}
	entries, _, err := b.Leaderboard.Board(ctx, inv.Channel.ID, period, "", LeaderboardSize, 0)
	if err != nil {
		return "", err
	}
	if len(entries) == 0 {
		return fmt.Sprintf("%s Leaderboard: %s", empty, url), nil
	}
	lines := make([]string, len(entries))
	for i, e := range entries {
		lines[i] = fmt.Sprintf("%d. %s - %d", e.Position, e.DisplayName, e.Points)
	}
	return fmt.Sprintf("Top %d Guessers %s 🏆 Perfect guess = %dpts, Closest guess = %dpt: %s | Full leaderboard: %s",
		LeaderboardSize, title, tm.PointsPerfectGuess, tm.PointsClosestGuess, strings.Join(lines, " | "), url), nil
}
