package commands

import _ "embed"

type GuessCommand struct {
	BaseCommand
}

func (c *GuessCommand) Execute(user string, args []string) *string {
	return nil
}
