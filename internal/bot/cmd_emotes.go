package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

func emotesCommand() *Command {
	return &Command{
		Name: "wgpt-emotes", Access: Moderator, Cooldown: 10 * time.Second,
		Usage:       "!wgpt-emotes",
		Description: "List the 7TV emotes WheelGPT uses. Add them to your channel.",
		Run:         runEmotes,
	}
}

func runEmotes(_ context.Context, _ *Bot, inv Invocation) (string, error) {
	names := make([]string, len(tm.Emotes))
	for i, e := range tm.Emotes {
		names[i] = e.Name
	}
	return fmt.Sprintf("%s %s %s", inv.mention(), strings.Join(names, " "), tm.EmoteYek.Name), nil
}
