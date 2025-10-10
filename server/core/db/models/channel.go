package models

import "github.com/SoWieMarkus/wheelgpt/core/db"

type Channel struct {
	// Unique identifier for the channel
	ID string `db:"id,primarykey" json:"id"`
	// Unique login name for the channel
	Login string `db:"login" json:"login"`
	// Seed used for generating the plugin token for the channel
	Seed string `db:"seed" json:"seed"`
	// Indicates if the channel is live
	IsLive bool `db:"is_live" json:"isLive"`
}

func (Channel) TableName() string {
	return "channels"
}

func (Channel) Indexes() []db.Index {
	return []db.Index{}
}
