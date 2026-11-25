package app

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/SoWieMarkus/wheelgpt/internal/platform/metrics"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/store/storetest"
)

const (
	legacySecret   = "legacy-jwt-secret"
	eventSubSecret = "0123456789abcdef"
)

// chat records the messages the bot sends.
type chat struct {
	mu   sync.Mutex
	sent []string
}

func (c *chat) Send(_ context.Context, _, text string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sent = append(c.sent, text)
	return nil
}

func (c *chat) last() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.sent) == 0 {
		return ""
	}
	return c.sent[len(c.sent)-1]
}

type noTMX struct{}

func (noTMX) MapID(context.Context, string) (int64, bool, error) { return 0, false, nil }

type testServer struct {
	t     *testing.T
	base  string
	http  *http.Client
	st    *store.Store
	chat  *chat
	comps components
}

func newTestServer(t *testing.T, mutate func(*Config)) *testServer {
	t.Helper()
	publicURL, _ := url.Parse("http://localhost:5173")
	cfg := Config{
		PublicURL:       publicURL,
		SessionTTL:      24 * time.Hour,
		BotUsername:     "wheelgpt",
		LegacyJWTSecret: legacySecret,
		EventSubSecret:  eventSubSecret,
		DevMode:         true,
	}
	if mutate != nil {
		mutate(&cfg)
	}
	st := storetest.New(t)
	storetest.Channel(t, st, "100", "streamer")
	c := &chat{}
	comps, err := build(cfg, "test", st, metrics.New("test", nil), externals{tmx: noTMX{}, sender: c})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { comps.bot.Stop(); comps.ingame.Wait() })
	srv := httptest.NewServer(comps.router)
	t.Cleanup(srv.Close)
	jar, _ := cookiejar.New(nil)
	return &testServer{t: t, base: srv.URL, http: &http.Client{Jar: jar}, st: st, chat: c, comps: comps}
}

type response struct {
	status int
	header http.Header
	body   map[string]any
	raw    string
}

