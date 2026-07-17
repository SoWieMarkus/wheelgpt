package bot

import (
	"fmt"
	"regexp"
	"strings"
)

type Medal int

const (
	MedalChampion Medal = iota
	MedalAuthor
	MedalGold
	MedalSilver
	MedalBronze
	MedalNone
)

var stylingRegex = regexp.MustCompile(`\$([wnoitsgzWNOITSGZ]|[0-9A-Fa-f]{3})`)

type MapData struct {
	ChannelID    string
	UID          string
	Name         string
	Author       string
	AuthorTime   int64
	GoldTime     int64
	SilverTime   int64
	BronzeTime   int64
	ChampionTime int64
	TmxID        *int64
	WorldRecord  *int64
}

func (m *MapData) CleanName() string {
	return strings.TrimSpace(stylingRegex.ReplaceAllString(m.Name, ""))
}

func (m *MapData) Medal(t TMTime) Medal {
	ms := int64(t)
	switch {
	case m.ChampionTime > 0 && ms <= m.ChampionTime:
		return MedalChampion
	case ms <= m.AuthorTime:
		return MedalAuthor
	case ms <= m.GoldTime:
		return MedalGold
	case ms <= m.SilverTime:
		return MedalSilver
	case ms <= m.BronzeTime:
		return MedalBronze
	default:
		return MedalNone
	}
}

func (m *MapData) DiffToMedal(medal Medal, t TMTime) int64 {
	ms := int64(t)
	switch medal {
	case MedalChampion:
		return ms - m.ChampionTime
	case MedalAuthor:
		return ms - m.AuthorTime
	case MedalGold:
		return ms - m.GoldTime
	case MedalSilver:
		return ms - m.SilverTime
	case MedalBronze:
		return ms - m.BronzeTime
	}
	return 0
}

func (m *MapData) String() string {
	name := m.CleanName()
	author := TMTime(m.AuthorTime).String()
	ioLink := fmt.Sprintf("https://trackmania.io/#/leaderboard/%s", m.UID)

	hasTMX := m.TmxID != nil
	hasChampion := m.ChampionTime > 0

	tmxLink := ""
	if hasTMX {
		tmxLink = fmt.Sprintf("https://trackmania.exchange/maps/%d", *m.TmxID)
	}
	champion := TMTime(m.ChampionTime).String()

	switch {
	case hasTMX && hasChampion:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s %s | %s | %s`,
			name, m.Author,
			EmoteChampionMedal.Name, champion,
			EmoteAuthorMedal.Name, author,
			tmxLink, ioLink)
	case !hasTMX && hasChampion:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s %s | %s`,
			name, m.Author,
			EmoteChampionMedal.Name, champion,
			EmoteAuthorMedal.Name, author,
			ioLink)
	case hasTMX && !hasChampion:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s | %s`,
			name, m.Author,
			EmoteAuthorMedal.Name, author,
			tmxLink, ioLink)
	default:
		return fmt.Sprintf(`"%s" by %s | %s %s | %s`,
			name, m.Author,
			EmoteAuthorMedal.Name, author,
			ioLink)
	}
}
