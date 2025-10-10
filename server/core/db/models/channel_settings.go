package models

import "github.com/SoWieMarkus/wheelgpt/core/db"

type ChannelSettings struct {
	// Unique identifier for the channel
	ID string `db:"id,primarykey" json:"id"`

	// Time in seconds before a guess result is revealed in the chat
	GuessDelayTime int `db:"guess_delay_time" json:"guessDelayTime"`
	// Indicates if the bot should remain active when the stream is offline
	BotActiveWhenStreamOffline bool `db:"bot_active_when_stream_offline" json:"botActiveWhenStreamOffline"`
	// Indicates that the channel is publicly shared on the wheelgpt website
	Public bool `db:"public" json:"public"`
}

func (ChannelSettings) TableName() string {
	return "channel_settings"
}

func (ChannelSettings) Indexes() []db.Index {
	return []db.Index{}
}
