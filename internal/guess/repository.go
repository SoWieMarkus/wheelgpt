// Package guess runs the guessing game: viewers guess the next personal best,
// a personal best ends the round, the closest guesses win points.
package guess

import (
	"context"
	"fmt"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store/dbtime"
	"github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

// Guess is an open guess of the current round.
type Guess struct {
	UserID      string
	DisplayName string
	Time        trackmania.Time
	CreatedAt   time.Time
}

// RoundRecord is what is stored when a round ends.
type RoundRecord struct {
	ChannelID string
	Result    trackmania.Time
	MapUID    *string
	MapName   *string
	At        time.Time
	// Participants are the guesses old enough to count.
	Participants []Guess
	Winners      map[string]bool
	Points       int64
	Perfect      bool
}

// Repository is the persistence the guess service needs.
type Repository interface {
	Place(ctx context.Context, channelID string, g Guess) error
	Get(ctx context.Context, channelID, userID string) (Guess, bool, error)
	ListOpen(ctx context.Context, channelID string) ([]Guess, error)
	Reset(ctx context.Context, channelID string) error
	// SaveRound stores the round and its results, adds the all-time points
	// and clears all open guesses, in one transaction. It returns the id.
	SaveRound(ctx context.Context, r RoundRecord) (int64, error)
}

type SQLRepository struct {
	st *store.Store
}

func NewRepository(st *store.Store) *SQLRepository {
	return &SQLRepository{st: st}
}

func (r *SQLRepository) Place(ctx context.Context, channelID string, g Guess) error {
	err := r.st.UpsertGuess(ctx, store.UpsertGuessParams{
		ChannelID: channelID, UserID: g.UserID, DisplayName: g.DisplayName,
		Time: g.Time.Milliseconds(), CreatedAt: dbtime.From(g.CreatedAt),
	})
	if err != nil {
		return fmt.Errorf("place guess: %w", err)
	}
	return nil
}

func (r *SQLRepository) Get(ctx context.Context, channelID, userID string) (Guess, bool, error) {
	row, err := r.st.GetGuess(ctx, store.GetGuessParams{ChannelID: channelID, UserID: userID})
	if store.IsNotFound(err) {
		return Guess{}, false, nil
	}
	if err != nil {
		return Guess{}, false, fmt.Errorf("get guess: %w", err)
	}
	return toGuess(row), true, nil
}

func toGuess(row store.Guess) Guess {
	return Guess{UserID: row.UserId, DisplayName: row.DisplayName, Time: trackmania.Time(row.Time), CreatedAt: row.CreatedAt.Time}
}

func (r *SQLRepository) ListOpen(ctx context.Context, channelID string) ([]Guess, error) {
	rows, err := r.st.ListGuesses(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("list guesses: %w", err)
	}
	out := make([]Guess, len(rows))
	for i, row := range rows {
		out[i] = toGuess(row)
	}
	return out, nil
}

func (r *SQLRepository) Reset(ctx context.Context, channelID string) error {
	if _, err := r.st.DeleteGuesses(ctx, channelID); err != nil {
		return fmt.Errorf("reset guesses: %w", err)
	}
	return nil
}

func (r *SQLRepository) SaveRound(ctx context.Context, rec RoundRecord) (int64, error) {
	var id int64
	err := r.st.InTx(ctx, func(q *store.Queries) error {
		at := rec.At.UnixMilli()
		round, err := q.CreateRound(ctx, store.CreateRoundParams{
			ChannelID: rec.ChannelID, ResultTime: rec.Result.Milliseconds(),
			MapUid: rec.MapUID, MapName: rec.MapName, CreatedAt: at,
		})
		if err != nil {
			return err
		}
		id = round.ID
		for _, g := range rec.Participants {
			var points int64
			if rec.Winners[g.UserID] {
				points = rec.Points
			}
			if err := q.CreateResult(ctx, store.CreateResultParams{
				RoundID: round.ID, ChannelID: rec.ChannelID, UserID: g.UserID, DisplayName: g.DisplayName,
				GuessTime: g.Time.Milliseconds(), ResultTime: rec.Result.Milliseconds(),
				Points: points, Perfect: points > 0 && rec.Perfect, CreatedAt: at,
			}); err != nil {
				return err
			}
			if points == 0 {
				continue
			}
			var perfect int64
			if rec.Perfect {
				perfect = 1
			}
			// The all-time board is the table of the old backend, so points
			// from before the switchover are kept.
			if err := q.AddAllTimePoints(ctx, store.AddAllTimePointsParams{
				ChannelID: rec.ChannelID, UserID: g.UserID, DisplayName: g.DisplayName,
				Points: points, Perfect: perfect, UpdatedAt: dbtime.From(rec.At),
			}); err != nil {
				return err
			}
		}
		_, err = q.DeleteGuesses(ctx, rec.ChannelID)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("save round: %w", err)
	}
	return id, nil
}
