package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

func (s *Server) handleEventSub(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}

	if !verifyEventSubSignature(r, body, s.cfg.EventSubSecret) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var payload struct {
		Challenge    string `json:"challenge"`
		Subscription struct {
			Type string `json:"type"`
		} `json:"subscription"`
		Event struct {
			BroadcasterUserLogin string `json:"broadcaster_user_login"`
			BroadcasterUserID    string `json:"broadcaster_user_id"`
		} `json:"event"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	msgType := r.Header.Get("Twitch-Eventsub-Message-Type")
	switch msgType {
	case "webhook_callback_verification":
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, payload.Challenge)

	case "notification":
		ctx := r.Context()
		channelID := payload.Event.BroadcasterUserID
		isOnline := payload.Subscription.Type == "stream.online"

		live := int64(0)
		if isOnline {
			live = 1
		}
		if err := s.store.SetChannelLive(ctx, db.SetChannelLiveParams{
			IsLive: live,
			ID:     channelID,
		}); err != nil {
			slog.Error("eventsub: set live", "err", err)
		}

		if !isOnline {
			ch, err := s.store.GetChannel(ctx, channelID)
			if err == nil && !ch.BotActiveWhenOffline.Bool {
				_ = s.store.DeleteTrackmaniaMap(ctx, channelID)
				_ = s.store.DeleteTrackmaniaRoom(ctx, channelID)
				_ = s.store.DeleteGuessesByChannel(ctx, channelID)
			}
		}
		w.WriteHeader(http.StatusNoContent)

	case "revocation":
		slog.Warn("eventsub subscription revoked",
			"type", payload.Subscription.Type,
			"channel", payload.Event.BroadcasterUserID)
		w.WriteHeader(http.StatusNoContent)

	default:
		w.WriteHeader(http.StatusNoContent)
	}
}

func verifyEventSubSignature(r *http.Request, body []byte, secret string) bool {
	msgID := r.Header.Get("Twitch-Eventsub-Message-Id")
	timestamp := r.Header.Get("Twitch-Eventsub-Message-Timestamp")
	sig := r.Header.Get("Twitch-Eventsub-Message-Signature")

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(msgID + timestamp))
	mac.Write(body)
	expected := fmt.Sprintf("sha256=%x", mac.Sum(nil))

	return hmac.Equal([]byte(expected), []byte(sig))
}
