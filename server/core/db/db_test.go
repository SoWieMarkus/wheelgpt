package db

import (
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/testlib"
)

type mockTable struct {
	ID   int    `db:"id,primarykey"`
	Name string `db:"name"`
}

func (m mockTable) TableName() string {
	return "mock_table"
}

func (m mockTable) Indexes() []Index {
	return nil
}

func TestHelper_getDatabaseUrl(t *testing.T) {
	conf := config.DatabaseConfig{
		User:         "testuser",
		Password:     "testpass",
		HostName:     "localhost",
		Port:         5432,
		DatabaseName: "testdb",
	}
	expected := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	actual := getDatabaseURL(&conf)
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestDB_CreateTable(t *testing.T) {
	dbEnv := testlib.SQLite(t)
	db := Database{DbMap: dbEnv.DbMap}
	defer db.Close()
	defer dbEnv.Close()

	table := db.AddTable(mockTable{})
	err := db.CreateTable(table)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !db.TableExists(mockTable{}) {
		t.Fatal("expected table to exist")
	}
}

func TestDB_AddTable(t *testing.T) {
	dbEnv := testlib.SQLite(t)
	db := Database{DbMap: dbEnv.DbMap}
	defer db.Close()
	defer dbEnv.Close()

	table := db.AddTable(mockTable{})
	if table == nil {
		t.Fatal("expected table to be added")
	}
}
