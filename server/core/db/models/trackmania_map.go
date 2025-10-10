package models

import (
	"fmt"

	"github.com/SoWieMarkus/wheelgpt/core/db"
)

type TrackmaniaMap struct {
	// Channel where this map is currently played
	ChannelID string `db:"channel_id,primarykey" json:"channelId"`
	// Unique identifier for the map. Multiple channels can play the same map, so this is not a primary key.
	UID string `db:"uid" json:"uid"`
	// Name of the map
	Name string `db:"name" json:"name"`
	// Name of the author of the map
	Author string `db:"author" json:"author"`
	// Trackmania medal times
	AuthorTime int64 `db:"author_time" json:"authorTime"`
	GoldTime   int64 `db:"gold_time" json:"goldTime"`
	SilverTime int64 `db:"silver_time" json:"silverTime"`
	BronzeTime int64 `db:"bronze_time" json:"bronzeTime"`
	// Trackmania champion time from Champion Medal plugin
	ChampionTime *int64 `db:"champion_time" json:"championTime"`
	// ID from Trackmania Exchange, if available, to share the link to the tmx page.
	TmxID *int64 `db:"tmx_id" json:"tmxId"`
}

func (TrackmaniaMap) TableName() string {
	return "trackmania_maps"
}

func (TrackmaniaMap) Indexes() []db.Index {
	return []db.Index{}
}

// Get the link to the Trackmania.io page for this map.
func (m *TrackmaniaMap) TrackmaniaIOLink() string {
	return fmt.Sprintf("https://trackmania.io/#/map/%s", m.UID)
}

// Get the link to the Trackmania Exchange page for this map, if TmxID is set.
func (m *TrackmaniaMap) TrackmaniaExchangeLink() *string {
	if m.TmxID == nil {
		return nil
	}
	link := fmt.Sprintf("https://trackmania.exchange/maps/%d", *m.TmxID)
	return &link
}
