package db

import (
	"testing"

	"github.com/SoWieMarkus/wheelgpt/pkg/config"
	testdb "github.com/SoWieMarkus/wheelgpt/pkg/database/tests"
)

type MockTable struct {
	ID   int    `db:"id,primarykey"`
	Name string `db:"name"`
}

func (m MockTable) TableName() string {
	return "mock_table"
}

func (m MockTable) Indexes() []Index {
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
	actual := getDatabaseURL(conf)
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestDB_CreateTable(t *testing.T) {
	dbEnv := testdb.SQLite(t)
	db := Database{DbMap: dbEnv.DbMap}
	defer db.Close()
	defer dbEnv.Close()

	table := db.AddTable(MockTable{})
	err := db.CreateTable(table)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !db.TableExists(MockTable{}) {
		t.Fatal("expected table to exist")
	}
}

func TestDB_AddTable(t *testing.T) {
	dbEnv := testdb.SQLite(t)
	db := Database{DbMap: dbEnv.DbMap}
	defer db.Close()
	defer dbEnv.Close()

	table := db.AddTable(MockTable{})
	if table == nil {
		t.Fatal("expected table to be added")
	}
}
