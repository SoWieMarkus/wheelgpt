package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/go-gorp/gorp/v3"
)

// Table is an interface that database models must implement to be added to the database.
type Table interface {
	// Name of the table in the database.
	TableName() string
	// List of indexes to create on the table.
	Indexes() []Index
}

// Index represents a database index.
type Index struct {
	// Name of the index.
	Name string
	// Columns that make up the index.
	ColumnNames []string
}

type Database struct {
	*gorp.DbMap
	Config config.DatabaseConfig
}

// Constructs the database URL from the given configuration.
func getDatabaseURL(conf config.DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.User, conf.Password, conf.HostName, conf.Port, conf.DatabaseName)
}

// Create a new Postgres database connection.
func NewPostgresDB(ctx context.Context, conf config.DatabaseConfig) Database {
	databaseURL := getDatabaseURL(conf)
	slog.Info("connecting to database", "url", strings.Replace(databaseURL, conf.Password, "****", -1))
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		panic(err)
	}

	// If the wait time exceeds 10 seconds, we will panic.
	maxRetries := 10
	for i := range maxRetries {
		err := db.PingContext(ctx)
		if err == nil {
			break
		}
		if i == maxRetries-1 {
			panic("giving up connecting to database")
		}
		slog.Error("failed to connect to database, retrying...", "error", err)
		time.Sleep(1 * time.Second)
	}

	db.SetMaxOpenConns(16)
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	slog.Info("database is ready")
	return Database{DbMap: dbMap, Config: conf}
}

// Adds a Model table to the database.
func (db *Database) AddTable(table Table) *gorp.TableMap {
	tableName := table.TableName()
	slog.Info("adding table", "table", tableName, "model", table)
	tablemap := db.AddTableWithName(table, tableName)
	for _, index := range table.Indexes() {
		slog.Info("adding index", "index", index.Name, "table", tableName, "columns", index.ColumnNames)
		tablemap.AddIndex(index.Name, "Btree", index.ColumnNames)
	}
	return tablemap
}

// Adds missing functionality to gorp.DbMap which creates one table.
func (db *Database) CreateTable(table ...*gorp.TableMap) error {
	tx, err := db.Begin()
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return tx.Rollback()
	}
	for _, t := range table {
		slog.Info("creating table if exists", "table", t.TableName)
		// true means to add IF NOT EXISTS
		sql := t.SqlForCreate(true)
		if _, err := tx.Exec(sql); err != nil {
			return tx.Rollback()
		}
	}
	return tx.Commit()
}

// Check if a table exists in the database.
func (db *Database) TableExists(t Table) bool {
	var query string
	switch db.Dialect.(type) {
	case gorp.PostgresDialect:
		query = `SELECT EXISTS (
			SELECT 1
			FROM   information_schema.tables
			WHERE  table_name = :table_name
		);`
	case gorp.SqliteDialect:
		query = `SELECT EXISTS (
			SELECT 1
			FROM sqlite_master
			WHERE type='table' AND name = :table_name
		);`
	default:
		slog.Error("unsupported database dialect")
		return false
	}
	var exists bool
	err := db.SelectOne(&exists, query, map[string]any{"table_name": t.TableName()})
	if err != nil {
		slog.Error("failed to check if table exists", "error", err)
		return false
	}
	return exists
}

// Convenience function to close the database connection.
func (db *Database) Close() {
	if err := db.Db.Close(); err != nil {
		slog.Error("failed to close database connection", "error", err)
	}
}
