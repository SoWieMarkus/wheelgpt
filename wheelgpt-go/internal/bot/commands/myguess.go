package commands

import (
	"context"
	"fmt"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

type MyGuessCmd struct{ Store *db.Store }

func (c *MyGuessCmd) Names() []string    { return []string{"myguess", "mg"} }
func (c *MyGuessCmd) Level() bot.AccessLevel { return bot.AccessUser }

func (c *MyGuessCmd) Execute(ctx *bot.CmdContext) (string, error) {
	g, err := c.Store.GetGuess(context.Background(), db.GetGuessParams{
		ChannelID: ctx.ChannelID,
		UserID:    ctx.User.ID,
	})
	if err != nil {
		return fmt.Sprintf("%s you haven't made a guess yet.", bot.MentionUser(ctx.User.DisplayName)), nil
	}
	t := bot.TMTime(g.Time)
	return fmt.Sprintf("%s your current guess is %s.", bot.MentionUser(ctx.User.DisplayName), t), nil
}
