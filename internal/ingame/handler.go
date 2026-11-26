package ingame

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/SoWieMarkus/gokit/httpx"

	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/slug"
	"github.com/SoWieMarkus/wheelgpt/internal/trackmania"
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
	r.Get("/channels/{channel}/map", h.getMap)
	r.Get("/channels/{channel}/room", h.getRoom)
}

// getMap returns the map the streamer is playing.
//
//	@Summary	Current map
//	@Tags		channels
//	@Produce	json
//	@Param		channel	path		string	true	"Twitch user id or login"
//	@Success	200		{object}	MapResponse
//	@Failure	404		{object}	httpx.Problem
//	@Router		/channels/{channel}/map [get]
func (h *Handler) getMap(w http.ResponseWriter, r *http.Request) {
	c, ok := h.channel(w, r)
	if !ok {
		return
	}
	m, err := h.service.Map(r.Context(), c.ID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, ToMapResponse(m))
}

// getRoom returns the server room the streamer is in.
//
//	@Summary	Current room
//	@Tags		channels
//	@Produce	json
//	@Param		channel	path		string	true	"Twitch user id or login"
//	@Success	200		{object}	RoomBody
//	@Failure	404		{object}	httpx.Problem
//	@Router		/channels/{channel}/room [get]
func (h *Handler) getRoom(w http.ResponseWriter, r *http.Request) {
	c, ok := h.channel(w, r)
	if !ok {
		return
	}
	room, err := h.service.Room(r.Context(), c.ID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, http.StatusOK, RoomBody(room))
}

func (h *Handler) channel(w http.ResponseWriter, r *http.Request) (channel.Channel, bool) {
	ref, err := slug.FromPath(r, "channel")
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

// Medals are medal times in milliseconds.
type Medals struct {
	// Only if the Champion Medals plugin knows it.
	Champion *int64 `json:"champion,omitempty" binding:"optional" example:"46850"`
	Author   int64  `json:"author"   example:"47904"`
	Gold     int64  `json:"gold"     example:"51000"`
	Silver   int64  `json:"silver"   example:"58000"`
	Bronze   int64  `json:"bronze"   example:"72000"`
}

type MapResponse struct {
	UID string `json:"uid" example:"Rq3h3pT5ZyvA0lqk1eYyQ2xYb1c"`
	// Name without Trackmania formatting codes.
	Name            string  `json:"name"   example:"Spring 2026 - 07"`
	Author          string  `json:"author" example:"Nadeo"`
	Medals          Medals  `json:"medals"`
	TMXID           *int64  `json:"tmxId,omitempty"  binding:"optional" example:"123456"`
	TMXURL          *string `json:"tmxUrl,omitempty" binding:"optional"`
	TrackmaniaIOURL string  `json:"trackmaniaIoUrl"`
}

// ToMapResponse converts a map for the API.
func ToMapResponse(m trackmania.Map) MapResponse {
	out := MapResponse{
		UID: m.UID, Name: trackmania.StripStyling(m.Name), Author: m.Author,
		Medals: Medals{
			Author: int64(m.AuthorTime), Gold: int64(m.GoldTime), Silver: int64(m.SilverTime), Bronze: int64(m.BronzeTime),
		},
		TrackmaniaIOURL: m.TrackmaniaIOURL(),
	}
	if m.HasChampionTime() {
		v := int64(m.ChampionTime)
		out.Medals.Champion = &v
	}
	if m.TMXID != 0 {
		id, link := m.TMXID, m.TMXURL()
		out.TMXID, out.TMXURL = &id, &link
	}
	return out
}

// MapBody is the map the plugin sends. Same body as the old
// POST /api/trackmania/map.
type MapBody struct {
	UID          string `json:"uid"          validate:"required,max=1000" example:"Rq3h3pT5ZyvA0lqk1eYyQ2xYb1c"`
	Name         string `json:"name"         validate:"required,max=1000" example:"$o$fffSpring 2026 - 07"`
	Author       string `json:"author"       validate:"required,max=1000" example:"Nadeo"`
	AuthorTime   int64  `json:"authorTime"   validate:"gte=0" example:"47904"`
	GoldTime     int64  `json:"goldTime"     validate:"gte=0" example:"51000"`
	SilverTime   int64  `json:"silverTime"   validate:"gte=0" example:"58000"`
	BronzeTime   int64  `json:"bronzeTime"   validate:"gte=0" example:"72000"`
	ChampionTime int64  `json:"championTime" validate:"gte=0" example:"0"`
}

// ToMap converts the plugin body.
func (b MapBody) ToMap() trackmania.Map {
	return trackmania.Map{
		UID: b.UID, Name: b.Name, Author: b.Author,
		AuthorTime: trackmania.Time(b.AuthorTime), GoldTime: trackmania.Time(b.GoldTime),
		SilverTime: trackmania.Time(b.SilverTime), BronzeTime: trackmania.Time(b.BronzeTime),
		ChampionTime: trackmania.Time(b.ChampionTime),
	}
}

// RoomBody is a server room.
type RoomBody struct {
	Login           string `json:"login"           validate:"required,max=1000" example:"cotd-practice-1"`
	Name            string `json:"name"            validate:"required,max=1000" example:"Cup of the Day Practice"`
	NumberOfPlayers int64  `json:"numberOfPlayers" validate:"gte=0,lte=2000" example:"42"`
	MaxPlayers      int64  `json:"maxPlayers"      validate:"gte=0,lte=2000" example:"64"`
}
