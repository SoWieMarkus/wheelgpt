package trackmania

import (
	"fmt"
	"regexp"

	"github.com/SoWieMarkus/wheelgpt/pkg/emote"
)

// stylingRegex is a regular expression that matches Trackmania styling codes in map names.
// See: https://wiki.trackmania.io/en/content-creation/text-styling
var stylingRegex = regexp.MustCompile(`\$([wnoitsgz$WNOITSGZ]|[0-9A-Fa-f]{3})`)

// Map represents a Trackmania map with its unique identifier, name, author, and medal times.
type Map struct {
	// UID is the unique identifier of a map.
	UID string

	// Name is the name of the map.
	Name string

	// Author is the name of the creator of a map.
	Author string

	// AuthorTime is the time that is required for author medal of a map.
	AuthorTime Time
	// GoldTime is the time that is required for gold medal of a map.
	GoldTime Time
	// SilverTime is the time that is required for silver medal of a map.
	SilverTime Time
	// BronzeTime is the time that is required for bronze medal of a map.
	BronzeTime Time

	// ChampionTime is the time that is required for champion medal of a map. Since this is only accessible through an OpenPlanet plugin this field is optional and may be nil if the information is not available.
	ChampionTime *Time

	// TmxID is the unique identifier of a map on Trackmania Exchange. Since not every map is available on Trackmania Exchange this field is optional and may be nil if the information is not available.
	TmxID *int64
}

// CleanName returns the name of the map with all trackmania styling codes removed.
// See: https://wiki.trackmania.io/en/content-creation/text-styling
func (m *Map) CleanName() string {
	return stylingRegex.ReplaceAllString(m.Name, "")
}

// GetMedal returns the medal that corresponds to the given time for the map.
func (m *Map) Medal(time Time) Medal {
	switch {
	case m.ChampionTime != nil && time <= *m.ChampionTime:
		return MedalChampion
	case time <= m.AuthorTime:
		return MedalAuthor
	case time <= m.GoldTime:
		return MedalGold
	case time <= m.SilverTime:
		return MedalSilver
	case time <= m.BronzeTime:
		return MedalBronze
	default:
		return MedalNone
	}
}

// TimeDifferenceToMedal returns the time difference between the given time and the required time for the specified medal on the map.
// If the medal is not available for the map, it returns 0.
func (m *Map) TimeDifferenceToMedal(medal Medal, time Time) Time {
	switch medal {
	case MedalChampion:
		if m.ChampionTime == nil {
			return 0
		}
		diff := *m.ChampionTime - time
		return diff
	case MedalAuthor:
		diff := m.AuthorTime - time
		return diff
	case MedalGold:
		diff := m.GoldTime - time
		return diff
	case MedalSilver:
		diff := m.SilverTime - time
		return diff
	case MedalBronze:
		diff := m.BronzeTime - time
		return diff
	default:
		return 0
	}
}

// String returns a string representation of the map, including its name, author, medal times, and links to Trackmania.io and Trackmania Exchange (if available).
func (m *Map) String() string {
	name := m.CleanName()

	trackmaniaIOLink := fmt.Sprintf("https://trackmania.io/#/leaderboard/%s", m.UID)
	hasChampionTime := m.ChampionTime != nil
	hasTmxID := m.TmxID != nil
	tmxLink := ""
	if hasTmxID {
		tmxLink = fmt.Sprintf("https://trackmania.exchange/map/%d", *m.TmxID)
	}

	switch {
	case hasChampionTime && hasTmxID:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s %s | %s | %s`,
			name, m.Author,
			emote.ChampionMedal, m.ChampionTime,
			emote.AuthorMedal, m.AuthorTime,
			trackmaniaIOLink, tmxLink,
		)
	case hasChampionTime && !hasTmxID:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s %s | %s`,
			name, m.Author,
			emote.ChampionMedal, m.ChampionTime,
			emote.AuthorMedal, m.AuthorTime,
			trackmaniaIOLink,
		)
	case !hasChampionTime && hasTmxID:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s | %s`,
			name, m.Author,
			emote.AuthorMedal, m.AuthorTime,
			trackmaniaIOLink, tmxLink,
		)
	default:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s`,
			name, m.Author,
			emote.AuthorMedal, m.AuthorTime,
			trackmaniaIOLink,
		)
	}
}
