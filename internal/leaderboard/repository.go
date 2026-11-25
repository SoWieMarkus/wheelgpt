// Package leaderboard ranks the guessers of a channel per calendar month and
// of all time, and keeps the history of rounds and guesses.
package leaderboard

import (
	"context"
	"fmt"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
)

// Entry is one leaderboard line.
type Entry struct {
	Position       int64
	UserID         string
	DisplayName    string
	Points         int64
	PerfectGuesses int64
}

// Repository is the persistence the leaderboard service needs.
type Repository interface {
	Board(ctx context.Context, channelID string, p Period, query string, limit, offset int64) ([]Entry, int64, error)
	// Find looks a viewer up by user id or, if userID is empty, by name.
	Find(ctx context.Context, channelID string, p Period, userID, displayName string) (Entry, bool, error)
	Months(ctx context.Context, channelID string) ([]string, error)
}

type SQLRepository struct {
	st *store.Store
}

func NewRepository(st *store.Store) *SQLRepository {
	return &SQLRepository{st: st}
}

func (r *SQLRepository) Board(ctx context.Context, channelID string, p Period, query string, limit, offset int64) ([]Entry, int64, error) {
	if p.IsAllTime() {
		rows, err := r.st.ListAllTimeLeaderboard(ctx, store.ListAllTimeLeaderboardParams{ChannelID: channelID, Query: query, Lim: limit, Off: offset})
		if err != nil {
			return nil, 0, fmt.Errorf("list all-time leaderboard: %w", err)
		}
		total, err := r.st.CountAllTimeLeaderboard(ctx, store.CountAllTimeLeaderboardParams{ChannelID: channelID, Query: query})
		if err != nil {
			return nil, 0, fmt.Errorf("count all-time leaderboard: %w", err)
		}
		out := make([]Entry, len(rows))
		for i, e := range rows {
			out[i] = Entry{e.Position, e.UserId, e.DisplayName, e.Points, e.PerfectGuessCount}
		}
		return out, total, nil
	}
	rows, err := r.st.ListPeriodLeaderboard(ctx, store.ListPeriodLeaderboardParams{
		ChannelID: channelID, FromMs: p.From.UnixMilli(), ToMs: p.To.UnixMilli(), Query: query, Lim: limit, Off: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list monthly leaderboard: %w", err)
	}
	total, err := r.st.CountPeriodLeaderboard(ctx, store.CountPeriodLeaderboardParams{
		ChannelID: channelID, FromMs: p.From.UnixMilli(), ToMs: p.To.UnixMilli(), Query: query,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("count monthly leaderboard: %w", err)
	}
	out := make([]Entry, len(rows))
	for i, e := range rows {
		out[i] = Entry{e.Position, e.UserId, e.DisplayName, e.Points, e.PerfectCount}
	}
	return out, total, nil
}

func (r *SQLRepository) Find(ctx context.Context, channelID string, p Period, userID, displayName string) (Entry, bool, error) {
	if p.IsAllTime() {
		e, err := r.st.FindAllTimeEntry(ctx, store.FindAllTimeEntryParams{ChannelID: channelID, UserID: userID, DisplayName: displayName})
		if store.IsNotFound(err) {
			return Entry{}, false, nil
		}
		if err != nil {
			return Entry{}, false, fmt.Errorf("find all-time entry: %w", err)
		}
		return Entry{e.Position, e.UserId, e.DisplayName, e.Points, e.PerfectGuessCount}, true, nil
	}
	e, err := r.st.FindPeriodEntry(ctx, store.FindPeriodEntryParams{
		ChannelID: channelID, FromMs: p.From.UnixMilli(), ToMs: p.To.UnixMilli(), UserID: userID, DisplayName: displayName,
	})
	if store.IsNotFound(err) {
		return Entry{}, false, nil
	}
	if err != nil {
		return Entry{}, false, fmt.Errorf("find monthly entry: %w", err)
	}
	return Entry{e.Position, e.UserId, e.DisplayName, e.Points, e.PerfectCount}, true, nil
}

func (r *SQLRepository) Months(ctx context.Context, channelID string) ([]string, error) {
	months, err := r.st.ListLeaderboardMonths(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("list months: %w", err)
	}
	return months, nil
}
