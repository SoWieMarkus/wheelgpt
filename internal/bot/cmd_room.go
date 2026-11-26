package bot

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
)

func roomCommand() *Command {
	return &Command{
		Name: "room", Access: Everyone, Cooldown: 10 * time.Second,
		Usage:       "!room",
		Description: "Show the server room the streamer is playing in.",
		Run:         runRoom,
	}
}

func runRoom(ctx context.Context, b *Bot, inv Invocation) (string, error) {
	room, err := b.Ingame.Room(ctx, inv.Channel.ID)
	if errors.Is(err, ingame.ErrRoomNotFound) {
		return inv.mention() + " I am currently not in a room.", nil
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s [%d/%d]", inv.mention(), room.Name, room.NumberOfPlayers, room.MaxPlayers), nil
}
