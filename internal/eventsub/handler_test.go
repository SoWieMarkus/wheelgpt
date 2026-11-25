package eventsub_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/eventsub"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/twitch"
)

const secret = "0123456789abcdef"

type recorder struct {
	online, offline []string
	chat            []bot.Message
	revoked         []twitch.Subscription
}

func (r *recorder) StreamOnline(_ context.Context, id string) error {
	r.online = append(r.online, id)
	return nil
}

func (r *recorder) StreamOffline(_ context.Context, id string) error {
	r.offline = append(r.offline, id)
	return nil
}

func (r *recorder) ChatMessage(_ context.Context, m bot.Message) error {
	r.chat = append(r.chat, m)
	return nil
}

func (r *recorder) Revoked(_ context.Context, s twitch.Subscription) error {
	r.revoked = append(r.revoked, s)
	return nil
}

var now = time.Date(2026, 11, 25, 12, 0, 0, 0, time.UTC)

func newHandler() (*eventsub.Handler, *recorder) {
	rec := &recorder{}
	h := eventsub.NewHandler(secret, rec)
	h.Now = func() time.Time { return now }
	return h, rec
}

func send(h http.Handler, id, typ string, ts time.Time, body string, sig string) *httptest.ResponseRecorder {
	stamp := ts.Format(time.RFC3339Nano)
	if sig == "" {
		sig = eventsub.Sign(secret, id, stamp, []byte(body))
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/twitch/eventsub", strings.NewReader(body))
	req.Header.Set("Twitch-Eventsub-Message-Id", id)
	req.Header.Set("Twitch-Eventsub-Message-Timestamp", stamp)
	req.Header.Set("Twitch-Eventsub-Message-Signature", sig)
	req.Header.Set("Twitch-Eventsub-Message-Type", typ)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestChallenge(t *testing.T) {
	h, _ := newHandler()
	w := send(h, "m1", eventsub.TypeVerification, now, `{"challenge":"pogchamp","subscription":{"type":"stream.online"}}`, "")
	if w.Code != http.StatusOK || w.Body.String() != "pogchamp" || w.Header().Get("Content-Type") != "text/plain" {
		t.Fatalf("got %d %q", w.Code, w.Body.String())
	}
}

func TestSignatureAndReplay(t *testing.T) {
	h, rec := newHandler()
	body := `{"subscription":{"type":"stream.online"},"event":{"broadcaster_user_id":"100"}}`
	if w := send(h, "m1", eventsub.TypeNotification, now, body, "sha256=00"); w.Code != http.StatusForbidden {
		t.Errorf("bad signature: %d", w.Code)
	}
	if w := send(h, "m2", eventsub.TypeNotification, now.Add(-11*time.Minute), body, ""); w.Code != http.StatusForbidden {
		t.Errorf("old message: %d", w.Code)
	}
	if w := send(h, "m3", eventsub.TypeNotification, now, body, ""); w.Code != http.StatusNoContent {
		t.Errorf("valid: %d", w.Code)
	}
	if w := send(h, "m3", eventsub.TypeNotification, now, body, ""); w.Code != http.StatusNoContent {
		t.Errorf("duplicate: %d", w.Code)
	}
	if len(rec.online) != 1 || rec.online[0] != "100" {
		t.Errorf("online events = %v (duplicates must be dropped)", rec.online)
	}
}

func TestChatMessage(t *testing.T) {
	h, rec := newHandler()
	body := `{"subscription":{"type":"channel.chat.message"},"event":{
		"broadcaster_user_id":"100","broadcaster_user_login":"streamer",
		"chatter_user_id":"7","chatter_user_login":"viewer","chatter_user_name":"Viewer",
		"message":{"text":"!guess 1.000"},
		"badges":[{"set_id":"subscriber","id":"12"},{"set_id":"vip","id":"1"}]}}`
	if w := send(h, "c1", eventsub.TypeNotification, now, body, ""); w.Code != http.StatusNoContent {
		t.Fatalf("status %d", w.Code)
	}
	if len(rec.chat) != 1 {
		t.Fatalf("chat = %v", rec.chat)
	}
	m := rec.chat[0]
	if m.BroadcasterID != "100" || m.ChatterName != "Viewer" || m.Text != "!guess 1.000" || m.Access != bot.VIP {
		t.Errorf("message = %+v", m)
	}
}

func TestRevocation(t *testing.T) {
	h, rec := newHandler()
	body := `{"subscription":{"type":"channel.chat.message","status":"authorization_revoked","condition":{"broadcaster_user_id":"100"}}}`
	if w := send(h, "r1", eventsub.TypeRevocation, now, body, ""); w.Code != http.StatusNoContent {
		t.Fatalf("status %d", w.Code)
	}
	if len(rec.revoked) != 1 || rec.revoked[0].Condition["broadcaster_user_id"] != "100" {
		t.Errorf("revoked = %+v", rec.revoked)
	}
}
