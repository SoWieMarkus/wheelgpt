package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/testlib"
)

func TestTrackmaniaMap_Insert(t *testing.T) {
	tests := []struct {
		Name             string
		Map              *TrackmaniaMap
		NextMap          *TrackmaniaMap
		CanInsertNextMap bool
	}{
		{
			Name: "Can't insert duplicate TrackmaniaMap for same ChannelID",
			Map: &TrackmaniaMap{
				ChannelID:  "test_channel",
				UID:        "test_uid",
				Name:       "Test Map",
				Author:     "Test Author",
				AuthorTime: 1234567890,
				GoldTime:   1234567891,
				SilverTime: 1234567892,
				BronzeTime: 1234567893,
			},
			NextMap: &TrackmaniaMap{
				ChannelID:  "test_channel",
				UID:        "test_uid2",
				Name:       "Test Map2",
				Author:     "Test Author2",
				AuthorTime: 1234567890,
				GoldTime:   1234567891,
				SilverTime: 1234567892,
				BronzeTime: 1234567893,
			},
			CanInsertNextMap: false,
		},
		{
			Name: "Can insert duplicate TrackmaniaMap with same UID but different ChannelID",
			Map: &TrackmaniaMap{
				ChannelID:  "test_channel2",
				UID:        "test_uid",
				Name:       "Test Map",
				Author:     "Test Author",
				AuthorTime: 1234567890,
				GoldTime:   1234567891,
				SilverTime: 1234567892,
				BronzeTime: 1234567893,
			},
			NextMap: &TrackmaniaMap{
				ChannelID:  "test_channel",
				UID:        "test_uid",
				Name:       "Test Map",
				Author:     "Test Author",
				AuthorTime: 1234567890,
				GoldTime:   1234567891,
				SilverTime: 1234567892,
				BronzeTime: 1234567893,
			},
			CanInsertNextMap: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			sqlite := testlib.SQLite(t)
			testDB := db.Database{DbMap: sqlite.DbMap}
			defer testDB.Close()
			defer sqlite.Close()

			// Create dependency tables
			if err := testDB.CreateTable(
				testDB.AddTable(TrackmaniaMap{}),
			); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if err := testDB.Insert(tt.Map); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			err := testDB.Insert(tt.NextMap)
			if tt.CanInsertNextMap != (err == nil) {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestTrackmaniaMap_TrackmaniaIOLink(t *testing.T) {
	trackmaniaMap := &TrackmaniaMap{
		ChannelID: "test_channel",
		UID:       "test_uid",
	}
	expectedLink := "https://trackmania.io/#/map/test_uid"
	if trackmaniaMap.TrackmaniaIOLink() != expectedLink {
		t.Errorf("expected %s, got %s", expectedLink, trackmaniaMap.TrackmaniaIOLink())
	}
}

func TestTrackmaniaMap_TrackmaniaExchange(t *testing.T) {
	tmxId := int64(12345)
	expectedLink := "https://trackmania.exchange/maps/12345"
	tests := []struct {
		Name         string
		Map          *TrackmaniaMap
		ExpectedLink *string
	}{
		{
			Name: "With TmxID",
			Map: &TrackmaniaMap{
				ChannelID: "test_channel",
				UID:       "test_uid",
				TmxID:     &tmxId,
			},
			ExpectedLink: &expectedLink,
		},
		{
			Name: "Without TmxID",
			Map: &TrackmaniaMap{
				ChannelID: "test_channel",
				UID:       "test_uid",
				TmxID:     nil,
			},
			ExpectedLink: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.Map.TrackmaniaExchangeLink(), tt.ExpectedLink) {
				t.Errorf("expected %v, got %v", tt.ExpectedLink, tt.Map.TrackmaniaExchangeLink())
			}
		})
	}
}

func TestTrackmaniaMap_JsonMarshalling(t *testing.T) {
	tests := []struct {
		Name     string
		Map      *TrackmaniaMap
		Expected string
	}{
		{
			Name: "Without ChampionTime and TmxID returns null for both fields",
			Map: &TrackmaniaMap{
				ChannelID:  "test_channel",
				UID:        "test_uid",
				Name:       "Test Map",
				Author:     "Test Author",
				AuthorTime: 1234567890,
				GoldTime:   1234567891,
				SilverTime: 1234567892,
				BronzeTime: 1234567893,
			},
			Expected: `{"channelId":"test_channel","uid":"test_uid","name":"Test Map","author":"Test Author","authorTime":1234567890,"goldTime":1234567891,"silverTime":1234567892,"bronzeTime":1234567893,"championTime":null,"tmxId":null}`,
		},
		{
			Name: "With ChampionTime and TmxID returns values for both fields",
			Map: &TrackmaniaMap{
				ChannelID:    "test_channel",
				UID:          "test_uid",
				Name:         "Test Map",
				Author:       "Test Author",
				AuthorTime:   1234567890,
				GoldTime:     1234567891,
				SilverTime:   1234567892,
				BronzeTime:   1234567893,
				ChampionTime: testlib.Ptr(int64(1234567894)),
				TmxID:        testlib.Ptr(int64(12345)),
			},
			Expected: `{"channelId":"test_channel","uid":"test_uid","name":"Test Map","author":"Test Author","authorTime":1234567890,"goldTime":1234567891,"silverTime":1234567892,"bronzeTime":1234567893,"championTime":1234567894,"tmxId":12345}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.Map)
			if err != nil {
				t.Fatalf("expected no error during JSON marshalling, got %v", err)
			}

			if string(jsonData) != tt.Expected {
				t.Fatalf("expected JSON %q, got %q", tt.Expected, jsonData)
			}
		})
	}
}
