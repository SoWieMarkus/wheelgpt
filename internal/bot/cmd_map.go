package bot

import (
	"context"
	"time"
)

func mapCommand() *Command {
	return &Command{
		Name: "map", Access: Everyone, Cooldown: 10 * time.Second,
		Usage:       "!map",
		Description: "Show the current map with its medals and links.",
		Run:         runMap,
	}
}

func runMap(ctx context.Context, b *Bot, inv Invocation) (string, error) {
	m, err := b.Ingame.CurrentMap(ctx, inv.Channel.ID)
	if err != nil {
		return "", err
	}
	if m == nil {
		return inv.mention() + " No map is currently set for this channel.", nil
	}
	return inv.mention() + " " + m.String(), nil
}
