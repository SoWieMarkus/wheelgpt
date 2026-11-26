package channel

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/SoWieMarkus/gokit/httpx"

	"github.com/SoWieMarkus/wheelgpt/internal/auth"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/paging"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/slug"
)

// Stream switches a channel live or offline the way the Twitch events do:
// going offline clears map, room and open guesses.
type Stream interface {
	StreamOnline(ctx context.Context, channelID string) error
	StreamOffline(ctx context.Context, channelID string) error
}

type HandlerConfig struct {
	// DevMode lets the dashboard pretend the stream is live, so the plugin
	// can be tested without streaming. Local development only.
	DevMode bool
	// ChatEvents is false when the server doesn't subscribe to chats
	// (EVENTSUB_SYNC off). The bot status is then reported as "off".
	ChatEvents bool
}

type Handler struct {
	service *Service
	stream  Stream
	cfg     HandlerConfig
}

func NewHandler(service *Service, stream Stream, cfg HandlerConfig) *Handler {
	return &Handler{service: service, stream: stream, cfg: cfg}
}

// Mount registers the public routes.
func (h *Handler) Mount(r chi.Router) {
	r.Get("/channels", h.list)
	r.Get("/channels/{channel}", h.get)
	r.Get("/stats", h.stats)
}

// MountProtected registers the routes of the logged in streamer.
func (h *Handler) MountProtected(r chi.Router) {
	r.Route("/me", func(r chi.Router) {
		r.Get("/", h.me)
		r.Delete("/", h.deleteMe)
		r.Get("/settings", h.settings)
		r.Patch("/settings", h.updateSettings)
		r.Put("/bot-connection", h.connect)
		r.Delete("/bot-connection", h.disconnect)
		if h.cfg.DevMode {
			r.Put("/live", h.goLive)
			r.Delete("/live", h.endStream)
		}
	})
}

// --------------------------------------------------------------- public

// list returns the channels that want to be listed, live ones first.
//
//	@Summary	List channels using WheelGPT
//	@Tags		channels
//	@Produce	json
//	@Param		live	query		bool	false	"only live channels"
//	@Param		q		query		string	false	"filter by login or display name"
//	@Param		limit	query		int		false	"page size (1-100)"	default(25)
//	@Param		offset	query		int		false	"items to skip"		default(0)
//	@Success	200		{object}	paging.PageDoc{items=[]SummaryResponse}
//	@Failure	400		{object}	httpx.Problem
//	@Router		/channels [get]
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := paging.Parse(r, 25)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	live, err := httpx.QueryBool(r, "live")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	q := r.URL.Query().Get("q")
	if len(q) > 50 {
		httpx.Error(w, r, httpx.BadRequest("invalid_query", "q must be at most 50 characters long"))
		return
	}
	list, total, err := h.service.ListPublic(r.Context(), live, q, limit, offset)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	items := make([]SummaryResponse, len(list))
	for i, c := range list {
		items[i] = SummaryResponse(c)
	}
	httpx.JSON(w, http.StatusOK, paging.Page[SummaryResponse]{Items: items, Total: total, Limit: limit, Offset: offset})
}

// get returns one channel, also channels that are not listed.
//
//	@Summary	Get a channel
//	@Tags		channels
//	@Produce	json
//	@Param		channel	path		string	true	"Twitch user id or login"
//	@Success	200		{object}	SummaryResponse
//	@Failure	404		{object}	httpx.Problem
//	@Router		/channels/{channel} [get]
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	ref, err := slug.FromPath(r, "channel")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	c, err := h.service.Resolve(r.Context(), ref)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toSummary(c))
}

// stats returns usage totals for the landing page.
//
//	@Summary	Usage totals
//	@Tags		channels
//	@Produce	json
//	@Success	200	{object}	StatsResponse
//	@Router		/stats [get]
func (h *Handler) stats(w http.ResponseWriter, r *http.Request) {
	s, err := h.service.Stats(r.Context())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, StatsResponse(s))
}

