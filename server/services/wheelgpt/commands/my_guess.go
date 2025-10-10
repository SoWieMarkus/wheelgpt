package commands

import (
	_ "embed"

	"github.com/SoWieMarkus/wheelgpt/core/db/models"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch"
)

// Command to get the user's last guess on the channel.
type MyGuessCommand struct {
	BaseCommand
}

//go:embed my_guess.sql
var myGuessSQL string

func (c *MyGuessCommand) Execute(user twitch.User, args []string) *CommandAnswer {
	var guess models.Guess
	err := c.DB.SelectOne(&guess, myGuessSQL, user.ID)
	if err != nil {
		return &CommandAnswer{
			Message: "I can't find any guess from you.",
			Reply:   true,
		}
	}
	return &CommandAnswer{
		Message: "TODO",
		Reply:   true,
	}
}
