package guess

import (
	"context"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

// Maps returns the map a channel is on, or nil.
type Maps interface {
	CurrentMap(ctx context.Context, channelID string) (*trackmania.Map, error)
}

// Recorder counts guesses (metrics).
type Recorder interface {
	Guess()
}

type Service struct {
	repo Repository
	maps Maps
	rec  Recorder
	now  func() time.Time
}

func NewService(repo Repository, maps Maps, rec Recorder) *Service {
	return &Service{repo: repo, maps: maps, rec: rec, now: time.Now}
}

// SetClock replaces the clock (tests).
func (s *Service) SetClock(now func() time.Time) { s.now = now }

// Now returns the current time of the service clock.
func (s *Service) Now() time.Time { return s.now() }

// Place stores or replaces a viewer's guess. A changed guess counts as new
// for the minimum guess age.
func (s *Service) Place(ctx context.Context, channelID, userID, displayName string, t trackmania.Time) error {
	if err := s.repo.Place(ctx, channelID, Guess{UserID: userID, DisplayName: displayName, Time: t, CreatedAt: s.now()}); err != nil {
		return err
	}
	s.rec.Guess()
	return nil
}

func (s *Service) Get(ctx context.Context, channelID, userID string) (Guess, bool, error) {
	return s.repo.Get(ctx, channelID, userID)
}

func (s *Service) Reset(ctx context.Context, channelID string) error {
	return s.repo.Reset(ctx, channelID)
}

// Round is the outcome of a personal best.
type Round struct {
	ID           int64
	Result       trackmania.Time
	MapUID       *string
	MapName      *string
	At           time.Time
	Participants int
	Winners      []trackmania.Guess
	Points       int64
	Perfect      bool
	// Message is the chat announcement.
	Message string
}

// EvaluatePB ends the round: guesses old enough are evaluated, the winners
// get points, the round is stored and all guesses are cleared.
func (s *Service) EvaluatePB(ctx context.Context, c channel.Channel, pb trackmania.Time) (Round, error) {
	now := s.now()
	open, err := s.repo.ListOpen(ctx, c.ID)
	if err != nil {
		return Round{}, err
	}
	current, err := s.maps.CurrentMap(ctx, c.ID)
	if err != nil {
		return Round{}, err
	}

	cutoff := now.Add(-c.Settings.GuessMinAge())
	var participants []Guess
	var candidates []trackmania.Guess
	for _, g := range open {
		if g.CreatedAt.After(cutoff) {
			continue // placed too late
		}
		participants = append(participants, g)
		candidates = append(candidates, trackmania.Guess{UserID: g.UserID, DisplayName: g.DisplayName, Time: g.Time})
	}
	winners := trackmania.Evaluate(candidates, pb)
	points, perfect := trackmania.PointsFor(winners, pb)
	isWinner := make(map[string]bool, len(winners))
	for _, w := range winners {
		isWinner[w.UserID] = true
	}

	round := Round{
		Result: pb, At: now, Participants: len(participants), Winners: winners,
		Points: points, Perfect: perfect, Message: trackmania.ResultMessage(current, pb, winners),
	}
	if current != nil {
		uid, name := current.UID, trackmania.StripStyling(current.Name)
		round.MapUID, round.MapName = &uid, &name
	}
	round.ID, err = s.repo.SaveRound(ctx, RoundRecord{
		ChannelID: c.ID, Result: pb, MapUID: round.MapUID, MapName: round.MapName, At: now,
		Participants: participants, Winners: isWinner, Points: points, Perfect: perfect,
	})
	return round, err
}