// --------------------------------------------------------------- me

// me returns the logged in streamer's channel.
//
//	@Summary	My channel
//	@Tags		me
//	@Produce	json
//	@Success	200	{object}	MeResponse
//	@Failure	401	{object}	httpx.Problem
//	@Router		/me [get]
func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	c, err := h.service.Get(r.Context(), auth.CurrentChannel(r).ChannelID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, MeResponse{
		Channel:  toSummary(c),
		Settings: toSettings(c.Settings),
		Bot:      h.toBot(c.Bot),
		DevMode:  h.cfg.DevMode,
	})
}

// goLive pretends the stream started, so the plugin's data is accepted.
//
//	@Summary		Pretend to be live (dev mode)
//	@Description	Only with DEV_MODE on a local server. Marks the channel live like a stream.online event from Twitch.
//	@Tags			me
//	@Produce		json
//	@Success		200	{object}	SummaryResponse
//	@Failure		401	{object}	httpx.Problem
//	@Router			/me/live [put]
func (h *Handler) goLive(w http.ResponseWriter, r *http.Request) {
	h.setLive(w, r, true)
}

// endStream pretends the stream ended: map, room and open guesses are
// cleared like after a real stream.
//
//	@Summary		Pretend the stream ended (dev mode)
//	@Description	Only with DEV_MODE on a local server. Works like a stream.offline event from Twitch: map, room and open guesses are cleared.
//	@Tags			me
//	@Produce		json
//	@Success		200	{object}	SummaryResponse
//	@Failure		401	{object}	httpx.Problem
//	@Router			/me/live [delete]
func (h *Handler) endStream(w http.ResponseWriter, r *http.Request) {
	h.setLive(w, r, false)
}

func (h *Handler) setLive(w http.ResponseWriter, r *http.Request, live bool) {
	id := auth.CurrentChannel(r).ChannelID
	var err error
	if live {
		err = h.stream.StreamOnline(r.Context(), id)
	} else {
		err = h.stream.StreamOffline(r.Context(), id)
	}
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	c, err := h.service.Get(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toSummary(c))
}

// deleteMe removes the bot from the channel and deletes all its data.
//
//	@Summary	Delete my channel and all its data
//	@Tags		me
//	@Success	204
//	@Failure	401	{object}	httpx.Problem
//	@Router		/me [delete]
func (h *Handler) deleteMe(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), auth.CurrentChannel(r).ChannelID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	auth.ClearSessionCookie(w, r)
	httpx.NoContent(w)
}

// settings returns the channel settings.
//
//	@Summary	My settings
//	@Tags		me
//	@Produce	json
//	@Success	200	{object}	SettingsResponse
//	@Failure	401	{object}	httpx.Problem
//	@Router		/me/settings [get]
func (h *Handler) settings(w http.ResponseWriter, r *http.Request) {
	c, err := h.service.Get(r.Context(), auth.CurrentChannel(r).ChannelID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toSettings(c.Settings))
}

// updateSettings changes the given settings, omitted fields stay as they are.
//
//	@Summary	Update my settings
//	@Tags		me
//	@Accept		json
//	@Produce	json
//	@Param		body	body		SettingsBody	true	"settings to change"
//	@Success	200		{object}	SettingsResponse
//	@Failure	400		{object}	httpx.Problem
//	@Failure	401		{object}	httpx.Problem
//	@Router		/me/settings [patch]
func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.Decode[SettingsBody](w, r)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	s, err := h.service.UpdateSettings(r.Context(), auth.CurrentChannel(r).ChannelID, SettingsPatch(body))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, toSettings(s))
}

