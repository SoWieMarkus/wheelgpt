// This file owns the schema migrations. Opening the database and running
// them is done by gokit/sqlite.
//
// Timestamps are not written through the driver's time format: the tables of
// the Prisma era store unix milliseconds, see package dbtime.

package store

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"

	"github.com/SoWieMarkus/gokit/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrations returns the goose migrations of the schema.
func Migrations() fs.FS {
	files, err := fs.Sub(migrations, "migrations")
	if err != nil {
		panic(err) // the directory is embedded, this cannot fail
	}
	return files
}

// Migrate applies all pending migrations.
//
// Migrations are additive only. 00001 recreates the schema the Prisma backend
// left behind with CREATE ... IF NOT EXISTS, so it is a no-op on the
// production database. Later migrations only add tables, the old backend can
// keep using the file. goose tracks them in goose_db_version, Prisma's
// _prisma_migrations is never touched.
func Migrate(ctx context.Context, db *sql.DB) error {
	return sqlite.Migrate(ctx, db, Migrations())
}
