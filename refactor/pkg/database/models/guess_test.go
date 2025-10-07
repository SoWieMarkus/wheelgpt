package models

import (
	"testing"

	db "github.com/SoWieMarkus/wheelgpt/pkg/database"
	testdb "github.com/SoWieMarkus/wheelgpt/pkg/database/tests"
)

func TestGuess_NoDuplicatePrimaryKey(t *testing.T) {
	dbEnv := testdb.SQLite(t)
	testDB := db.Database{DbMap: dbEnv.DbMap}
	defer testDB.Close()
	defer dbEnv.Close()

	// Create dependency tables
	if err := testDB.CreateTable(
		testDB.AddTable(Guess{}),
	); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	guess := &Guess{
		ChannelID:   "test_channel",
		UserID:      "test_user",
		Time:        1234567890,
		DisplayName: "Test User",
	}
	if err := testDB.Insert(guess); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	duplicateGuess := &Guess{
		ChannelID:   "test_channel",
		UserID:      "test_user",
		Time:        1234567891,
		DisplayName: "Test User 2",
	}
	if err := testDB.Insert(duplicateGuess); err == nil {
		t.Fatalf("expected error due to duplicate primary key, got nil")
	}
}
