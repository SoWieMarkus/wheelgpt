package bot

import (
	"context"

	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

func guessCommand() *Command {
	return &Command{
		Name: "guess", Aliases: []string{"g"}, Access: Everyone,
		Usage:       "!guess 1:23.456",
		Description: "Guess the streamer's next personal best on the current map.",
		Run:         runGuess,
	}
}

func runGuess(ctx context.Context, b *Bot, inv Invocation) (string, error) {
	if len(inv.Args) == 0 {
		return inv.mention() + " " + tm.TimeFormatHelp, nil
	}
	t, ok := tm.ParseTime(inv.Args[0])
	if !ok {
		return inv.mention() + " " + tm.TimeFormatHelp, nil
	}
	return "", b.Guesses.Place(ctx, inv.Channel.ID, inv.User.ChatterID, inv.User.ChatterName, t)
}
