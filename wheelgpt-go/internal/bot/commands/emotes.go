package commands

import (
	"fmt"
	"strings"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
)

type EmotesCmd struct{}

func (c *EmotesCmd) Names() []string    { return []string{"wgpt-emotes"} }
func (c *EmotesCmd) Level() bot.AccessLevel { return bot.AccessMod }

func (c *EmotesCmd) Execute(ctx *bot.CmdContext) (string, error) {
	names := make([]string, 0, len(bot.AllEmotes()))
	for _, e := range bot.AllEmotes() {
		names = append(names, e.Name)
	}
	return fmt.Sprintf("Available emotes: %s", strings.Join(names, ", ")), nil
}
