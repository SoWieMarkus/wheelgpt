package commands

import (
	_ "embed"
)

type MyGuessCommand struct {
	BaseCommand
}

//go:embed my_guess.sql
var myGuessSQL string

func (c *MyGuessCommand) Execute(user string, args []string) *string {
	return nil
}