// connect enables the bot and subscribes to the chat via EventSub.
//
//	@Summary		Connect the bot to my chat
//	@Description	If Twitch refuses, chatSubscription is needs_authorization: log in again to grant channel:bot or make the bot a moderator.
//	@Tags			me
//	@Produce		json
//	@Success		200	{object}	BotResponse
//	@Failure		401	{object}	httpx.Problem
//	@Failure		502	{object}	httpx.Problem
//	@Router			/me/bot-connection [put]
func (h *Handler) connect(w http.ResponseWriter, r *http.Request) {
	b, err := h.service.ConnectBot(r.Context(), auth.CurrentChannel(r).ChannelID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, h.toBot(b))
}

// disconnect disables the bot in the chat.
//
//	@Summary	Disconnect the bot from my chat
//	@Tags		me
//	@Produce	json
//	@Success	200	{object}	BotResponse
//	@Failure	401	{object}	httpx.Problem
//	@Router		/me/bot-connection [delete]
func (h *Handler) disconnect(w http.ResponseWriter, r *http.Request) {
	b, err := h.service.DisconnectBot(r.Context(), auth.CurrentChannel(r).ChannelID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, h.toBot(b))
}

// --------------------------------------------------------------- dto

type SummaryResponse struct {
	ID           string `json:"id"           example:"12345678"`
	Login        string `json:"login"        example:"granadyy"`
	DisplayName  string `json:"displayName"  example:"GranaDyy"`
	ProfileImage string `json:"profileImage" example:"https://static-cdn.jtvnw.net/user.png"`
	IsLive       bool   `json:"isLive"`
}

func toSummary(c Channel) SummaryResponse {
	return SummaryResponse{ID: c.ID, Login: c.Login, DisplayName: c.DisplayName, ProfileImage: c.ProfileImage, IsLive: c.IsLive}
}

type StatsResponse struct {
	Channels int64 `json:"channels" example:"80"`
}

type SettingsResponse struct {
	// Seconds between the PB and the chat announcement (stream delay).
	GuessDelaySeconds float64 `json:"guessDelaySeconds" example:"2"`
	// Guesses younger than this are ignored when a PB comes in.
	GuessMinAgeSeconds float64 `json:"guessMinAgeSeconds" example:"0"`
	// List the channel on the website.
	Public bool `json:"public" example:"true"`
}

func toSettings(s Settings) SettingsResponse {
	return SettingsResponse(s)
}

// SettingsBody changes settings. Omitted fields stay unchanged.
type SettingsBody struct {
	GuessDelaySeconds  *float64 `json:"guessDelaySeconds,omitempty"    validate:"omitempty,gte=0,lte=60" example:"2.5"`
	GuessMinAgeSeconds *float64 `json:"guessMinAgeSeconds,omitempty"   validate:"omitempty,gte=0,lte=60" example:"5"`
	Public             *bool    `json:"public,omitempty"`
}

type BotResponse struct {
	Enabled bool `json:"enabled"`
	// unknown, enabled, needs_authorization, disabled, or off when the
	// server doesn't subscribe to chats (local development)
	ChatSubscription string `json:"chatSubscription" example:"enabled" enums:"unknown,enabled,needs_authorization,disabled,off"`
	// The streamer granted channel:bot through the new login.
	Authorized   bool       `json:"authorized"`
	AuthorizedAt *time.Time `json:"authorizedAt,omitempty" binding:"optional"`
	Scopes       []string   `json:"scopes" example:"channel:bot"`
}

func (h *Handler) toBot(b Bot) BotResponse {
	chat := b.ChatSubscription
	if b.Enabled && !h.cfg.ChatEvents {
		chat = ChatOff
	}
	return BotResponse{Enabled: b.Enabled, ChatSubscription: chat, Authorized: b.Authorized(), AuthorizedAt: b.AuthorizedAt, Scopes: b.Scopes}
}

type MeResponse struct {
	Channel  SummaryResponse  `json:"channel"`
	Settings SettingsResponse `json:"settings"`
	Bot      BotResponse      `json:"bot"`
	// The server runs in dev mode: the live status can be faked.
	DevMode bool `json:"devMode"`
}
