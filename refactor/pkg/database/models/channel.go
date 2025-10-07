package models

import db "github.com/SoWieMarkus/wheelgpt/pkg/database"

type Channel struct {
	// Unique identifier for the channel
	ID string `db:"id,primarykey" json:"id"`
	// Unique login name for the channel
	Login string `db:"login" json:"login"`
	// Seed used for generating the plugin token for the channel
	Seed string `db:"seed" json:"seed"`
	// Indicates if the channel is live
	IsLive bool `db:"is_live" json:"isLive"`
	// Time in seconds before a guess result is revealed in the chat
	GuessDelayTime int `db:"guess_delay_time" json:"guessDelayTime"`
	// Indicates if the bot should remain active when the stream is offline
	BotActiveWhenStreamOffline bool `db:"bot_active_when_stream_offline" json:"botActiveWhenStreamOffline"`
	// Indicates that the channel is public shared on the wheelgpt website
	Public bool `db:"public" json:"public"`
}

func (Channel) TableName() string {
	return "channels"
}

func (Channel) Indexes() []db.Index {
	return []db.Index{}
}
