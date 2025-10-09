package models

import (
	"encoding/json"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/testlib"
)

func TestGuess_NoDuplicatePrimaryKey(t *testing.T) {
	sqlite := testlib.SQLite(t)
	testDB := db.Database{DbMap: sqlite.DbMap}
	defer testDB.Close()
	defer sqlite.Close()

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

func TestGuess_JsonMarshalling(t *testing.T) {
	guess := &Guess{
		ChannelID:   "test_channel",
		UserID:      "test_user",
		Time:        1234567890,
		DisplayName: "Test User",
	}

	jsonData, err := json.Marshal(guess)
	if err != nil {
		t.Fatalf("expected no error during JSON marshalling, got %v", err)
	}

	expectedJson := `{"channelId":"test_channel","userId":"test_user","time":1234567890,"displayName":"Test User"}`
	if string(jsonData) != expectedJson {
		t.Fatalf("expected JSON %q, got %q", expectedJson, jsonData)
	}
}
