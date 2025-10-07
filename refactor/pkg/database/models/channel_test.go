package models

import (
	"testing"

	db "github.com/SoWieMarkus/wheelgpt/pkg/database"
	testdb "github.com/SoWieMarkus/wheelgpt/pkg/database/tests"
)

func TestChannel_NoDuplicatePrimaryKey(t *testing.T) {
	sqlite := testdb.SQLite(t)
	testDB := db.Database{DbMap: sqlite.DbMap}
	defer testDB.Close()
	defer sqlite.Close()

	// Create dependency tables
	if err := testDB.CreateTable(
		testDB.AddTable(Channel{}),
	); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	channel := &Channel{
		ID:    "test_channel",
		Login: "test_user",
		Seed:  "test_seed",
	}
	if err := testDB.Insert(channel); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	duplicateChannel := &Channel{
		ID:    "test_channel",
		Login: "test_user_2",
		Seed:  "test_seed_2",
	}
	if err := testDB.Insert(duplicateChannel); err == nil {
		t.Fatalf("expected error due to duplicate primary key, got nil")
	}
}
