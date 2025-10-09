package commands

import (
	"reflect"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/core/db/models"
	"github.com/SoWieMarkus/wheelgpt/testlib"
)

func TestMyGuess_Query(t *testing.T) {
	tests := []struct {
		Name           string
		UserId         string
		ChannelId      string
		Guesses        []any
		ExpectedGuess  models.Guess
		ExpectNoResult bool
	}{
		{
			Name:           "No guesses",
			UserId:         "user1",
			ChannelId:      "channel1",
			Guesses:        []any{},
			ExpectedGuess:  models.Guess{},
			ExpectNoResult: true,
		},
		{
			Name:      "Single guess should",
			UserId:    "user1",
			ChannelId: "channel1",
			Guesses: []any{
				&models.Guess{
					DisplayName: "user1",
					UserID:      "user1",
					ChannelID:   "channel1",
					Time:        1000,
				},
			},
			ExpectedGuess: models.Guess{
				DisplayName: "user1",
				UserID:      "user1",
				ChannelID:   "channel1",
				Time:        1000,
			},
			ExpectNoResult: false,
		},
		{
			Name:      "No guesses on this channel",
			UserId:    "user1",
			ChannelId: "channel1",
			Guesses: []any{
				&models.Guess{
					DisplayName: "user1",
					UserID:      "user1",
					ChannelID:   "channel2",
					Time:        1000,
				},
			},
			ExpectedGuess:  models.Guess{},
			ExpectNoResult: true,
		},
		{
			Name:      "No guesses from user on this channel",
			UserId:    "user1",
			ChannelId: "channel1",
			Guesses: []any{
				&models.Guess{
					DisplayName: "user2",
					UserID:      "user2",
					ChannelID:   "channel1",
					Time:        1000,
				},
			},
			ExpectedGuess:  models.Guess{},
			ExpectNoResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sqlite := testlib.SQLite(t)
			testDB := db.Database{DbMap: sqlite.DbMap}
			defer testDB.Close()
			defer sqlite.Close()

			if err := testDB.CreateTable(
				testDB.AddTable(models.Guess{}),
			); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if err := testDB.Insert(tt.Guesses...); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			var guess models.Guess
			err := testDB.SelectOne(&guess, myGuessSQL, tt.UserId, tt.ChannelId)
			if (err != nil) != tt.ExpectNoResult {
				t.Errorf("Database query error: %v", err)
			}

			if !reflect.DeepEqual(guess, tt.ExpectedGuess) {
				t.Errorf("Expected: %+v, got: %+v", tt.ExpectedGuess, guess)
			}
		})
	}
}
