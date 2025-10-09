package models

import (
	"encoding/json"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/testlib"
)

func TestChannel_NoDuplicatePrimaryKey(t *testing.T) {
	sqlite := testlib.SQLite(t)
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

func TestChannel_JsonMarshalling(t *testing.T) {
	channel := &Channel{
		ID:     "test_channel",
		Login:  "test_user",
		Seed:   "test_seed",
		IsLive: true,
	}

	jsonData, err := json.Marshal(channel)
	if err != nil {
		t.Fatalf("expected no error during JSON marshalling, got %v", err)
	}

	expectedJson := `{"id":"test_channel","login":"test_user","seed":"test_seed","isLive":true}`
	if string(jsonData) != expectedJson {
		t.Fatalf("expected JSON %q, got %q", expectedJson, jsonData)
	}
}
