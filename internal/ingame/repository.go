// Package ingame keeps what the streamer is doing in Trackmania right now:
// the current map and server room, as sent by the Openplanet plugin.
package ingame

import (
	"context"
	"fmt"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
	"github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

// Room is the server room the streamer plays in.
type Room struct {
	Login           string
	Name            string
	NumberOfPlayers int64
	MaxPlayers      int64
}

// Repository is the persistence the ingame service needs.
type Repository interface {
	Map(ctx context.Context, channelID string) (trackmania.Map, error)
	// SetMap replaces the map and clears the open guesses, which belonged to
	// the previous map.
	SetMap(ctx context.Context, channelID string, m trackmania.Map) error
	SetMapTMXID(ctx context.Context, channelID, uid string, tmxID int64) error
	ClearMap(ctx context.Context, channelID string) error
	Room(ctx context.Context, channelID string) (Room, error)
	SetRoom(ctx context.Context, channelID string, r Room) error
	ClearRoom(ctx context.Context, channelID string) error
	// ClearAll removes map, room and open guesses (stream went offline).
	ClearAll(ctx context.Context, channelID string) error
}

type SQLRepository struct {
	st *store.Store
}

func NewRepository(st *store.Store) *SQLRepository {
	return &SQLRepository{st: st}
}

func (r *SQLRepository) Map(ctx context.Context, channelID string) (trackmania.Map, error) {
	row, err := r.st.GetMap(ctx, channelID)
	if store.IsNotFound(err) {
		return trackmania.Map{}, ErrMapNotFound
	}
	if err != nil {
		return trackmania.Map{}, fmt.Errorf("get map: %w", err)
	}
	m := trackmania.Map{
		UID: row.Uid, Name: row.Name, Author: row.Author,
		AuthorTime: trackmania.Time(row.AuthorTime), GoldTime: trackmania.Time(row.GoldTime),
		SilverTime: trackmania.Time(row.SilverTime), BronzeTime: trackmania.Time(row.BronzeTime),
		ChampionTime: trackmania.Time(row.ChampionTime),
	}
	if row.TmxId != nil {
		m.TMXID = *row.TmxId
	}
	return m, nil
}

func (r *SQLRepository) SetMap(ctx context.Context, channelID string, m trackmania.Map) error {
	err := r.st.InTx(ctx, func(q *store.Queries) error {
		if _, err := q.DeleteGuesses(ctx, channelID); err != nil {
			return err
		}
		return q.UpsertMap(ctx, store.UpsertMapParams{
			ChannelID: channelID, Uid: m.UID, Name: m.Name, Author: m.Author,
			AuthorTime: int64(m.AuthorTime), GoldTime: int64(m.GoldTime), SilverTime: int64(m.SilverTime),
			BronzeTime: int64(m.BronzeTime), ChampionTime: int64(m.ChampionTime),
		})
	})
	if err != nil {
		return fmt.Errorf("set map: %w", err)
	}
	return nil
}

func (r *SQLRepository) SetMapTMXID(ctx context.Context, channelID, uid string, tmxID int64) error {
	if err := r.st.SetMapTmxID(ctx, store.SetMapTmxIDParams{TmxID: &tmxID, ChannelID: channelID, Uid: uid}); err != nil {
		return fmt.Errorf("set tmx id: %w", err)
	}
	return nil
}

func (r *SQLRepository) ClearMap(ctx context.Context, channelID string) error {
	err := r.st.InTx(ctx, func(q *store.Queries) error {
		if err := q.DeleteMap(ctx, channelID); err != nil {
			return err
		}
		_, err := q.DeleteGuesses(ctx, channelID)
		return err
	})
	if err != nil {
		return fmt.Errorf("clear map: %w", err)
	}
	return nil
}

func (r *SQLRepository) Room(ctx context.Context, channelID string) (Room, error) {
	row, err := r.st.GetRoom(ctx, channelID)
	if store.IsNotFound(err) {
		return Room{}, ErrRoomNotFound
	}
	if err != nil {
		return Room{}, fmt.Errorf("get room: %w", err)
	}
	return Room{Login: row.Login, Name: row.Name, NumberOfPlayers: row.NumberOfPlayers, MaxPlayers: row.MaxPlayers}, nil
}

func (r *SQLRepository) SetRoom(ctx context.Context, channelID string, room Room) error {
	err := r.st.UpsertRoom(ctx, store.UpsertRoomParams{
		ChannelID: channelID, Login: room.Login, Name: room.Name,
		NumberOfPlayers: room.NumberOfPlayers, MaxPlayers: room.MaxPlayers,
	})
	if err != nil {
		return fmt.Errorf("set room: %w", err)
	}
	return nil
}

func (r *SQLRepository) ClearRoom(ctx context.Context, channelID string) error {
	if err := r.st.DeleteRoom(ctx, channelID); err != nil {
		return fmt.Errorf("clear room: %w", err)
	}
	return nil
}

func (r *SQLRepository) ClearAll(ctx context.Context, channelID string) error {
	err := r.st.InTx(ctx, func(q *store.Queries) error {
		if err := q.DeleteMap(ctx, channelID); err != nil {
			return err
		}
		if err := q.DeleteRoom(ctx, channelID); err != nil {
			return err
		}
		_, err := q.DeleteGuesses(ctx, channelID)
		return err
	})
	if err != nil {
		return fmt.Errorf("clear ingame state: %w", err)
	}
	return nil
}
