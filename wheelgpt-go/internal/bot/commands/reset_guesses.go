package commands

import (
	"context"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

type ResetGuessesCmd struct{ Store *db.Store }

func (c *ResetGuessesCmd) Names() []string    { return []string{"resetguesses", "rg"} }
func (c *ResetGuessesCmd) Level() bot.AccessLevel { return bot.AccessMod }

func (c *ResetGuessesCmd) Execute(ctx *bot.CmdContext) (string, error) {
	if err := c.Store.DeleteGuessesByChannel(context.Background(), ctx.ChannelID); err != nil {
		return "", err
	}
	return "All guesses have been reset.", nil
}
