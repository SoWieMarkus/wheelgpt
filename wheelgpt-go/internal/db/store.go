package db

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Open(dsn string) *sql.DB {
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		panic(err)
	}
	conn.SetMaxOpenConns(1) // SQLite does not support concurrent writes
	return conn
}

func Migrate(conn *sql.DB) error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	return goose.Up(conn, "migrations")
}

// Store wraps the sqlc-generated Queries with the raw *sql.DB for transactions.
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{Queries: New(db), db: db}
}

func (s *Store) WithTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := fn(s.Queries.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}
