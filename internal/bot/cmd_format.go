package bot

import (
	"context"
	"time"

	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

func formatCommand() *Command {
	return &Command{
		Name: "format", Access: Everyone, Cooldown: 5 * time.Second,
		Usage:       "!format @viewer",
		Description: "Explain the time format for guesses.",
		Run:         runFormat,
	}
}

func runFormat(_ context.Context, _ *Bot, inv Invocation) (string, error) {
	if len(inv.Args) > 0 {
		return tm.Mention(inv.Args[0]) + " " + tm.TimeFormatHelp, nil
	}
	return inv.mention() + " " + tm.TimeFormatHelp, nil
}
