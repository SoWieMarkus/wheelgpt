// Package plugin is the API of the Openplanet plugin: it authenticates
// plugin tokens and takes the map, room and personal bests of a streamer.
package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
)

// Token is the metadata of a plugin token. The token itself is never stored.
type Token struct {
	ChannelID  string
	CreatedAt  time.Time
	LastUsedAt *time.Time
}

// Repository is the persistence the plugin service needs.
type Repository interface {
	TokenByHash(ctx context.Context, hash string) (Token, bool, error)
	Token(ctx context.Context, channelID string) (Token, bool, error)
	TouchToken(ctx context.Context, channelID string, at time.Time) error
	// ReplaceToken stores a new token hash and clears the secret of old-style
	// tokens, which revokes them.
	ReplaceToken(ctx context.Context, channelID, hash string, at time.Time) error
}

type SQLRepository struct {
	st *store.Store
}

func NewRepository(st *store.Store) *SQLRepository {
	return &SQLRepository{st: st}
}

func toToken(row store.PluginToken) Token {
	t := Token{ChannelID: row.ChannelId, CreatedAt: time.UnixMilli(row.CreatedAt).UTC()}
	if row.LastUsedAt != nil {
		used := time.UnixMilli(*row.LastUsedAt).UTC()
		t.LastUsedAt = &used
	}
	return t
}

func (r *SQLRepository) TokenByHash(ctx context.Context, hash string) (Token, bool, error) {
	row, err := r.st.GetPluginTokenByHash(ctx, hash)
	if store.IsNotFound(err) {
		return Token{}, false, nil
	}
	if err != nil {
		return Token{}, false, fmt.Errorf("get plugin token: %w", err)
	}
	return toToken(row), true, nil
}

func (r *SQLRepository) Token(ctx context.Context, channelID string) (Token, bool, error) {
	row, err := r.st.GetPluginToken(ctx, channelID)
	if store.IsNotFound(err) {
		return Token{}, false, nil
	}
	if err != nil {
		return Token{}, false, fmt.Errorf("get plugin token: %w", err)
	}
	return toToken(row), true, nil
}

func (r *SQLRepository) TouchToken(ctx context.Context, channelID string, at time.Time) error {
	ms := at.UnixMilli()
	if err := r.st.TouchPluginToken(ctx, store.TouchPluginTokenParams{ChannelID: channelID, LastUsedAt: &ms}); err != nil {
		return fmt.Errorf("touch plugin token: %w", err)
	}
	return nil
}

func (r *SQLRepository) ReplaceToken(ctx context.Context, channelID, hash string, at time.Time) error {
	err := r.st.InTx(ctx, func(q *store.Queries) error {
		if err := q.UpsertPluginToken(ctx, store.UpsertPluginTokenParams{ChannelID: channelID, TokenHash: hash, CreatedAt: at.UnixMilli()}); err != nil {
			return err
		}
		return q.SetChannelLegacyToken(ctx, store.SetChannelLegacyTokenParams{ID: channelID, Token: ""})
	})
	if err != nil {
		return fmt.Errorf("replace plugin token: %w", err)
	}
	return nil
}
