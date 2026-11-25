package bot

import (
	"context"
	"fmt"

	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

func myGuessCommand() *Command {
	return &Command{
		Name: "myguess", Aliases: []string{"mg"}, Access: Everyone,
		Usage:       "!myguess",
		Description: "Show your current guess.",
		Run:         runMyGuess,
	}
}

func runMyGuess(ctx context.Context, b *Bot, inv Invocation) (string, error) {
	g, ok, err := b.Guesses.Get(ctx, inv.Channel.ID, inv.User.ChatterID)
	if err != nil {
		return "", err
	}
	if !ok {
		return fmt.Sprintf("%s I can't find any guess from you %s", inv.mention(), tm.EmoteYek.Name), nil
	}
	return inv.mention() + " " + g.Time.String(), nil
}
