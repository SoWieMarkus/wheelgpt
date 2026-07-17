package commands

import (
	"context"
	"fmt"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

type GuessCmd struct{ Store *db.Store }

func (c *GuessCmd) Names() []string    { return []string{"guess", "g"} }
func (c *GuessCmd) Level() bot.AccessLevel { return bot.AccessUser }

func (c *GuessCmd) Execute(ctx *bot.CmdContext) (string, error) {
	if len(ctx.Args) == 0 {
		return fmt.Sprintf("%s %s", bot.MentionUser(ctx.User.DisplayName), bot.ExampleFormat), nil
	}
	t, ok := bot.ParseTMTime(ctx.Args[0])
	if !ok {
		return fmt.Sprintf("%s %s", bot.MentionUser(ctx.User.DisplayName), bot.ExampleFormat), nil
	}
	err := c.Store.UpsertGuess(context.Background(), db.UpsertGuessParams{
		ChannelID:   ctx.ChannelID,
		UserID:      ctx.User.ID,
		DisplayName: ctx.User.DisplayName,
		Time:        int64(t),
	})
	return "", err
}
