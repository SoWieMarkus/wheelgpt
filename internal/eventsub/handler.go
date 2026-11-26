// Package eventsub receives Twitch EventSub webhooks and keeps the
// subscriptions in sync with the channels in the database.
//
// Webhooks (not WebSockets) are the transport: the bot is a server with a
// public https endpoint, webhooks work with an app access token for any
// number of channels, and there is no connection to keep alive.
package eventsub

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/SoWieMarkus/gokit/httpx"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/twitch"
)

// Header names (Go canonicalizes header keys).
const (
	headerID        = "Twitch-Eventsub-Message-Id"
	headerTimestamp = "Twitch-Eventsub-Message-Timestamp"
	headerSignature = "Twitch-Eventsub-Message-Signature"
	headerType      = "Twitch-Eventsub-Message-Type"
)

// Message types.
const (
	TypeVerification = "webhook_callback_verification"
	TypeNotification = "notification"
	TypeRevocation   = "revocation"
)

// MaxAge is how old a message may be before it is rejected as a replay.
const MaxAge = 10 * time.Minute

// Dispatcher reacts to events.
type Dispatcher interface {
	StreamOnline(ctx context.Context, broadcasterID string) error
	StreamOffline(ctx context.Context, broadcasterID string) error
	ChatMessage(ctx context.Context, msg bot.Message) error
	Revoked(ctx context.Context, sub twitch.Subscription) error
}

// Handler is the webhook endpoint.
type Handler struct {
	Secret     string
	Dispatcher Dispatcher
	Now        func() time.Time
	// OnMessage is called for every verified message (metrics), may be nil.
	OnMessage func(messageType, subscriptionType string)

	mu   sync.Mutex
	seen map[string]time.Time
}

// NewHandler returns a webhook handler.
func NewHandler(secret string, d Dispatcher) *Handler {
	return &Handler{Secret: secret, Dispatcher: d, Now: time.Now, seen: map[string]time.Time{}}
}

// Mount registers the webhook endpoint and the callback of the old backend.
//
//	@Summary		Twitch EventSub webhook
//	@Description	Receives stream.online, stream.offline and channel.chat.message. Requests are verified with the Twitch-Eventsub-Message-Signature HMAC, must be younger than 10 minutes and are deduplicated by message id.
//	@Tags			twitch
//	@Accept			json
//	@Success		200
//	@Success		204
//	@Failure		403	{object}	httpx.Problem
//	@Router			/twitch/eventsub [post]
func (h *Handler) Mount(r chi.Router) {
	r.Post("/twitch/eventsub", h.ServeHTTP)
}

// Sign computes the signature Twitch sends for a message.
func Sign(secret, id, timestamp string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(id))
	mac.Write([]byte(timestamp))
	mac.Write(body)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

type envelope struct {
	Challenge    string              `json:"challenge"`
	Subscription twitch.Subscription `json:"subscription"`
	Event        json.RawMessage     `json:"event"`
}

type streamEvent struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
}

type chatEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	ChatterUserID        string `json:"chatter_user_id"`
	ChatterUserLogin     string `json:"chatter_user_login"`
	ChatterUserName      string `json:"chatter_user_name"`
	Message              struct {
		Text string `json:"text"`
	} `json:"message"`
	Badges []struct {
		SetID string `json:"set_id"`
	} `json:"badges"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		httpx.Error(w, r, httpx.BadRequest("invalid_body", "could not read the body"))
		return
	}

	id := r.Header.Get(headerID)
	timestamp := r.Header.Get(headerTimestamp)
	signature := r.Header.Get(headerSignature)
	if id == "" || timestamp == "" || signature == "" ||
		!hmac.Equal([]byte(signature), []byte(Sign(h.Secret, id, timestamp, body))) {
		slog.WarnContext(r.Context(), "eventsub: invalid signature", "message_id", id)
		httpx.Error(w, r, ErrInvalidSignature)
		return
	}
	sent, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil || h.Now().Sub(sent).Abs() > MaxAge {
		httpx.Error(w, r, ErrInvalidSignature)
		return
	}
	if h.duplicate(id) {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		httpx.Error(w, r, httpx.BadRequest("invalid_body", "malformed JSON"))
		return
	}
	messageType := r.Header.Get(headerType)
	if h.OnMessage != nil {
		h.OnMessage(messageType, env.Subscription.Type)
	}

	switch messageType {
	case TypeVerification:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, env.Challenge)
		return
	case TypeRevocation:
		slog.WarnContext(r.Context(), "eventsub: subscription revoked", "type", env.Subscription.Type,
			"status", env.Subscription.Status, "condition", env.Subscription.Condition)
		err = h.Dispatcher.Revoked(r.Context(), env.Subscription)
	case TypeNotification:
		err = h.notification(r.Context(), env)
	default:
		httpx.Error(w, r, httpx.BadRequest("unknown_message_type", "unknown EventSub message type"))
		return
	}
	if err != nil {
		// Answer 2xx anyway: Twitch would retry and eventually revoke the
		// subscription, and a retry rarely fixes an internal error.
		slog.ErrorContext(r.Context(), "eventsub: handling failed", "message_id", id, "type", env.Subscription.Type, "error", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) notification(ctx context.Context, env envelope) error {
	switch env.Subscription.Type {
	case twitch.SubStreamOnline, twitch.SubStreamOffline:
		var ev streamEvent
		if err := json.Unmarshal(env.Event, &ev); err != nil || ev.BroadcasterUserID == "" {
			return errors.New("invalid stream event")
		}
		if env.Subscription.Type == twitch.SubStreamOnline {
			return h.Dispatcher.StreamOnline(ctx, ev.BroadcasterUserID)
		}
		return h.Dispatcher.StreamOffline(ctx, ev.BroadcasterUserID)
	case twitch.SubChatMessage:
		var ev chatEvent
		if err := json.Unmarshal(env.Event, &ev); err != nil || ev.BroadcasterUserID == "" {
			return errors.New("invalid chat event")
		}
		badges := make([]string, len(ev.Badges))
		for i, b := range ev.Badges {
			badges[i] = b.SetID
		}
		access := bot.AccessFromBadges(badges)
		if ev.ChatterUserID == ev.BroadcasterUserID {
			access = bot.Broadcaster
		}
		return h.Dispatcher.ChatMessage(ctx, bot.Message{
			BroadcasterID: ev.BroadcasterUserID,
			ChatterID:     ev.ChatterUserID,
			ChatterName:   ev.ChatterUserName,
			Text:          ev.Message.Text,
			Access:        access,
		})
	default:
		return nil
	}
}

// duplicate records id and reports whether it was seen within MaxAge.
func (h *Handler) duplicate(id string) bool {
	now := h.Now()
	h.mu.Lock()
	defer h.mu.Unlock()
	if t, ok := h.seen[id]; ok && now.Sub(t) < MaxAge {
		return true
	}
	h.seen[id] = now
	if len(h.seen) > 1000 {
		for k, t := range h.seen {
			if now.Sub(t) >= MaxAge {
				delete(h.seen, k)
			}
		}
	}
	return false
}
