package leaderboard

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/httpx"
)

// Channels resolves the channel of a path.
type Channels interface {
	Resolve(ctx context.Context, ref string) (channel.Channel, error)
}

type Handler struct {
	service  *Service
	channels Channels
}

func NewHandler(service *Service, channels Channels) *Handler {
	return &Handler{service: service, channels: channels}
}

func (h *Handler) Mount(r chi.Router) {
	// Single routes, a chi sub-router on /channels/{channel} would shadow
	// the channel route of package channel.
	r.Get("/channels/{channel}/leaderboards", h.periods)
	r.Get("/channels/{channel}/leaderboards/{period}", h.board)
}

// periods lists the leaderboard periods of a channel.
//
//	@Summary	Leaderboard periods
//	@Tags		leaderboards
//	@Produce	json
//	@Param		channel	path		string	true	"Twitch user id or login"
//	@Success	200		{object}	PeriodsResponse
//	@Failure	404		{object}	httpx.Problem
//	@Router		/channels/{channel}/leaderboards [get]
func (h *Handler) periods(w http.ResponseWriter, r *http.Request) {
	c, ok := h.channel(w, r)
	if !ok {
		return
	}
	current, periods, err := h.service.Periods(r.Context(), c.ID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, PeriodsResponse{Current: current, Periods: periods})
}

// board returns one page of a leaderboard.
//
//	@Summary		Leaderboard
//	@Description	all-time includes points from before monthly leaderboards existed. A month (YYYY-MM, UTC) sums the rounds of that month. Ties go to more perfect guesses, then to who got there first.
//	@Tags			leaderboards
//	@Produce		json
//	@Param			channel	path		string	true	"Twitch user id or login"
//	@Param			period	path		string	true	"all-time or YYYY-MM"
//	@Param			q		query		string	false	"filter by display name, positions stay those of the full board"
//	@Param			limit	query		int		false	"page size (1-100)"	default(50)
//	@Param			offset	query		int		false	"items to skip"		default(0)
//	@Success		200		{object}	httpx.PageDoc{items=[]EntryResponse}
//	@Failure		400		{object}	httpx.Problem
//	@Failure		404		{object}	httpx.Problem
//	@Router			/channels/{channel}/leaderboards/{period} [get]
func (h *Handler) board(w http.ResponseWriter, r *http.Request) {
	c, ok := h.channel(w, r)
	if !ok {
		return
	}
	key, err := httpx.PathSlug(r, "period")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	period, err := ParsePeriod(key)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	limit, offset, err := httpx.Pagination(r, 50)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	q := r.URL.Query().Get("q")
	if len(q) > 50 {
		httpx.Error(w, r, httpx.BadRequest("invalid_query", "q must be at most 50 characters long"))
		return
	}
	entries, total, err := h.service.Board(r.Context(), c.ID, period, q, limit, offset)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	items := make([]EntryResponse, len(entries))
	for i, e := range entries {
		items[i] = EntryResponse(e)
	}
	httpx.JSON(w, http.StatusOK, httpx.Page[EntryResponse]{Items: items, Total: total, Limit: limit, Offset: offset})
}

func (h *Handler) channel(w http.ResponseWriter, r *http.Request) (channel.Channel, bool) {
	ref, err := httpx.PathSlug(r, "channel")
	if err != nil {
		httpx.Error(w, r, err)
		return channel.Channel{}, false
	}
	c, err := h.channels.Resolve(r.Context(), ref)
	if err != nil {
		httpx.Error(w, r, err)
		return channel.Channel{}, false
	}
	return c, true
}

// --------------------------------------------------------------- dto

type EntryResponse struct {
	Position       int64  `json:"position"       example:"1"`
	UserID         string `json:"userId"         example:"98765"`
	DisplayName    string `json:"displayName"    example:"tinytires"`
	Points         int64  `json:"points"         example:"26"`
	PerfectGuesses int64  `json:"perfectGuesses" example:"2"`
}

type PeriodsResponse struct {
	// The month that is running now.
	Current string `json:"current" example:"2026-11"`
	// This month, every month with points (newest first) and all-time.
	Periods []string `json:"periods" example:"2026-11,2026-10,all-time"`
}
