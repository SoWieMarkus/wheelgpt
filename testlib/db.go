package testlib

import (
	"database/sql"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/go-gorp/gorp/v3"
	_ "github.com/mattn/go-sqlite3"
)

type TestDatabase struct {
	*gorp.DbMap
	Close func()
}

func SQLite(t *testing.T) TestDatabase {
	var db TestDatabase
	// To run tests faster, the default is running with sqlite.
	slog.Info("Using sqlite")
	tmpDir := t.TempDir()
	sqlDB, err := sql.Open("sqlite3", tmpDir+"/test.db")
	if err != nil {
		t.Fatal(err)
	}
	db.DbMap = &gorp.DbMap{Db: sqlDB, Dialect: gorp.SqliteDialect{}}
	db.Close = func() {}
	db.TraceOn("[gorp]", log.New(os.Stdout, "cortex:", log.Lmicroseconds))
	return db
}
