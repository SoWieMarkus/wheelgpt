package commands

import (
	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch"
)

type CommandAnswer struct {
	Message string
	Reply   bool
}

type Command interface {
	Execute(user twitch.User, args []string) *CommandAnswer
	Initialize() error
}

type BaseCommand struct {
	// Database connection.
	DB db.Database
}
