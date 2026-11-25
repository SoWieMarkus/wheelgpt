package trackmania

import (
	"fmt"
	"regexp"
	"strings"
)

// Medal is a Trackmania medal, best first.
type Medal int

// Medals from best to worst.
const (
	MedalChampion Medal = iota
	MedalAuthor
	MedalGold
	MedalSilver
	MedalBronze
	MedalNone
)

func (m Medal) String() string {
	return [...]string{"champion", "author", "gold", "silver", "bronze", "none"}[m]
}

// Map is the map a streamer is currently playing.
type Map struct {
	UID          string
	Name         string
	Author       string
	AuthorTime   Time
	GoldTime     Time
	SilverTime   Time
	BronzeTime   Time
	ChampionTime Time // 0 if the map has no champion medal
	TMXID        int64
}

var stylingPattern = regexp.MustCompile(`\$([wnoitsgzWNOITSGZ$]|[0-9A-Fa-f]{3})`)

// StripStyling removes Trackmania text formatting codes like $fff or $o.
func StripStyling(name string) string {
	return strings.TrimSpace(stylingPattern.ReplaceAllString(name, ""))
}

// HasChampionTime reports whether a champion medal time is known.
func (m Map) HasChampionTime() bool { return m.ChampionTime > 0 }

// Medal returns the best medal t earns on this map.
func (m Map) Medal(t Time) Medal {
	switch {
	case m.HasChampionTime() && t <= m.ChampionTime:
		return MedalChampion
	case t <= m.AuthorTime:
		return MedalAuthor
	case t <= m.GoldTime:
		return MedalGold
	case t <= m.SilverTime:
		return MedalSilver
	case t <= m.BronzeTime:
		return MedalBronze
	default:
		return MedalNone
	}
}

// TMXURL links to Trackmania Exchange, or "" if the map is not on TMX.
func (m Map) TMXURL() string {
	if m.TMXID == 0 {
		return ""
	}
	return fmt.Sprintf("https://trackmania.exchange/maps/%d", m.TMXID)
}

// TrackmaniaIOURL links to the trackmania.io leaderboard of the map.
func (m Map) TrackmaniaIOURL() string {
	return "https://trackmania.io/#/leaderboard/" + m.UID
}

// String is the text of the !map command.
func (m Map) String() string {
	parts := []string{`"` + StripStyling(m.Name) + `" by ` + m.Author}
	if m.HasChampionTime() {
		parts = append(parts, fmt.Sprintf("%s  %s", EmoteChampionMedal.Name, m.ChampionTime))
	}
	parts = append(parts, fmt.Sprintf("%s  %s", EmoteAuthorMedal.Name, m.AuthorTime))
	if tmx := m.TMXURL(); tmx != "" {
		parts = append(parts, tmx)
	}
	parts = append(parts, m.TrackmaniaIOURL())
	return strings.Join(parts, " | ")
}
