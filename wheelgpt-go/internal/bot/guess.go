package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

const (
	PointsPerfect = 5
	PointsClosest = 1
)

type GuessRow struct {
	ChannelID   string
	UserID      string
	DisplayName string
	Time        int64
}

func EvaluateGuesses(guesses []db.Guess, result TMTime) []db.Guess {
	var (
		minDiff int64 = -1
		winners []db.Guess
	)
	for _, g := range guesses {
		diff := result.Diff(TMTime(g.Time))
		if diff < 0 {
			diff = -diff
		}
		switch {
		case minDiff < 0 || diff < minDiff:
			minDiff = diff
			winners = []db.Guess{g}
		case diff == minDiff:
			winners = append(winners, g)
		}
	}
	return winners
}

func UpdateLeaderboard(ctx context.Context, q *db.Queries, winners []db.Guess, result TMTime) error {
	if len(winners) == 0 {
		return nil
	}
	perfect := TMTime(winners[0].Time) == result
	points := int64(PointsClosest)
	perfectCount := int64(0)
	if perfect {
		points = PointsPerfect
		perfectCount = 1
	}
	for _, w := range winners {
		if err := q.UpsertLeaderboard(ctx, db.UpsertLeaderboardParams{
			ChannelID:          w.ChannelID,
			UserID:             w.UserID,
			DisplayName:        w.DisplayName,
			Points:             points,
			PerfectGuessCount:  perfectCount,
		}); err != nil {
			return err
		}
	}
	return nil
}

func BuildResultMessage(m *MapData, newPB TMTime, winners []db.Guess) string {
	pbMsg := buildPBMessage(m, newPB)
	guessMsg := buildWinnerMessage(winners, newPB)
	return pbMsg + " " + guessMsg
}

func buildPBMessage(m *MapData, t TMTime) string {
	if m == nil {
		return fmt.Sprintf("%s I got a new PB but I didn't know you are on a map?", EmoteYek.Name)
	}
	medal := m.Medal(t)
	diffChampion := m.DiffToMedal(MedalChampion, t)
	diffAuthor := m.DiffToMedal(MedalAuthor, t)

	switch medal {
	case MedalChampion:
		return fmt.Sprintf("NEW PERSONAL BEST %s %s That's Champion Medal %s %s",
			t, EmoteDinkDonk.Name, EmoteChampionMedal.Name, EmoteBwoah.Name)
	case MedalAuthor:
		if diffChampion < 10 {
			return fmt.Sprintf("NEW PERSONAL BEST %s %s HAHAHAHAH HE MISSED CHAMPION BY %ss %s %s %s",
				t, EmoteDinkDonk.Name, TMTime(-diffChampion), EmoteOmegalul.Name, EmoteICant.Name, EmotePepePoint.Name)
		}
		return fmt.Sprintf("NEW PERSONAL BEST %s %s That's %s %s",
			t, EmoteDinkDonk.Name, EmoteAuthorMedal.Name, EmoteOK.Name)
	case MedalGold:
		if diffAuthor < 10 {
			return fmt.Sprintf("NEW PERSONAL BEST %s %s HAHAHAHAH HE MISSED AUTHOR BY %ss %s %s %s",
				t, EmoteDinkDonk.Name, TMTime(-diffAuthor), EmoteOmegalul.Name, EmoteICant.Name, EmotePepePoint.Name)
		}
		return fmt.Sprintf("NEW PERSONAL BEST %s %s That's only %s",
			t, EmoteDinkDonk.Name, EmoteGoldMedal.Name)
	case MedalSilver:
		return fmt.Sprintf("NEW PERSONAL BEST %s %s That's only %s",
			t, EmoteDinkDonk.Name, EmoteSilverMedal.Name)
	case MedalBronze:
		return fmt.Sprintf("NEW PERSONAL BEST %s %s That's only %s",
			t, EmoteDinkDonk.Name, EmoteBronzeMedal.Name)
	default:
		return fmt.Sprintf("NEW PERSONAL BEST %s %s Not even %s %s",
			t, EmoteDinkDonk.Name, EmoteBronzeMedal.Name, EmotePepePoint.Name)
	}
}

func buildWinnerMessage(winners []db.Guess, result TMTime) string {
	if len(winners) == 0 {
		return fmt.Sprintf("But no chatter participated %s", EmoteReallyFuckingMad.Name)
	}
	best := winners[0]
	bestTime := TMTime(best.Time)
	diff := result.Diff(bestTime)
	sign := "-"
	if diff < 0 {
		sign = "+"
		diff = -diff
	}
	diffStr := fmt.Sprintf("%s%s", sign, TMTime(diff))
	perfect := result == bestTime

	if len(winners) == 1 {
		if perfect {
			return fmt.Sprintf("%s the %s guessed it correctly! %s (+%d points)",
				MentionUser(best.DisplayName), EmoteGigachad.Name, EmoteBwoah.Name, PointsPerfect)
		}
		return fmt.Sprintf("Nobody guessed it correctly but %s guessed %s (%s) %s (+%d point)",
			MentionUser(best.DisplayName), bestTime, diffStr, EmoteOK.Name, PointsClosest)
	}

	names := make([]string, len(winners))
	for i, w := range winners {
		names[i] = MentionUser(w.DisplayName)
	}
	joined := strings.Join(names, ", ")
	if perfect {
		return fmt.Sprintf("%s the %s 's guessed it correctly! %s (+%d points)",
			joined, EmoteGigachad.Name, EmoteBwoah.Name, PointsPerfect)
	}
	return fmt.Sprintf("Nobody guessed it correctly but %s guessed %s (%s) %s (+%d point)",
		joined, bestTime, diffStr, EmoteOK.Name, PointsClosest)
}
