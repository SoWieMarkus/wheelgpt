package commands

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

// ResultCmd is the manual !result trigger (MOD only).
// The same core logic is reused by the plugin webhook path via HandleNewPB.
type ResultCmd struct{ Store *db.Store }

func (c *ResultCmd) Names() []string    { return []string{"result"} }
func (c *ResultCmd) Level() bot.AccessLevel { return bot.AccessMod }

func (c *ResultCmd) Execute(ctx *bot.CmdContext) (string, error) {
	if len(ctx.Args) == 0 {
		return fmt.Sprintf("%s %s", bot.MentionUser(ctx.User.DisplayName), bot.ExampleFormat), nil
	}
	t, ok := bot.ParseTMTime(ctx.Args[0])
	if !ok {
		return fmt.Sprintf("%s smh granadyy mods are all degens. Wrong format you idiot. %s",
			bot.MentionUser(ctx.User.DisplayName), bot.ExampleFormat), nil
	}
	msg, err := HandleNewPB(context.Background(), c.Store, ctx.ChannelID, t)
	if err != nil {
		return "", err
	}
	return msg, nil
}

// HandleNewPB is the shared path used by both !result and the plugin webhook.
func HandleNewPB(ctx context.Context, store *db.Store, channelID string, t bot.TMTime) (string, error) {
	guesses, err := store.GetGuessesByChannel(ctx, channelID)
	if err != nil {
		return "", err
	}

	row, _ := store.GetTrackmaniaMap(ctx, channelID)
	var mapData *bot.MapData
	if row.ChannelID != "" {
		mapData = dbRowToMapData(row)
	}

	winners := bot.EvaluateGuesses(guesses, t)

	if err := store.DeleteGuessesByChannel(ctx, channelID); err != nil {
		slog.Error("delete guesses", "channel", channelID, "err", err)
	}
	if err := bot.UpdateLeaderboard(ctx, store.Queries, winners, t); err != nil {
		slog.Error("update leaderboard", "channel", channelID, "err", err)
	}

	return bot.BuildResultMessage(mapData, t, winners), nil
}

func dbRowToMapData(row db.TrackmaniaMap) *bot.MapData {
	return &bot.MapData{
		ChannelID:    row.ChannelID,
		UID:          row.Uid,
		Name:         row.Name,
		Author:       row.Author,
		AuthorTime:   row.AuthorTime,
		GoldTime:     row.GoldTime,
		SilverTime:   row.SilverTime,
		BronzeTime:   row.BronzeTime,
		ChampionTime: row.ChampionTime,
		TmxID:        row.TmxID,
		WorldRecord:  row.WorldRecord,
	}
}
