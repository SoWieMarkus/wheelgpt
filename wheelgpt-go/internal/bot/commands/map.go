package commands

import (
	"context"
	"fmt"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

type MapCmd struct{ Store *db.Store }

func (c *MapCmd) Names() []string    { return []string{"map"} }
func (c *MapCmd) Level() bot.AccessLevel { return bot.AccessUser }

func (c *MapCmd) Execute(ctx *bot.CmdContext) (string, error) {
	row, err := c.Store.GetTrackmaniaMap(context.Background(), ctx.ChannelID)
	if err != nil || row.ChannelID == "" {
		return fmt.Sprintf("%s No map data available.", bot.MentionUser(ctx.User.DisplayName)), nil
	}
	m := dbRowToMapData(row)
	return m.String(), nil
}
