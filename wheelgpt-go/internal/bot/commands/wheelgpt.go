package commands

import (
	"github.com/SoWieMarkus/wheelgpt/internal/bot"
)

type WheelGPTCmd struct{}

func (c *WheelGPTCmd) Names() []string    { return []string{"wheelgpt", "wgpt"} }
func (c *WheelGPTCmd) Level() bot.AccessLevel { return bot.AccessUser }

func (c *WheelGPTCmd) Execute(ctx *bot.CmdContext) (string, error) {
	return "WheelGPT — guess the streamer's PB time! Use !guess <time> to participate. More info: https://wheelgpt.dev", nil
}
