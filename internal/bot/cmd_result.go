package bot

import (
	"context"

	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

func resultCommand() *Command {
	return &Command{
		Name: "result", Access: Moderator,
		Usage:       "!result 1:23.456",
		Description: "End the round by hand with the final time.",
		Run:         runResult,
	}
}

func runResult(ctx context.Context, b *Bot, inv Invocation) (string, error) {
	if len(inv.Args) == 0 {
		return inv.mention() + " " + tm.TimeFormatHelp, nil
	}
	t, ok := tm.ParseTime(inv.Args[0])
	if !ok {
		return inv.mention() + " smh granadyy mods are all degens. Wrong format you idiot. " + tm.TimeFormatHelp, nil
	}
	round, err := b.Guesses.EvaluatePB(ctx, inv.Channel, t)
	if err != nil {
		return "", err
	}
	return round.Message, nil
}
