// This file owns the SQLite connection and the schema migrations.

package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/url"

	"github.com/pressly/goose/v3"
	// Pure Go SQLite driver, no cgo needed.
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Open opens the SQLite database at path, creating it if needed.
//
//   - WAL lets readers work while a write is in progress.
//   - busy_timeout makes writers queue up instead of failing with SQLITE_BUSY.
//   - _txlock=immediate takes the write lock at BEGIN, which serializes all
//     transactions and avoids deadlocks from lock upgrades.
//
// Timestamps are not written through the driver's time format: the tables of
// the Prisma era store unix milliseconds, see package dbtime.
func Open(ctx context.Context, path string) (*sql.DB, error) {
	params := url.Values{}
	params.Add("_pragma", "foreign_keys(1)")
	params.Add("_pragma", "journal_mode(WAL)")
	params.Add("_pragma", "synchronous(NORMAL)")
	params.Add("_pragma", "busy_timeout(10000)")
	params.Set("_txlock", "immediate")

	db, err := sql.Open("sqlite", "file:"+path+"?"+params.Encode())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

// Migrate applies all pending migrations.
//
// Migrations are additive only. 00001 recreates the schema the Prisma backend
// left behind with CREATE ... IF NOT EXISTS, so it is a no-op on the
// production database. Later migrations only add tables, the old backend can
// keep using the file. goose tracks them in goose_db_version, Prisma's
// _prisma_migrations is never touched.
func Migrate(ctx context.Context, db *sql.DB) error {
	files, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("open migrations: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectSQLite3, db, files, goose.WithDisableGlobalRegistry(true))
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}
	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	for _, r := range results {
		slog.InfoContext(ctx, "applied migration", "version", r.Source.Version, "duration", r.Duration)
	}
	return nil
}
