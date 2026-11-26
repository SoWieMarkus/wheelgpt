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

	"github.com/SoWieMarkus/gokit/sqlite"
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
	return sqlite.InTx(ctx, s.db, func(tx *sql.Tx) error {
		return fn(s.WithTx(tx))
	})
}

// Ping verifies the database is reachable.
func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// DB exposes the raw handle for metrics.
func (s *Store) DB() *sql.DB {
	return s.db
}

// Error helpers from gokit/sqlite, re-exported so repositories only depend
// on this package.
var (
	IsNotFound        = sqlite.IsNotFound
	IsUniqueViolation = sqlite.IsUniqueViolation
)
