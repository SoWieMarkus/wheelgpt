package commands

import (
	"context"
	"fmt"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

type RoomCmd struct{ Store *db.Store }

func (c *RoomCmd) Names() []string    { return []string{"room"} }
func (c *RoomCmd) Level() bot.AccessLevel { return bot.AccessUser }

func (c *RoomCmd) Execute(ctx *bot.CmdContext) (string, error) {
	row, err := c.Store.GetTrackmaniaRoom(context.Background(), ctx.ChannelID)
	if err != nil || row.ChannelID == "" {
		return fmt.Sprintf("%s No room data available.", bot.MentionUser(ctx.User.DisplayName)), nil
	}
	return fmt.Sprintf("Current room: %s (%d/%d players)", row.Name, row.NumberOfPlayers, row.MaxPlayers), nil
}
