package models

import "github.com/SoWieMarkus/wheelgpt/core/db"

// Guess represents a user's guess in a channel.
type Guess struct {
	// Channel where the guess was made
	ChannelID string `db:"channel_id,primarykey" json:"channelId"`
	// User who made the guess
	UserID string `db:"user_id,primarykey" json:"userId"`
	// The guess made by the user
	Time int64 `db:"time" json:"time"`
	// User's display name at the time of the guess
	DisplayName string `db:"display_name" json:"displayName"`
}

func (Guess) TableName() string {
	return "guesses"
}

func (Guess) Indexes() []db.Index {
	return []db.Index{
		{Name: "idx_guess_channel_user_id", ColumnNames: []string{"channel_id", "user_id"}},
	}
}
