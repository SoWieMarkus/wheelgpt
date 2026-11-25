package trackmania

import (
	"fmt"
	"strings"
)

// Points awarded per round.
const (
	PointsPerfectGuess = 5
	PointsClosestGuess = 1
)

// Guess is a viewer's guess for the next personal best.
type Guess struct {
	UserID      string
	DisplayName string
	Time        Time
}

// Evaluate returns the guesses closest to result. Ties all win.
func Evaluate(guesses []Guess, result Time) []Guess {
	var winners []Guess
	best := Time(-1)
	for _, g := range guesses {
		diff := (result - g.Time).Abs()
		switch {
		case best < 0 || diff < best:
			best = diff
			winners = []Guess{g}
		case diff == best:
			winners = append(winners, g)
		}
	}
	return winners
}

// PointsFor returns the points each winner gets and whether it was perfect.
func PointsFor(winners []Guess, result Time) (points int64, perfect bool) {
	if len(winners) == 0 {
		return 0, false
	}
	if winners[0].Time == result {
		return PointsPerfectGuess, true
	}
	return PointsClosestGuess, false
}

// Mention formats a chat mention.
func Mention(name string) string {
	name = strings.TrimPrefix(name, "@")
	if name == "" {
		return ""
	}
	return "@" + name
}

// ResultMessage is the chat message after a new personal best.
func ResultMessage(current *Map, pb Time, winners []Guess) string {
	return PBMessage(current, pb) + " " + WinnerMessage(winners, pb)
}

// PBMessage comments on the personal best itself.
func PBMessage(current *Map, pb Time) string {
	if current == nil {
		return EmoteYek.Name + " I got a new PB but I didn't know you are on a map?"
	}
	prefix := fmt.Sprintf("NEW PERSONAL BEST %s %s", pb, EmoteDinkDonk.Name)
	missed := func(medal string, diff Time) string {
		return fmt.Sprintf("%s HAHAHAHAH HE MISSED %s BY %ss %s %s %s",
			prefix, medal, diff.Abs(), EmoteOmegalul.Name, EmoteICANT.Name, EmotePepePoint.Name)
	}

	switch current.Medal(pb) {
	case MedalChampion:
		return fmt.Sprintf("%s That's Champion Medal %s %s", prefix, EmoteChampionMedal.Name, EmoteBwoah.Name)
	case MedalAuthor:
		if diff := pb - current.ChampionTime; diff < 10 {
			return missed("CHAMPION", diff)
		}
		return fmt.Sprintf("%s That's %s %s", prefix, EmoteAuthorMedal.Name, EmoteOk.Name)
	case MedalGold:
		if diff := pb - current.AuthorTime; diff < 10 {
			return missed("AUTHOR", diff)
		}
		return fmt.Sprintf("%s That's only %s", prefix, EmoteGoldMedal.Name)
	case MedalSilver:
		return fmt.Sprintf("%s That's only %s", prefix, EmoteSilverMedal.Name)
	case MedalBronze:
		return fmt.Sprintf("%s That's only %s", prefix, EmoteBronzeMedal.Name)
	default:
		return fmt.Sprintf("%s Not even %s %s", prefix, EmoteBronzeMedal.Name, EmotePepePoint.Name)
	}
}

// WinnerMessage announces who guessed best.
func WinnerMessage(winners []Guess, pb Time) string {
	if len(winners) == 0 {
		return "But no chatter participated " + EmoteReallyFuckingMad.Name
	}

	best := winners[0]
	diff := pb - best.Time
	sign := "-"
	if diff < 0 {
		sign = "+"
	}
	diffText := sign + diff.Abs().String()

	names := make([]string, len(winners))
	for i, w := range winners {
		names[i] = Mention(w.DisplayName)
	}
	mentions := strings.Join(names, ", ")

	if diff == 0 {
		who := mentions + " the " + EmoteGigachad.Name
		if len(winners) > 1 {
			who += " 's"
		}
		return fmt.Sprintf("%s guessed it correctly! %s (+%d points)", who, EmoteBwoah.Name, PointsPerfectGuess)
	}
	return fmt.Sprintf("Nobody guessed it correctly but %s guessed %s (%s) %s (+%d point)",
		mentions, best.Time, diffText, EmoteOk.Name, PointsClosestGuess)
}
