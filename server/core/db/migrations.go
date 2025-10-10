package db

import (
	"embed"
	"log/slog"
)

// Migration files directory
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// Table to track executed migrations
type Migration struct {
	File string `db:"file,primarykey"`
}

func (Migration) TableName() string {
	return "migrations"
}

func (Migration) Indexes() []Index {
	return []Index{}
}

type Migrater struct {
	// Map of all available migrations. Key: Filename, Value: SQL content
	Migrations map[string]string
	DB         Database
}

func NewMigrater(db Database) Migrater {
	migrations := make(map[string]string)

	files, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			slog.Info("Skipping directory in migrations", "name", file.Name())
			continue
		}

		// Read the content of the migration file (sql)
		fileName := "migrations/" + file.Name()
		slog.Info("Found migration file", "name", fileName)

		content, err := migrationFiles.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
		migrations[file.Name()] = string(content)
	}

	return Migrater{
		DB:         db,
		Migrations: migrations,
	}
}

// Create all the required tables if they do not exist
func (m *Migrater) CreateTableIfNotExists(tables ...Table) error {
	for _, table := range tables {
		if m.DB.TableExists(table) {
			slog.Info("Table already exists, skipping creation", "table", table.TableName())
			continue
		}
		slog.Info("Creating table", "table", table.TableName())
		if err := m.DB.CreateTable(m.DB.AddTable(table)); err != nil {
			return err
		}
	}
	return nil
}

func (m *Migrater) Migrate() error {
	// Ensure the migrations table exists
	if !m.DB.TableExists(&Migration{}) {
		slog.Info("Creating migrations table")
		if err := m.DB.CreateTable(m.DB.AddTable(Migration{})); err != nil {
			return err
		}
	}

	// Fetch already executed migrations
	executedMigrations := make(map[string]bool)
	var migrations []Migration
	if _, err := m.DB.Select(&migrations, "SELECT * FROM migrations"); err != nil {
		return err
	}
	for _, migration := range migrations {
		executedMigrations[migration.File] = true
	}

	// Execute pending migrations
	for fileName, sql := range m.Migrations {
		if executedMigrations[fileName] {
			slog.Info("Migration already executed, skipping", "file", fileName)
			continue
		}
		slog.Info("Executing migration", "file", fileName)
		if _, err := m.DB.Exec(sql); err != nil {
			return err
		}
		// Record the executed migration
		if err := m.DB.Insert(&Migration{File: fileName}); err != nil {
			return err
		}
	}
	slog.Info("All migrations executed successfully")
	return nil
}
