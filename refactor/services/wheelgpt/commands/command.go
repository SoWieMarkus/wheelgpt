package commands

import (
	db "github.com/SoWieMarkus/wheelgpt/pkg/database"
)

type Command interface {
	Execute(user string, args []string) *string
	Initialize() error
}

type BaseCommand struct {
	// Database connection.
	DB db.Database
}
