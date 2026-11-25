package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"

	"github.com/SoWieMarkus/wheelgpt/internal/auth"
	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/guess"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/httpx"
	"github.com/SoWieMarkus/wheelgpt/internal/trackmania"
)

type HandlerConfig struct {
	// RateLimit is the number of plugin requests per channel and minute.
	// Zero disables the limit (tests).
	RateLimit int
}

type Handler struct {
	service *Service
	cfg     HandlerConfig
}

func NewHandler(service *Service, cfg HandlerConfig) *Handler {
	return &Handler{service: service, cfg: cfg}
}

type callerKey struct{}

func currentCaller(r *http.Request) Caller {
	c, _ := r.Context().Value(callerKey{}).(Caller)
	return c
}

// requirePlugin authenticates the plugin token. The old plugin sends the
// token without "Bearer ", that is accepted too.
func (h *Handler) requirePlugin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		caller, err := h.service.Authenticate(r.Context(), bearerToken(r.Header.Get("Authorization")))
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Bearer realm="wheelgpt"`)
			httpx.Error(w, r, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), callerKey{}, caller)))
	})
}

func (h *Handler) rateLimit(r chi.Router) {
	if h.cfg.RateLimit > 0 {
		r.Use(httprate.LimitBy(h.cfg.RateLimit, time.Minute,
			func(r *http.Request) (string, error) { return currentCaller(r).ChannelID, nil },
			httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
				httpx.Error(w, r, domain.ErrTooManyRequests)
			}),
		))
	}
}

func bearerToken(header string) string {
	header = strings.TrimSpace(header)
	if len(header) > 7 && strings.EqualFold(header[:7], "bearer ") {
		return strings.TrimSpace(header[7:])
	}
	return header
}

// Mount registers the plugin routes under /plugin.
func (h *Handler) Mount(r chi.Router) {
	r.Route("/plugin", func(r chi.Router) {
		r.Use(h.requirePlugin)
		h.rateLimit(r)
		r.Get("/status", h.status)
		r.Put("/map", h.putMap)
		r.Delete("/map", h.deleteMap)
		r.Put("/room", h.putRoom)
		r.Delete("/room", h.deleteRoom)
		r.Post("/personal-bests", h.personalBest)
	})
}

// MountProtected registers the token routes of the logged in streamer.
func (h *Handler) MountProtected(r chi.Router) {
	r.Get("/me/plugin-token", h.tokenInfo)
	r.Post("/me/plugin-token", h.createToken)
}

// --------------------------------------------------------------- plugin

// status checks the token and tells the plugin whether data is processed.
//
//	@Summary	Check the plugin token
//	@Tags		plugin
//	@Produce	json
//	@Security	PluginToken
//	@Success	200	{object}	StatusResponse
//	@Failure	401	{object}	httpx.Problem
//	@Router		/plugin/status [get]
func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	caller := currentCaller(r)
	st, err := h.service.Status(r.Context(), caller.ChannelID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, StatusResponse{
		ChannelID: st.Channel.ID, Login: st.Channel.Login, Accepting: st.Accepting,
		BotEnabled: st.Channel.Bot.Enabled, LegacyToken: caller.Legacy,
	})
}

// putMap sets the current map and clears the open guesses.
//
//	@Summary		Set the current map
//	@Description	Open guesses are cleared, they belonged to the previous map. The Trackmania Exchange id is looked up in the background. Plugin data is only accepted while the stream is live, otherwise 409 channel_offline.
//	@Tags			plugin
//	@Accept			json
//	@Produce		json
//	@Security		PluginToken
//	@Param			body	body		ingame.MapBody	true	"map"
//	@Success		200		{object}	ingame.MapResponse
//	@Failure		409		{object}	httpx.Problem	"the stream is not live"
//	@Failure		400		{object}	httpx.Problem
//	@Failure		401		{object}	httpx.Problem
//	@Router			/plugin/map [put]
func (h *Handler) putMap(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.Decode[ingame.MapBody](w, r)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	m, err := h.service.SetMap(r.Context(), currentCaller(r).ChannelID, body.ToMap())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, ingame.ToMapResponse(m))
}

// deleteMap clears the current map and the open guesses.
//
//	@Summary	Clear the current map
//	@Tags		plugin
//	@Security	PluginToken
//	@Success	204
//	@Failure	409	{object}	httpx.Problem	"the stream is not live"
//	@Failure	401	{object}	httpx.Problem
//	@Router		/plugin/map [delete]
func (h *Handler) deleteMap(w http.ResponseWriter, r *http.Request) {
	h.done(w, r, h.service.ClearMap(r.Context(), currentCaller(r).ChannelID))
}

// putRoom sets the server room.
//
//	@Summary	Set the current room
//	@Tags		plugin
//	@Accept		json
//	@Produce	json
//	@Security	PluginToken
//	@Param		body	body		ingame.RoomBody	true	"room"
//	@Success	200		{object}	ingame.RoomBody
//	@Failure	409		{object}	httpx.Problem	"the stream is not live"
//	@Failure	400		{object}	httpx.Problem
//	@Failure	401		{object}	httpx.Problem
//	@Router		/plugin/room [put]
func (h *Handler) putRoom(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.Decode[ingame.RoomBody](w, r)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	if err := h.service.SetRoom(r.Context(), currentCaller(r).ChannelID, ingame.Room(body)); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, body)
}

// deleteRoom clears the server room.
//
//	@Summary	Clear the current room
//	@Tags		plugin
//	@Security	PluginToken
//	@Success	204
//	@Failure	409	{object}	httpx.Problem	"the stream is not live"
//	@Failure	401	{object}	httpx.Problem
//	@Router		/plugin/room [delete]
func (h *Handler) deleteRoom(w http.ResponseWriter, r *http.Request) {
	h.done(w, r, h.service.ClearRoom(r.Context(), currentCaller(r).ChannelID))
}

// personalBest reports a new PB: the round is evaluated right away and the
// chat announcement follows after the channel's delay, hence 202.
//
//	@Summary	Report a personal best
//	@Tags		plugin
//	@Accept		json
//	@Produce	json
//	@Security	PluginToken
//	@Param		body	body		PersonalBestBody	true	"personal best"
//	@Success	202		{object}	PersonalBestResponse
//	@Failure	400		{object}	httpx.Problem
//	@Failure	409		{object}	httpx.Problem	"the stream is not live"
//	@Failure	401		{object}	httpx.Problem
//	@Router		/plugin/personal-bests [post]
func (h *Handler) personalBest(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.Decode[PersonalBestBody](w, r)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	res, err := h.service.PersonalBest(r.Context(), currentCaller(r).ChannelID, trackmania.Time(body.Time))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusAccepted, toPersonalBestResponse(res))
}

func (h *Handler) done(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.NoContent(w)
}

// --------------------------------------------------------------- token

// tokenInfo returns plugin token metadata. The token is only shown once.
//
//	@Summary	Plugin token metadata
//	@Tags		me
//	@Produce	json
//	@Success	200	{object}	TokenInfoResponse
//	@Failure	401	{object}	httpx.Problem
//	@Router		/me/plugin-token [get]
func (h *Handler) tokenInfo(w http.ResponseWriter, r *http.Request) {
	info, err := h.service.TokenInfo(r.Context(), auth.CurrentChannel(r).ChannelID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	out := TokenInfoResponse{Exists: info.Token != nil, LegacyTokenActive: info.LegacyTokenActive}
	if info.Token != nil {
		out.CreatedAt, out.LastUsedAt = &info.Token.CreatedAt, info.Token.LastUsedAt
	}
	httpx.JSON(w, http.StatusOK, out)
}

// createToken issues a new plugin token and revokes the previous one.
//
//	@Summary	Create a plugin token
//	@Tags		me
//	@Produce	json
//	@Success	201	{object}	TokenCreatedResponse
//	@Failure	401	{object}	httpx.Problem
//	@Router		/me/plugin-token [post]
func (h *Handler) createToken(w http.ResponseWriter, r *http.Request) {
	token, at, err := h.service.CreateToken(r.Context(), auth.CurrentChannel(r).ChannelID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	w.Header().Set("Location", "/api/v1/me/plugin-token")
	httpx.JSON(w, http.StatusCreated, TokenCreatedResponse{Token: token, CreatedAt: at})
}

// --------------------------------------------------------------- legacy

// legacySunset is announced in the Sunset header of the old plugin paths.
const legacySunset = "Wed, 30 Jun 2027 00:00:00 GMT"

// MountLegacy keeps the plugin paths of the old backend working. Released
// plugin builds call them, treat only 200 as success and send the raw
// token. Bodies are parsed leniently, like the old backend did.
func (h *Handler) MountLegacy(r chi.Router) {
	r.Route("/api/trackmania", func(r chi.Router) {
		r.Use(deprecated, h.requirePlugin)
		h.rateLimit(r)
		r.Post("/map", legacy(func(r *http.Request, id string, body ingame.MapBody) error {
			_, err := h.service.SetMap(r.Context(), id, body.ToMap())
			return err
		}, "Successfully updated the map."))
		r.Delete("/map", legacy(func(r *http.Request, id string, _ struct{}) error {
			return h.service.ClearMap(r.Context(), id)
		}, "Map deleted successfully."))
		r.Post("/room", legacy(func(r *http.Request, id string, body ingame.RoomBody) error {
			return h.service.SetRoom(r.Context(), id, ingame.Room(body))
		}, "Room updated successfully."))
		r.Delete("/room", legacy(func(r *http.Request, id string, _ struct{}) error {
			return h.service.ClearRoom(r.Context(), id)
		}, "Room deleted successfully."))
		r.Post("/pb", legacy(func(r *http.Request, id string, body PersonalBestBody) error {
			_, err := h.service.PersonalBest(r.Context(), id, trackmania.Time(body.Time))
			return err
		}, "Personal best updated successfully."))
	})
}

func deprecated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Deprecation", "@1767225600") // 2026-01-01
		h.Set("Sunset", legacySunset)
		h.Set("Link", `</api/v1/plugin>; rel="successor-version"`)
		next.ServeHTTP(w, r)
	})
}

type legacyMessage struct {
	Message string `json:"message"`
}

// legacy adapts a plugin action to the old API. Offline channels get 200 with
// "Channel is offline." like the old backend sent: old plugin builds retry
// every other status.
func legacy[T any](fn func(r *http.Request, channelID string, body T) error, success string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body T
		if r.Method == http.MethodPost {
			if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(&body); err != nil {
				httpx.Error(w, r, httpx.BadRequest("invalid_body", "malformed JSON"))
				return
			}
			if err := httpx.Validate(body); err != nil {
				httpx.Error(w, r, err)
				return
			}
		}
		err := fn(r, currentCaller(r).ChannelID, body)
		switch {
		case errors.Is(err, domain.ErrChannelOffline):
			httpx.JSON(w, http.StatusOK, legacyMessage{Message: "Channel is offline."})
		case err != nil:
			httpx.Error(w, r, err)
		default:
			httpx.JSON(w, http.StatusOK, legacyMessage{Message: success})
		}
	}
}

// --------------------------------------------------------------- dto

type StatusResponse struct {
	ChannelID string `json:"channelId" example:"12345678"`
	Login     string `json:"login"     example:"granadyy"`
	// Whether map, room and PB data is accepted right now (stream is live).
	Accepting  bool `json:"accepting"`
	BotEnabled bool `json:"botEnabled"`
	// The request used a token of the old backend.
	LegacyToken bool `json:"legacyToken"`
}

type PersonalBestBody struct {
	// Race time in milliseconds.
	Time int64 `json:"time" validate:"gte=0" example:"47881"`
}

type WinnerResponse struct {
	UserID      string `json:"userId"      example:"98765"`
	DisplayName string `json:"displayName" example:"tinytires"`
	GuessTime   int64  `json:"guessTime"   example:"47850"`
}

type RoundResponse struct {
	ID           int64            `json:"id"           example:"42"`
	ResultTime   int64            `json:"resultTime"   example:"47881"`
	MapName      *string          `json:"mapName,omitempty" binding:"optional" example:"Spring 2026 - 07"`
	Participants int              `json:"participants" example:"4"`
	Points       int64            `json:"points"       example:"1"`
	Perfect      bool             `json:"perfect"`
	Winners      []WinnerResponse `json:"winners"`
}

type PersonalBestResponse struct {
	// Why the result is not announced, if it is not.
	Reason     string        `json:"reason,omitempty"     binding:"optional"`
	AnnounceAt *time.Time    `json:"announceAt,omitempty" binding:"optional"`
	Round      RoundResponse `json:"round"`
}

func toPersonalBestResponse(res PersonalBestResult) PersonalBestResponse {
	return PersonalBestResponse{Reason: res.Reason, AnnounceAt: res.AnnounceAt, Round: toRound(res.Round)}
}

func toRound(rd guess.Round) RoundResponse {
	out := RoundResponse{
		ID: rd.ID, ResultTime: rd.Result.Milliseconds(), MapName: rd.MapName, Participants: rd.Participants,
		Points: rd.Points, Perfect: rd.Perfect, Winners: make([]WinnerResponse, len(rd.Winners)),
	}
	for i, w := range rd.Winners {
		out.Winners[i] = WinnerResponse{UserID: w.UserID, DisplayName: w.DisplayName, GuessTime: w.Time.Milliseconds()}
	}
	return out
}

type TokenInfoResponse struct {
	Exists     bool       `json:"exists"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"  binding:"optional"`
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty" binding:"optional"`
	// A token of the old backend still works. Creating a new token revokes it.
	LegacyTokenActive bool `json:"legacyTokenActive"`
}

type TokenCreatedResponse struct {
	// Shown only once.
	Token     string    `json:"token" example:"wgpt_6Q1q..."`
	CreatedAt time.Time `json:"createdAt"`
}
