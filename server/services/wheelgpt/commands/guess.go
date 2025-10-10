package commands

import (
	_ "embed"

	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch"
)

type GuessCommand struct {
	BaseCommand
}

func (c *GuessCommand) Execute(user twitch.User, args []string) *string {
	return nil
}
