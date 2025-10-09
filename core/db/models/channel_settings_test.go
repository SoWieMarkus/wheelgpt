package models

import (
	"encoding/json"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/testlib"
)

func TestChannelSettings_NoDuplicatePrimaryKey(t *testing.T) {
	sqlite := testlib.SQLite(t)
	testDB := db.Database{DbMap: sqlite.DbMap}
	defer testDB.Close()
	defer sqlite.Close()

	// Create dependency tables
	if err := testDB.CreateTable(
		testDB.AddTable(ChannelSettings{}),
	); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	settings := &ChannelSettings{
		ID:                         "test_channel",
		GuessDelayTime:             10,
		BotActiveWhenStreamOffline: true,
		Public:                     true,
	}
	if err := testDB.Insert(settings); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	duplicateSettings := &ChannelSettings{
		ID:                         "test_channel",
		GuessDelayTime:             15,
		BotActiveWhenStreamOffline: false,
		Public:                     false,
	}
	if err := testDB.Insert(duplicateSettings); err == nil {
		t.Fatalf("expected error due to duplicate primary key, got nil")
	}
}

func TestChannelSettings_JsonMarshalling(t *testing.T) {
	settings := &ChannelSettings{
		ID:                         "test_channel",
		GuessDelayTime:             10,
		BotActiveWhenStreamOffline: true,
		Public:                     true,
	}

	jsonData, err := json.Marshal(settings)
	if err != nil {
		t.Fatalf("expected no error during JSON marshalling, got %v", err)
	}

	expectedJson := `{"id":"test_channel","guessDelayTime":10,"botActiveWhenStreamOffline":true,"public":true}`
	if string(jsonData) != expectedJson {
		t.Fatalf("expected JSON %q, got %q", expectedJson, jsonData)
	}
}
