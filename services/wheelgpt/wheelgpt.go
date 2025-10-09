package wheelgpt

import (
	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch"
)

type WheelGPT struct {
	twitch.Bot
	DB db.Database
}

func (bot *WheelGPT) Initialize() error {
	return nil
}
