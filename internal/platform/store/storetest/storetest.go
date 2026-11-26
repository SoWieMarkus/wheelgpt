// Package storetest provides real SQLite databases for tests. Every call gets
// its own file in a temp directory, so tests are isolated and can run in
// parallel with nothing to set up.
package storetest

import (
	"context"
	"testing"

	"github.com/SoWieMarkus/gokit/sqlite/sqlitetest"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store/dbtime"
)

// New opens a fresh, fully migrated database.
func New(t *testing.T) *store.Store {
	t.Helper()
	return store.NewStore(sqlitetest.Open(t, store.Migrations()))
}

// Channel creates a channel like the old backend did: with a legacy plugin
// token secret and no ChannelState row.
func Channel(t *testing.T, s *store.Store, id, login string) {
	t.Helper()
	err := s.UpsertChannel(context.Background(), store.UpsertChannelParams{
		ID: id, Login: login, Token: "legacy-secret", DisplayName: login,
		ProfileImage: "https://example.com/" + login + ".png", CreatedAt: dbtime.Now(),
	})
	if err != nil {
		t.Fatalf("create channel %s: %v", login, err)
	}
}

// Live marks a channel live.
func Live(t *testing.T, s *store.Store, id string) {
	t.Helper()
	if err := s.SetChannelLive(context.Background(), store.SetChannelLiveParams{ID: id, IsLive: true}); err != nil {
		t.Fatal(err)
	}
}
