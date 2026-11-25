// Package store is the data access layer for SQLite.
//
// The SQL lives in queries/*.sql and migrations/*.sql, sqlc generates the
// typed query methods (*.gen.go) from it. This file adds the hand written
// transaction support on top. Feature packages wrap these queries in their
// own repositories, so services never see SQL or sqlc types.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// Store bundles the generated queries with the database handle so
// repositories can run several queries atomically.
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{Queries: New(db), db: db}
}

// InTx runs fn inside a transaction. The transaction is committed when fn
// returns nil and rolled back otherwise.
func (s *Store) InTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	// Rollback after a successful commit is a no-op.
	defer tx.Rollback() //nolint:errcheck

	if err := fn(s.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit()
}

// Ping verifies the database is reachable.
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// DB exposes the raw handle for metrics.
func (s *Store) DB() *sql.DB {
	return s.db
}

// IsNotFound reports whether err means "no row matched".
func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

// IsUniqueViolation reports whether err is a UNIQUE constraint failure.
// modernc.org/sqlite reports constraint errors only through their message.
func IsUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
