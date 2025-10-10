package wheelgpt

import (
	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt/twitch"
)

type WheelGPT struct {
	twitch.Bot
	DB db.Database
}

func NewInstance(db db.Database) WheelGPT {
	return WheelGPT{DB: db}
}

func (bot *WheelGPT) Initialize() error {
	// Load all channels from the database
	// Update the details of all channels
	// Update the current live status of all channels
	// Join all the channels
	return nil
}
