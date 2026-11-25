package ingame

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

// TMX looks up Trackmania Exchange ids.
type TMX interface {
	MapID(ctx context.Context, uid string) (id int64, ok bool, err error)
}

type Service struct {
	repo Repository
	tmx  TMX
	wg   sync.WaitGroup
}

func NewService(repo Repository, tmx TMX) *Service {
	return &Service{repo: repo, tmx: tmx}
}

func (s *Service) Map(ctx context.Context, channelID string) (trackmania.Map, error) {
	return s.repo.Map(ctx, channelID)
}

// CurrentMap returns the map or nil when the streamer is not on one.
func (s *Service) CurrentMap(ctx context.Context, channelID string) (*trackmania.Map, error) {
	m, err := s.repo.Map(ctx, channelID)
	if errors.Is(err, domain.ErrMapNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// SetMap stores the map, clears open guesses and looks the map up on
// Trackmania Exchange in the background.
func (s *Service) SetMap(ctx context.Context, channelID string, m trackmania.Map) (trackmania.Map, error) {
	m.TMXID = 0
	if err := s.repo.SetMap(ctx, channelID, m); err != nil {
		return trackmania.Map{}, err
	}
	if s.tmx != nil {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
			defer cancel()
			id, ok, err := s.tmx.MapID(ctx, m.UID)
			if err != nil {
				slog.WarnContext(ctx, "tmx lookup failed", "uid", m.UID, "error", err)
				return
			}
			if ok {
				if err := s.repo.SetMapTMXID(ctx, channelID, m.UID, id); err != nil {
					slog.WarnContext(ctx, "store tmx id failed", "uid", m.UID, "error", err)
				}
			}
		}()
	}
	return m, nil
}

func (s *Service) ClearMap(ctx context.Context, channelID string) error {
	return s.repo.ClearMap(ctx, channelID)
}

func (s *Service) Room(ctx context.Context, channelID string) (Room, error) {
	return s.repo.Room(ctx, channelID)
}

func (s *Service) SetRoom(ctx context.Context, channelID string, r Room) error {
	return s.repo.SetRoom(ctx, channelID, r)
}

func (s *Service) ClearRoom(ctx context.Context, channelID string) error {
	return s.repo.ClearRoom(ctx, channelID)
}

func (s *Service) ClearAll(ctx context.Context, channelID string) error {
	return s.repo.ClearAll(ctx, channelID)
}

// Wait blocks until background lookups are done.
func (s *Service) Wait() { s.wg.Wait() }
