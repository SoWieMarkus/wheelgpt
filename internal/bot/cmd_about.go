package bot

import (
	"context"
	"fmt"
	"time"

	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

func aboutCommand() *Command {
	return &Command{
		Name: "wheelgpt", Aliases: []string{"wgpt"}, Access: Everyone, Cooldown: 10 * time.Second,
		Usage:       "!wheelgpt",
		Description: "What is this bot?",
		Run:         runAbout,
	}
}

func runAbout(_ context.Context, b *Bot, _ Invocation) (string, error) {
	return fmt.Sprintf("%s I am a Trackmania Twitch Bot! You can ask me about the current map, the current room "+
		"and guess the next best time. More details: %s | ⭐ Star me on Github: https://github.com/SoWieMarkus/wheelgpt",
		tm.EmoteHeyGuys.Name, b.PublicURL), nil
}
