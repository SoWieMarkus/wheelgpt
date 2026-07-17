package command

import "github.com/SoWieMarkus/wheelgpt/pkg/bot"

type Guess struct {
}

func (g *Guess) Names() []string {
	return []string{
		"guess",
		"g",
	}
}

func (g *Guess) Execute(ctx bot.CommandContext) (string, error) {
	if len(ctx.Args) == 0 {
		return "Please provide a guess.", nil
	}
}
