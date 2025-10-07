package commands

import _ "embed"

//go:embed guess.sql
var guessSQL string

type GuessCommand struct {
	BaseCommand
}

func (c *GuessCommand) Execute(user string, args []string) *string {
	return nil
}
