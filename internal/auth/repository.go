// Package auth logs streamers in with Twitch and keeps their website
// session in an HttpOnly cookie. Only the hash of a session token is stored.
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
)

// Session is a stored website session.
type Session struct {
	ID        string // SHA-256 of the cookie value
	ChannelID string
	ExpiresAt time.Time
}

// Repository is the persistence the auth service needs.
type Repository interface {
	CreateSession(ctx context.Context, s Session, now time.Time) error
	GetSession(ctx context.Context, id string, now time.Time) (Session, error)
	DeleteSession(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context, now time.Time) error
}

type SQLRepository struct {
	st *store.Store
}

func NewRepository(st *store.Store) *SQLRepository {
	return &SQLRepository{st: st}
}

func (r *SQLRepository) CreateSession(ctx context.Context, s Session, now time.Time) error {
	err := r.st.CreateSession(ctx, store.CreateSessionParams{
		ID: s.ID, ChannelID: s.ChannelID, CreatedAt: now.UnixMilli(), ExpiresAt: s.ExpiresAt.UnixMilli(),
	})
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}
	return nil
}

func (r *SQLRepository) GetSession(ctx context.Context, id string, now time.Time) (Session, error) {
	row, err := r.st.GetSession(ctx, store.GetSessionParams{ID: id, Now: now.UnixMilli()})
	if store.IsNotFound(err) {
		return Session{}, domain.ErrUnauthenticated
	}
	if err != nil {
		return Session{}, fmt.Errorf("get session: %w", err)
	}
	return Session{ID: row.ID, ChannelID: row.ChannelId, ExpiresAt: time.UnixMilli(row.ExpiresAt).UTC()}, nil
}

func (r *SQLRepository) DeleteSession(ctx context.Context, id string) error {
	if err := r.st.DeleteSession(ctx, id); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (r *SQLRepository) DeleteExpired(ctx context.Context, now time.Time) error {
	if err := r.st.DeleteExpiredSessions(ctx, now.UnixMilli()); err != nil {
		return fmt.Errorf("delete expired sessions: %w", err)
	}
	return nil
}
