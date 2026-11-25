package store_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
)

// legacyDB builds a database exactly like the Prisma backend did, by running
// its migrations, and fills it with rows the way Prisma wrote them.
func legacyDB(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "legacy.db")
	conn, err := sql.Open("sqlite", "file:"+path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	files, err := filepath.Glob("testdata/prisma_migrations/*/migration.sql")
	if err != nil || len(files) == 0 {
		t.Fatalf("no prisma migrations: %v", err)
	}
	sort.Strings(files)
	for _, f := range files {
		body, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Exec(string(body)); err != nil {
			t.Fatalf("%s: %v", f, err)
		}
	}
	// Prisma writes DateTime as unix milliseconds; DEFAULT CURRENT_TIMESTAMP
	// (rows created by a migration) writes text.
	for _, s := range []string{
		`CREATE TABLE "_prisma_migrations" ("id" TEXT PRIMARY KEY, "migration_name" TEXT NOT NULL)`,
		`INSERT INTO "_prisma_migrations" VALUES ('1', '20260716095641_add_guess_min_required_age')`,
		`INSERT INTO "Channel" ("id","login","token","displayName","profileImage","isLive","createdAt","guessDelayTime","guessMinRequiredAgeTime","botActiveWhenOffline","usagePublic")
		 VALUES ('100','granadyy','secret','GranaDyy','https://img',1,1717000000000,2.5,1,0,1)`,
		`INSERT INTO "Channel" ("id","login","token","displayName","profileImage") VALUES ('200','other','secret2','Other','https://img2')`,
		`INSERT INTO "Guess" ("channelId","time","userId","displayName","createdAt") VALUES ('100',61234,'1','Viewer',1717000000123)`,
		`INSERT INTO "Guess" ("channelId","time","userId","displayName") VALUES ('100',62000,'2','Late')`,
		`INSERT INTO "GuesserLeaderboard" ("channelId","userId","displayName","points","perfectGuessCount","updatedAt") VALUES ('100','1','Viewer',12,2,1717000000000)`,
		`INSERT INTO "GuesserLeaderboard" ("channelId","userId","displayName","points") VALUES ('100','2','Late',3)`,
		`INSERT INTO "TrackmaniaMap" ("channelId","uid","name","author","authorTime","goldTime","silverTime","bronzeTime","championTime","tmxId") VALUES ('100','uid1','$fffMap','Me',1,2,3,4,0,99)`,
	} {
		if _, err := conn.Exec(s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
	return path
}

func open(t *testing.T, path string) *store.Store {
	t.Helper()
	ctx := context.Background()
	db, err := store.Open(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := store.Migrate(ctx, db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return store.NewStore(db)
}

func TestMigrateLegacyDatabase(t *testing.T) {
	ctx := context.Background()
	path := legacyDB(t)
	st := open(t, path)

	ch, err := st.GetChannel(ctx, "100")
	if err != nil {
		t.Fatal(err)
	}
	if ch.Login != "granadyy" || !ch.IsLive || ch.GuessDelayTime != 2.5 || ch.BotActiveWhenOffline || !ch.UsagePublic {
		t.Errorf("channel = %+v", ch)
	}
	if want := time.UnixMilli(1717000000000).UTC(); !ch.CreatedAt.Equal(want) {
		t.Errorf("createdAt = %v, want %v", ch.CreatedAt, want)
	}
	// No ChannelState row: legacy channel, bot enabled, not re-authorized.
	if !ch.BotEnabled || ch.AuthorizedAt != nil || ch.ChatSubscription != "unknown" {
		t.Errorf("legacy state = %+v", ch)
	}

	other, err := st.GetChannel(ctx, "200")
	if err != nil {
		t.Fatal(err)
	}
	if other.CreatedAt.Year() < 2020 {
		t.Errorf("text default createdAt not parsed: %v", other.CreatedAt)
	}

	guesses, err := st.ListGuesses(ctx, "100")
	if err != nil || len(guesses) != 2 || guesses[0].CreatedAt.UnixMilli() != 1717000000123 {
		t.Fatalf("guesses = %+v, %v", guesses, err)
	}

	board, err := st.ListAllTimeLeaderboard(ctx, store.ListAllTimeLeaderboardParams{ChannelID: "100", Lim: 10})
	if err != nil || len(board) != 2 || board[0].UserId != "1" || board[0].Position != 1 || board[1].Points != 3 {
		t.Fatalf("leaderboard = %+v, %v", board, err)
	}

	m, err := st.GetMap(ctx, "100")
	if err != nil || m.TmxId == nil || *m.TmxId != 99 {
		t.Fatalf("map = %+v, %v", m, err)
	}

	// Reopening applies nothing twice, and the Prisma tables keep their columns.
	st = open(t, path)
	var cols int
	if err := st.DB().QueryRowContext(ctx, `SELECT COUNT(*) FROM pragma_table_info('Channel')`).Scan(&cols); err != nil || cols != 11 {
		t.Errorf("Channel has %d columns (%v), want the 11 of the Prisma schema", cols, err)
	}
	var prisma int
	if err := st.DB().QueryRowContext(ctx, `SELECT COUNT(*) FROM "_prisma_migrations"`).Scan(&prisma); err != nil || prisma != 1 {
		t.Errorf("_prisma_migrations must be left alone: %d, %v", prisma, err)
	}
}
