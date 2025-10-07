package commands

import (
	"testing"

	db "github.com/SoWieMarkus/wheelgpt/pkg/database"
	"github.com/SoWieMarkus/wheelgpt/pkg/database/models"
	testdb "github.com/SoWieMarkus/wheelgpt/pkg/database/tests"
)

func TestGuessCommand_Execute(t *testing.T) {
	tests := []struct {
		Name          string
		ExistingGuess *models.Guess
		ExpectedGuess *models.Guess
		User          string
		Args          []string
	}{}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sqlite := testdb.SQLite(t)
			testDB := db.Database{DbMap: sqlite.DbMap}
			defer testDB.Close()
			defer sqlite.Close()
		})
	}
}
