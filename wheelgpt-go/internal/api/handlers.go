package api

import (
	"encoding/json"
	"net/http"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/bot/commands"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

func (s *Server) handleNewPB(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Time int64 `json:"time"`
	}
	if err := decodeJSON(r, &body); err != nil || body.Time <= 0 {
		jsonError(w, http.StatusBadRequest, "invalid time")
		return
	}

	ch, err := s.store.GetChannel(r.Context(), channelIDFrom(r.Context()))
	if err != nil {
		jsonError(w, http.StatusNotFound, "channel not found")
		return
	}

	msg, err := commands.HandleNewPB(r.Context(), s.store, ch.ID, bot.TMTime(body.Time))
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "pb handling failed")
		return
	}

	// Trigger the bot to say the message after the configured delay
	s.bot.NotifyNewPB(ch.Login, body.Time)

	writeJSON(w, http.StatusOK, map[string]string{"message": msg})
}

func (s *Server) handleUpsertMap(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UID          string `json:"uid"`
		Name         string `json:"name"`
		Author       string `json:"author"`
		AuthorTime   int64  `json:"authorTime"`
		GoldTime     int64  `json:"goldTime"`
		SilverTime   int64  `json:"silverTime"`
		BronzeTime   int64  `json:"bronzeTime"`
		ChampionTime int64  `json:"championTime"`
		TmxID        *int64 `json:"tmxId"`
		WorldRecord  *int64 `json:"worldRecord"`
	}
	if err := decodeJSON(r, &body); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid body")
		return
	}
	channelID := channelIDFrom(r.Context())
	if err := s.store.UpsertTrackmaniaMap(r.Context(), db.UpsertTrackmaniaMapParams{
		ChannelID:    channelID,
		Uid:          body.UID,
		Name:         body.Name,
		Author:       body.Author,
		AuthorTime:   body.AuthorTime,
		GoldTime:     body.GoldTime,
		SilverTime:   body.SilverTime,
		BronzeTime:   body.BronzeTime,
		ChampionTime: body.ChampionTime,
		TmxID:        body.TmxID,
		WorldRecord:  body.WorldRecord,
	}); err != nil {
		jsonError(w, http.StatusInternalServerError, "db error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleUpsertRoom(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Login           string `json:"login"`
		Name            string `json:"name"`
		NumberOfPlayers int64  `json:"numberOfPlayers"`
		MaxPlayers      int64  `json:"maxPlayers"`
	}
	if err := decodeJSON(r, &body); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid body")
		return
	}
	channelID := channelIDFrom(r.Context())
	if err := s.store.UpsertTrackmaniaRoom(r.Context(), db.UpsertTrackmaniaRoomParams{
		ChannelID:       channelID,
		Login:           body.Login,
		Name:            body.Name,
		NumberOfPlayers: body.NumberOfPlayers,
		MaxPlayers:      body.MaxPlayers,
	}); err != nil {
		jsonError(w, http.StatusInternalServerError, "db error")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetChannel(w http.ResponseWriter, r *http.Request) {
	ch, err := s.store.GetChannel(r.Context(), channelIDFrom(r.Context()))
	if err != nil {
		jsonError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, ch)
}

func (s *Server) handleUpdateChannel(w http.ResponseWriter, r *http.Request) {
	var body struct {
		GuessDelayTime       float64 `json:"guessDelayTime"`
		BotActiveWhenOffline bool    `json:"botActiveWhenOffline"`
		UsagePublic          bool    `json:"usagePublic"`
	}
	if err := decodeJSON(r, &body); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid body")
		return
	}
	channelID := channelIDFrom(r.Context())
	inactive := int64(0)
	if body.BotActiveWhenOffline {
		inactive = 1
	}
	public := int64(0)
	if body.UsagePublic {
		public = 1
	}
	if err := s.store.UpdateChannelSettings(r.Context(), db.UpdateChannelSettingsParams{
		GuessDelayTime:       body.GuessDelayTime,
		BotActiveWhenOffline: inactive,
		UsagePublic:          public,
		ID:                   channelID,
	}); err != nil {
		jsonError(w, http.StatusInternalServerError, "db error")
		return
	}
	_ = s.bot.Reload(r.Context(), channelID)
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleGetLeaderboard(w http.ResponseWriter, r *http.Request) {
	rows, err := s.store.GetLeaderboard(r.Context(), db.GetLeaderboardParams{
		ChannelID: channelIDFrom(r.Context()),
		Limit:     50,
		Offset:    0,
	})
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func botCfg(ch db.Channel) bot.ChannelConfig {
	return bot.ChannelConfig{
		ID:                   ch.ID,
		Login:                ch.Login,
		DisplayName:          ch.DisplayName,
		GuessDelayTime:       ch.GuessDelayTime,
		BotActiveWhenOffline: ch.BotActiveWhenOffline != 0,
	}
}

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
