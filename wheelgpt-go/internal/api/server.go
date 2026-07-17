package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/SoWieMarkus/wheelgpt/internal/auth"
	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/config"
	"github.com/SoWieMarkus/wheelgpt/internal/db"
	"github.com/SoWieMarkus/wheelgpt/internal/twitch"
)

type Server struct {
	cfg    config.Config
	store  *db.Store
	tc     *twitch.Client
	bot    *bot.Bot
	router *chi.Mux
}

func NewServer(cfg config.Config, store *db.Store, tc *twitch.Client, b *bot.Bot) *Server {
	s := &Server{cfg: cfg, store: store, tc: tc, bot: b}
	s.router = chi.NewRouter()
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.mount()
	return s
}

func (s *Server) Handler() http.Handler { return s.router }

func (s *Server) mount() {
	s.router.Get("/metrics", promhttp.Handler().ServeHTTP)

	s.router.Route("/api", func(r chi.Router) {
		// Auth
		r.Get("/authentication/twitch", s.handleTwitchRedirect)
		r.Post("/authentication/login", s.handleLogin)

		// Twitch EventSub
		r.Post("/twitch/webhook", s.handleEventSub)

		// Plugin routes (plugin JWT required)
		r.Group(func(r chi.Router) {
			r.Use(s.pluginAuthMiddleware)
			r.Post("/trackmania/map", s.handleUpsertMap)
			r.Post("/trackmania/room", s.handleUpsertRoom)
			r.Post("/trackmania/pb", s.handleNewPB)
		})

		// Web routes (web JWT required)
		r.Group(func(r chi.Router) {
			r.Use(s.webAuthMiddleware)
			r.Get("/channel", s.handleGetChannel)
			r.Put("/channel", s.handleUpdateChannel)
			r.Get("/leaderboard", s.handleGetLeaderboard)
			r.Put("/authentication/renew", s.handleRenewToken)
		})
	})
}

// --- middleware ---

func (s *Server) webAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		channelID, err := auth.VerifyWebToken(token, s.cfg.JWTSecretWeb)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(withChannelID(r.Context(), channelID)))
	})
}

func (s *Server) pluginAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		channelID, channelToken, err := auth.VerifyPluginToken(token, s.cfg.JWTSecretChannel)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ch, err := s.store.GetChannel(r.Context(), channelID)
		if err != nil || ch.Token != channelToken {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !ch.BotActiveWhenOffline.Bool && !ch.IsLive {
			writeJSON(w, http.StatusOK, map[string]string{"message": "Channel is offline."})
			return
		}
		next.ServeHTTP(w, r.WithContext(withChannelID(r.Context(), channelID)))
	})
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
