package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

const leaderboardLimit = 5

type LeaderboardCmd struct{ Store *db.Store }

func (c *LeaderboardCmd) Names() []string    { return []string{"wgpt-leaderboard", "lb"} }
func (c *LeaderboardCmd) Level() bot.AccessLevel { return bot.AccessUser }

func (c *LeaderboardCmd) Execute(ctx *bot.CmdContext) (string, error) {
	bctx := context.Background()

	if len(ctx.Args) > 0 {
		name := strings.TrimPrefix(ctx.Args[0], "@")
		entry, err := c.Store.GetLeaderboardByName(bctx, db.GetLeaderboardByNameParams{
			ChannelID: ctx.ChannelID,
			Lower:     strings.ToLower(name),
		})
		if err != nil {
			return fmt.Sprintf("%s No leaderboard entry found for %s.",
				bot.MentionUser(ctx.User.DisplayName), name), nil
		}
		pts := "pts"
		if entry.Points == 1 {
			pts = "pt"
		}
		return fmt.Sprintf("%s is currently in position %d with %d %s.",
			bot.MentionUser(entry.DisplayName), entry.Position, entry.Points, pts), nil
	}

	rows, err := c.Store.GetLeaderboard(bctx, db.GetLeaderboardParams{
		ChannelID: ctx.ChannelID,
		Limit:     leaderboardLimit,
		Offset:    0,
	})
	if err != nil || len(rows) == 0 {
		return "No guesses have been made yet.", nil
	}

	entries := make([]string, len(rows))
	for i, r := range rows {
		entries[i] = fmt.Sprintf("%d. %s - %d", r.Position, r.DisplayName, r.Points)
	}
	url := fmt.Sprintf("https://wheelgpt.dev/c/%s", ctx.ChannelID)
	return fmt.Sprintf("Top %d Guessers 🏆 Perfect guess = 5pts, Closest guess = 1pt: %s | Full leaderboard: %s",
		leaderboardLimit, strings.Join(entries, " | "), url), nil
}
