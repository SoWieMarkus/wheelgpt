package leaderboard

import (
	"context"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/domain"
)

// AllTime is the key of the all-time leaderboard.
const AllTime = "all-time"

// Period is a leaderboard period: all-time or one calendar month (UTC).
type Period struct {
	Key  string
	From time.Time
	To   time.Time
}

// IsAllTime reports whether p is the all-time period.
func (p Period) IsAllTime() bool { return p.Key == AllTime }

// AllTimePeriod returns the all-time period.
func AllTimePeriod() Period { return Period{Key: AllTime} }

// MonthOf returns the calendar month (UTC) containing t.
func MonthOf(t time.Time) Period {
	t = t.UTC()
	from := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return Period{Key: from.Format("2006-01"), From: from, To: from.AddDate(0, 1, 0)}
}

// ParsePeriod parses "all-time" or "YYYY-MM".
func ParsePeriod(key string) (Period, error) {
	if key == AllTime {
		return AllTimePeriod(), nil
	}
	t, err := time.Parse("2006-01", key)
	if err != nil {
		return Period{}, domain.ErrInvalidPeriod
	}
	return MonthOf(t), nil
}

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// SetClock replaces the clock (tests).
func (s *Service) SetClock(now func() time.Time) { s.now = now }

// CurrentMonth is the month that is running now.
func (s *Service) CurrentMonth() Period { return MonthOf(s.now()) }

func (s *Service) Board(ctx context.Context, channelID string, p Period, query string, limit, offset int64) ([]Entry, int64, error) {
	return s.repo.Board(ctx, channelID, p, query, limit, offset)
}

// Find looks a viewer up by user id or, if userID is empty, by name.
func (s *Service) Find(ctx context.Context, channelID string, p Period, userID, displayName string) (Entry, bool, error) {
	return s.repo.Find(ctx, channelID, p, userID, displayName)
}

// Periods lists the current month, every month with points and all-time.
func (s *Service) Periods(ctx context.Context, channelID string) (current string, periods []string, err error) {
	months, err := s.repo.Months(ctx, channelID)
	if err != nil {
		return "", nil, err
	}
	current = s.CurrentMonth().Key
	if len(months) == 0 || months[0] != current {
		periods = append(periods, current)
	}
	periods = append(periods, months...)
	return current, append(periods, AllTime), nil
}
