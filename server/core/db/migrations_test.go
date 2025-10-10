package db

import (
	"reflect"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/testlib"
)

type migrationMockTable struct {
	ID int `db:"id,primarykey"`
}

func (m migrationMockTable) TableName() string {
	return "migration_mock_table"
}

func (m migrationMockTable) Indexes() []Index {
	return nil
}

type otherMigrationMockTable struct {
	ID int `db:"id,primarykey"`
}

func (m otherMigrationMockTable) TableName() string {
	return "other_migration_mock_table"
}

func (m otherMigrationMockTable) Indexes() []Index {
	return nil
}

func TestMigrater_NewMigrater(t *testing.T) {
	dbEnv := testlib.SQLite(t)
	db := Database{DbMap: dbEnv.DbMap}
	defer db.Close()
	defer dbEnv.Close()

	migrater := NewMigrater(db)

	// Add all new migration files here to ensure they are loaded
	expectedMigrations := map[string]bool{
		"000_initial.sql": true,
	}

	if len(migrater.Migrations) != len(expectedMigrations) {
		t.Fatalf("expected %d migrations, got %d", len(expectedMigrations), len(migrater.Migrations))
	}

	for fileName := range expectedMigrations {
		if _, exists := migrater.Migrations[fileName]; !exists {
			t.Errorf("expected migration file %s to be loaded", fileName)
		}
	}
}

func TestMigrater_CreateTableIfNotExists(t *testing.T) {
	tests := []struct {
		Name            string
		MockData        []any
		ExistingTables  []Table
		TablesToCreate  []Table
		ExpectedCreated []Table
	}{
		{
			Name:           "Create new tables when none exist",
			ExistingTables: []Table{},
			TablesToCreate: []Table{
				migrationMockTable{},
				otherMigrationMockTable{},
			},
			ExpectedCreated: []Table{
				migrationMockTable{},
				otherMigrationMockTable{},
			},
			MockData: []any{},
		},
		{
			Name: "Skip existing tables",
			MockData: []any{
				&migrationMockTable{ID: 1},
				&migrationMockTable{ID: 2},
			},
			ExistingTables: []Table{
				migrationMockTable{},
			},
			TablesToCreate: []Table{
				migrationMockTable{},
				otherMigrationMockTable{},
			},
			ExpectedCreated: []Table{
				migrationMockTable{},
				otherMigrationMockTable{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			dbEnv := testlib.SQLite(t)
			db := Database{DbMap: dbEnv.DbMap}
			defer db.Close()
			defer dbEnv.Close()

			// Create existing tables and insert mock data
			for _, table := range tt.ExistingTables {
				if err := db.CreateTable(db.AddTable(table)); err != nil {
					t.Fatalf("failed to create existing table: %v", err)
				}
			}

			if err := db.Insert(tt.MockData...); err != nil {
				t.Fatalf("failed to insert mock data: %v", err)
			}

			migrater := NewMigrater(db)

			if err := migrater.CreateTableIfNotExists(tt.TablesToCreate...); err != nil {
				t.Fatalf("expected no error during table creation, got %v", err)
			}

			// Check that all expected tables exist
			for _, expectedTable := range tt.ExpectedCreated {
				if !db.TableExists(expectedTable) {
					t.Errorf("expected table %s to exist", expectedTable.TableName())
				}
			}

			// Check that existing data is still present
			for _, data := range tt.MockData {
				var count int
				if err := db.DbMap.SelectOne(&count, "SELECT COUNT(*) FROM "+data.(Table).TableName()); err != nil {
					t.Fatalf("failed to query table %s: %v", data.(Table).TableName(), err)
				}
				if count == 0 {
					t.Errorf("expected data in table %s to still exist", data.(Table).TableName())
				}
			}
		})
	}
}

func TestMigrater_Migrate(t *testing.T) {
	dbEnv := testlib.SQLite(t)
	db := Database{DbMap: dbEnv.DbMap}
	defer db.Close()
	defer dbEnv.Close()

	// Create a mock table to apply migrations to
	if err := db.CreateTable(db.AddTable(migrationMockTable{})); err != nil {
		t.Fatalf("failed to create migrations table: %v", err)
	}

	mockData := &migrationMockTable{ID: 1}
	if err := db.Insert(mockData); err != nil {
		t.Fatalf("failed to insert mock data: %v", err)
	}

	// Ensure migrations table does not exist before the first Migrate() call
	if db.TableExists(Migration{}) {
		t.Fatalf("migrations table should not exist yet")
	}

	migrater := &Migrater{
		DB: db,
		Migrations: map[string]string{
			"test.sql": `
                ALTER TABLE migration_mock_table ADD COLUMN name TEXT DEFAULT 'default_name';
            `,
		},
	}

	// Run migration
	if err := migrater.Migrate(); err != nil {
		t.Fatalf("expected no error during migration, got %v", err)
	}

	// Check that migrations table was created
	if !db.TableExists(Migration{}) {
		t.Fatalf("migrations table should exist after migration")
	}

	// Check that the new migration entry was added
	var appliedMigrations []Migration
	if _, err := db.Select(&appliedMigrations, "SELECT * FROM migrations"); err != nil {
		t.Fatalf("failed to query migrations table: %v", err)
	}
	if len(appliedMigrations) != 1 || appliedMigrations[0].File != "test.sql" {
		t.Fatalf("expected migration 'test.sql' to be recorded, got %v", appliedMigrations)
	}

	// Check that the schema change was applied
	type newMigrationMockTable struct {
		ID   int    `db:"id,primarykey"`
		Name string `db:"name"`
	}
	var results []newMigrationMockTable
	if _, err := db.Select(&results, "SELECT * FROM migration_mock_table"); err != nil {
		t.Fatalf("failed to query migration_mock_table: %v", err)
	}

	// Check that table was altered and data is intact
	expected := []newMigrationMockTable{
		{ID: 1, Name: "default_name"},
	}
	// Check that the deeop equal comparison works as expected
	notExpected := []newMigrationMockTable{
		{ID: 1, Name: "other_name"},
	}

	// Check if results match expected
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("expected data %v, got %v", expected, results)
	}

	if reflect.DeepEqual(results, notExpected) {
		t.Fatalf("did not expect data %v", notExpected)
	}
}