func (s *testServer) do(method, path string, body any, headers ...string) response {
	s.t.Helper()
	var reader io.Reader
	switch b := body.(type) {
	case nil:
	case string:
		reader = strings.NewReader(b)
	default:
		raw, err := json.Marshal(b)
		if err != nil {
			s.t.Fatal(err)
		}
		reader = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(s.t.Context(), method, s.base+path, reader)
	if err != nil {
		s.t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}
	res, err := s.http.Do(req)
	if err != nil {
		s.t.Fatal(err)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	out := response{status: res.StatusCode, header: res.Header, raw: string(raw)}
	_ = json.Unmarshal(raw, &out.body)
	return out
}

func (s *testServer) expect(r response, status int, code string) {
	s.t.Helper()
	if r.status != status {
		s.t.Fatalf("status = %d, want %d: %s", r.status, status, r.raw)
	}
	if code != "" && r.body["code"] != code {
		s.t.Fatalf("code = %v, want %s: %s", r.body["code"], code, r.raw)
	}
}

// login logs in through the dev login.
func (s *testServer) login(name string) {
	s.t.Helper()
	noRedirect := *s.http
	noRedirect.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	req, _ := http.NewRequestWithContext(s.t.Context(), http.MethodGet, s.base+"/api/v1/auth/dev?login="+name, http.NoBody)
	res, err := noRedirect.Do(req)
	if err != nil {
		s.t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode != http.StatusSeeOther {
		s.t.Fatalf("dev login: %d", res.StatusCode)
	}
}

// chatMessage delivers a signed EventSub chat message to the webhook.
func (s *testServer) chatMessage(userID, name, text string, badges ...string) {
	s.t.Helper()
	b := make([]map[string]string, len(badges))
	for i, id := range badges {
		b[i] = map[string]string{"set_id": id}
	}
	s.event("channel.chat.message", map[string]any{
		"broadcaster_user_id": "100", "chatter_user_id": userID, "chatter_user_name": name,
		"message": map[string]any{"text": text}, "badges": b,
	})
}

var eventIDs atomic.Int64

// event delivers a signed EventSub notification to the webhook.
func (s *testServer) event(typ string, event map[string]any) {
	s.t.Helper()
	body, _ := json.Marshal(map[string]any{"subscription": map[string]any{"type": typ}, "event": event})
	id := fmt.Sprintf("test-message-%d", eventIDs.Add(1))
	ts := time.Now().UTC().Format(time.RFC3339Nano)
	mac := hmac.New(sha256.New, []byte(eventSubSecret))
	mac.Write([]byte(id + ts))
	mac.Write(body)
	r := s.do(http.MethodPost, "/api/v1/twitch/eventsub", string(body),
		"Twitch-Eventsub-Message-Id", id, "Twitch-Eventsub-Message-Timestamp", ts,
		"Twitch-Eventsub-Message-Signature", "sha256="+hex.EncodeToString(mac.Sum(nil)),
		"Twitch-Eventsub-Message-Type", "notification")
	s.expect(r, http.StatusNoContent, "")
}

func legacyToken(t *testing.T) string {
	t.Helper()
	tok, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": "100", "token": "legacy-secret"}).SignedString([]byte(legacySecret))
	if err != nil {
		t.Fatal(err)
	}
	return tok
}

const mapBody = `{"uid":"u1","name":"$fffSpring 07","author":"Nadeo","authorTime":41000,"goldTime":44000,"silverTime":50000,"bronzeTime":60000,"championTime":0}`

func TestPluginWithLegacyToken(t *testing.T) {
	s := newTestServer(t, nil)
	auth := "Bearer " + legacyToken(t)

	storetest.Live(t, s.st, "100")
	r := s.do(http.MethodPut, "/api/v1/plugin/map", mapBody, "Authorization", auth)
	s.expect(r, http.StatusOK, "")
	if r.body["name"] != "Spring 07" {
		t.Fatalf("map = %s", r.raw)
	}

	// A guess through chat, then the PB from the plugin.
	s.chatMessage("7", "Viewer", "!guess 45.000")
	s.chatMessage("8", "Other", "!g 50.000")
	r = s.do(http.MethodPost, "/api/v1/plugin/personal-bests", `{"time":45000}`, "Authorization", auth)
	s.expect(r, http.StatusAccepted, "")
	round, _ := r.body["round"].(map[string]any)
	if round["participants"] != float64(2) || round["points"] != float64(5) {
		t.Fatalf("pb = %s", r.raw)
	}
	if got := s.chat.last(); got != "" {
		t.Fatalf("result announced before the stream delay: %q", got)
	}
	s.comps.bot.Stop() // sends the delayed announcement right away
	if got := s.chat.last(); !strings.Contains(got, "@Viewer the GIGACHAD guessed it correctly!") {
		t.Fatalf("announcement = %q", got)
	}

	r = s.do(http.MethodGet, "/api/v1/channels/streamer/leaderboards/all-time", nil)
	s.expect(r, http.StatusOK, "")
	if !strings.Contains(r.raw, `"displayName":"Viewer"`) || !strings.Contains(r.raw, `"points":5`) {
		t.Fatalf("leaderboard = %s", r.raw)
	}

	// The old plugin: raw token, old path, expects 200.
	r = s.do(http.MethodPost, "/api/trackmania/room", `{"login":"r","name":"Room","numberOfPlayers":3,"maxPlayers":10,"extra":1}`, "Authorization", legacyToken(t))
	s.expect(r, http.StatusOK, "")
	if r.body["message"] != "Room updated successfully." || r.header.Get("Deprecation") == "" {
		t.Fatalf("legacy room = %s %v", r.raw, r.header)
	}
	s.chatMessage("9", "Asker", "!room")
	if got := s.chat.last(); got != "@Asker Room [3/10]" {
		t.Fatalf("!room = %q", got)
	}
}

// Plugin data only counts while the stream is live.
func TestOfflineRejectsPluginData(t *testing.T) {
	s := newTestServer(t, nil)
	auth := "Bearer " + legacyToken(t)

	s.expect(s.do(http.MethodPut, "/api/v1/plugin/map", mapBody, "Authorization", auth), http.StatusConflict, "channel_offline")
	s.expect(s.do(http.MethodPut, "/api/v1/plugin/room", `{"login":"r","name":"Room","numberOfPlayers":1,"maxPlayers":2}`, "Authorization", auth), http.StatusConflict, "channel_offline")
	s.expect(s.do(http.MethodPost, "/api/v1/plugin/personal-bests", `{"time":1}`, "Authorization", auth), http.StatusConflict, "channel_offline")
	s.expect(s.do(http.MethodDelete, "/api/v1/plugin/map", nil, "Authorization", auth), http.StatusConflict, "channel_offline")

	// Old plugin builds retry anything but 200, so they get the old answer.
	r := s.do(http.MethodPost, "/api/trackmania/map", mapBody, "Authorization", legacyToken(t))
	s.expect(r, http.StatusOK, "")
	if r.body["message"] != "Channel is offline." {
		t.Fatalf("legacy offline = %s", r.raw)
	}
	s.expect(s.do(http.MethodGet, "/api/v1/channels/streamer/map", nil), http.StatusNotFound, "map_not_found")

	// Going live accepts data, going offline clears it again.
	s.event("stream.online", map[string]any{"broadcaster_user_id": "100"})
	s.expect(s.do(http.MethodPut, "/api/v1/plugin/map", mapBody, "Authorization", auth), http.StatusOK, "")
	s.chatMessage("7", "Viewer", "!guess 45.000")
	s.event("stream.offline", map[string]any{"broadcaster_user_id": "100"})
	s.expect(s.do(http.MethodGet, "/api/v1/channels/streamer/map", nil), http.StatusNotFound, "map_not_found")
	s.chatMessage("7", "Viewer", "!mg")
	if got := s.chat.last(); !strings.Contains(got, "can't find any guess") {
		t.Fatalf("guesses must be cleared when the stream ends: %q", got)
	}
}

// Dev mode fakes the live status, so the plugin can be tested offline.
func TestDevModeLiveSwitch(t *testing.T) {
	s := newTestServer(t, nil)
	s.login("tester")
	r := s.do(http.MethodGet, "/api/v1/me", nil)
	if r.body["devMode"] != true {
		t.Fatalf("me = %s", r.raw)
	}
	r = s.do(http.MethodPost, "/api/v1/me/plugin-token", nil)
	auth := "Bearer " + r.body["token"].(string)

	s.expect(s.do(http.MethodPut, "/api/v1/plugin/map", mapBody, "Authorization", auth), http.StatusConflict, "channel_offline")
	r = s.do(http.MethodPut, "/api/v1/me/live", nil)
	s.expect(r, http.StatusOK, "")
	if r.body["isLive"] != true {
		t.Fatalf("go live = %s", r.raw)
	}
	s.expect(s.do(http.MethodPut, "/api/v1/plugin/map", mapBody, "Authorization", auth), http.StatusOK, "")

	s.expect(s.do(http.MethodDelete, "/api/v1/me/live", nil), http.StatusOK, "")
	ref := s.do(http.MethodGet, "/api/v1/me", nil).body["channel"].(map[string]any)["login"].(string)
	s.expect(s.do(http.MethodGet, "/api/v1/channels/"+ref+"/map", nil), http.StatusNotFound, "map_not_found")

	// Without dev mode there is no switch.
	off := newTestServer(t, func(c *Config) { c.DevMode = false })
	issued, err := off.comps.auth.DevLogin(t.Context(), "tester")
	if err != nil {
		t.Fatal(err)
	}
	cookie := "wheelgpt_session=" + issued.Token
	off.expect(off.do(http.MethodGet, "/api/v1/me", nil, "Cookie", cookie), http.StatusOK, "")
	off.expect(off.do(http.MethodPut, "/api/v1/me/live", nil, "Cookie", cookie), http.StatusNotFound, "route_not_found")
	off.expect(off.do(http.MethodGet, "/api/v1/auth/dev?login=x", nil), http.StatusNotFound, "route_not_found")
}

func TestPluginAuthAndValidation(t *testing.T) {
	s := newTestServer(t, func(c *Config) { c.PluginRateLimit = 5 })
	s.expect(s.do(http.MethodPut, "/api/v1/plugin/map", mapBody), http.StatusUnauthorized, "invalid_plugin_token")
	s.expect(s.do(http.MethodGet, "/api/v1/plugin/status", nil, "Authorization", "Bearer wgpt_nope"), http.StatusUnauthorized, "invalid_plugin_token")

	auth := "Bearer " + legacyToken(t)
	s.expect(s.do(http.MethodPut, "/api/v1/plugin/map", `{"uid":`, "Authorization", auth), http.StatusBadRequest, "invalid_body")
	s.expect(s.do(http.MethodPut, "/api/v1/plugin/map", strings.Replace(mapBody, `"uid":"u1"`, `"uid":"u1","extra":1`, 1), "Authorization", auth), http.StatusBadRequest, "invalid_body")
	r := s.do(http.MethodPut, "/api/v1/plugin/map", strings.Replace(mapBody, `"authorTime":41000`, `"authorTime":-1`, 1), "Authorization", auth)
	s.expect(r, http.StatusBadRequest, "validation_failed")
	if !strings.Contains(r.raw, `"field":"authorTime"`) {
		t.Fatalf("errors = %s", r.raw)
	}

	limited := false
	for range 10 {
		if s.do(http.MethodGet, "/api/v1/plugin/status", nil, "Authorization", auth).status == http.StatusTooManyRequests {
			limited = true
			break
		}
	}
	if !limited {
		t.Fatal("plugin rate limit not applied")
	}
}

func TestDashboard(t *testing.T) {
	s := newTestServer(t, nil)
	s.expect(s.do(http.MethodGet, "/api/v1/me", nil), http.StatusUnauthorized, "unauthenticated")

	s.login("newstreamer")
	r := s.do(http.MethodGet, "/api/v1/me", nil)
	s.expect(r, http.StatusOK, "")
	if r.header.Get("Cache-Control") != "no-store" {
		t.Errorf("Cache-Control = %q", r.header.Get("Cache-Control"))
	}

	r = s.do(http.MethodPatch, "/api/v1/me/settings", map[string]any{"guessDelaySeconds": 7.5})
	s.expect(r, http.StatusOK, "")
	if r.body["guessDelaySeconds"] != 7.5 || r.body["public"] != true {
		t.Fatalf("settings = %s", r.raw)
	}
	s.expect(s.do(http.MethodPatch, "/api/v1/me/settings", map[string]any{"guessDelaySeconds": 100}), http.StatusBadRequest, "validation_failed")

	// Cross-site writes are rejected.
	s.expect(s.do(http.MethodPatch, "/api/v1/me/settings", map[string]any{"public": false},
		"Sec-Fetch-Site", "cross-site", "Origin", "https://evil.example"), http.StatusForbidden, "cross_origin")

	// A new channel never had an old-style token.
	r = s.do(http.MethodGet, "/api/v1/me/plugin-token", nil)
	if r.body["legacyTokenActive"] != false || r.body["exists"] != false {
		t.Fatalf("token info = %s", r.raw)
	}
	r = s.do(http.MethodPost, "/api/v1/me/plugin-token", nil)
	s.expect(r, http.StatusCreated, "")
	token, _ := r.body["token"].(string)
	s.expect(s.do(http.MethodGet, "/api/v1/plugin/status", nil, "Authorization", "Bearer "+token), http.StatusOK, "")

	r = s.do(http.MethodDelete, "/api/v1/me/bot-connection", nil)
	s.expect(r, http.StatusOK, "")
	if r.body["chatSubscription"] != "disabled" {
		t.Fatalf("disconnect = %s", r.raw)
	}

	s.expect(s.do(http.MethodPost, "/api/v1/auth/logout", nil), http.StatusNoContent, "")
	s.expect(s.do(http.MethodGet, "/api/v1/me", nil), http.StatusUnauthorized, "unauthenticated")
}

func TestNewTokenRevokesLegacyToken(t *testing.T) {
	s := newTestServer(t, nil)
	auth := "Bearer " + legacyToken(t)
	s.expect(s.do(http.MethodGet, "/api/v1/plugin/status", nil, "Authorization", auth), http.StatusOK, "")

	token, _, err := s.comps.plugin.CreateToken(t.Context(), "100")
	if err != nil {
		t.Fatal(err)
	}
	s.expect(s.do(http.MethodGet, "/api/v1/plugin/status", nil, "Authorization", "Bearer "+token), http.StatusOK, "")
	s.expect(s.do(http.MethodGet, "/api/v1/plugin/status", nil, "Authorization", auth), http.StatusUnauthorized, "invalid_plugin_token")
}

func TestErrorsAndDocs(t *testing.T) {
	s := newTestServer(t, nil)
	s.expect(s.do(http.MethodGet, "/api/v1/nope", nil), http.StatusNotFound, "route_not_found")
	s.expect(s.do(http.MethodDelete, "/api/v1/commands", nil), http.StatusMethodNotAllowed, "method_not_allowed")
	s.expect(s.do(http.MethodGet, "/api/v1/channels/streamer/leaderboards/2026-13", nil), http.StatusBadRequest, "invalid_period")
	s.expect(s.do(http.MethodGet, "/api/v1/channels?limit=1000", nil), http.StatusBadRequest, "invalid_query")
	s.expect(s.do(http.MethodGet, "/api/v1/channels/unknown", nil), http.StatusNotFound, "channel_not_found")

	r := s.do(http.MethodGet, "/api/v1/commands", nil)
	s.expect(r, http.StatusOK, "")
	if !strings.Contains(r.raw, `"name":"guess"`) || r.header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("commands = %s", r.raw)
	}
	for _, p := range []string{"/healthz", "/readyz", "/api/docs/doc.json"} {
		s.expect(s.do(http.MethodGet, p, nil), http.StatusOK, "")
	}
}
