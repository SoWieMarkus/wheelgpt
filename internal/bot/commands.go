package bot

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	tm "github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

// Invocation is one execution of a command.
type Invocation struct {
	Channel channel.Channel
	User    Message
	Args    []string
}

func (i Invocation) mention() string { return tm.Mention(i.User.ChatterName) }

// Command is a chat command.
type Command struct {
	Name        string
	Aliases     []string
	Access      Access
	Cooldown    time.Duration
	Usage       string
	Description string
	Run         func(ctx context.Context, b *Bot, inv Invocation) (string, error)
}

// Matches reports whether key names this command.
func (c *Command) Matches(key string) bool {
	return c.Name == key || slices.Contains(c.Aliases, key)
}

// defaultCommands is the command registry: every chat command the bot knows.
// Each command's definition and Run function live in its own cmd_<name>.go file.
func defaultCommands() []*Command {
	return []*Command{
		guessCommand(),
		myGuessCommand(),
		resultCommand(),
		mapCommand(),
		emotesCommand(),
		formatCommand(),
		roomCommand(),
		resetGuessesCommand(),
		aboutCommand(),
		leaderboardCommand(),
	}
}

func pts(n int64) string {
	if n == 1 {
		return "1 pt"
	}
	return fmt.Sprintf("%d pts", n)
}
