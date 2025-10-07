package wheelgpt

import (
	db "github.com/SoWieMarkus/wheelgpt/pkg/database"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch"
)

type WheelGPT struct {
	twitch.Bot
	DB db.Database
}

func (bot *WheelGPT) Initialize() error {
	return nil
}
