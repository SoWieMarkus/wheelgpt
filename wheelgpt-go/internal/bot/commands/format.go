package commands

import (
	"github.com/SoWieMarkus/wheelgpt/internal/bot"
)

type FormatCmd struct{}

func (c *FormatCmd) Names() []string    { return []string{"format"} }
func (c *FormatCmd) Level() bot.AccessLevel { return bot.AccessUser }

func (c *FormatCmd) Execute(ctx *bot.CmdContext) (string, error) {
	return bot.ExampleFormat, nil
}
