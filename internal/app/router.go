package app

// This file assembles the HTTP API: global middleware, health checks, the
// OpenAPI docs and the routes of every feature package.

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"

	// Registers the generated OpenAPI spec with swag.
	_ "github.com/SoWieMarkus/wheelgpt/api/openapi"
	"github.com/SoWieMarkus/wheelgpt/internal/auth"
	"github.com/SoWieMarkus/wheelgpt/internal/bot"
	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/eventsub"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/leaderboard"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/httpx"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/metrics"
	"github.com/SoWieMarkus/wheelgpt/internal/plugin"
)

type pinger interface {
	Ping(ctx context.Context) error
}

// handlers are the HTTP handlers of all feature packages.
type handlers struct {
	auth        *auth.Handler
	channel     *channel.Handler
	ingame      *ingame.Handler
	leaderboard *leaderboard.Handler
	plugin      *plugin.Handler
	bot         *bot.Handler
	eventsub    *eventsub.Handler
}

type routerConfig struct {
	version           string
	publicOrigin      string
	trustProxyHeaders bool
}

// newRouter builds the router for the whole API.
func newRouter(cfg routerConfig, db pinger, m *metrics.Metrics, h handlers) (http.Handler, error) {
	csrf, err := crossOriginProtection(cfg.publicOrigin)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	if cfg.trustProxyHeaders {
		r.Use(realIPFromProxy)
	}
	r.Use(
		requestLogger,
		middleware.Recoverer,
		m.Middleware,
		securityHeaders,
		middleware.Timeout(30*time.Second),
	)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteProblem(w, r, http.StatusNotFound, "route_not_found", "this route does not exist", nil)
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteProblem(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed on this route", nil)
	})

	r.Get("/healthz", health(cfg.version))
	r.Get("/readyz", ready(cfg.version, db))
	r.Group(func(r chi.Router) {
		// The swagger UI needs scripts and styles, relax the API policy.
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Del("Content-Security-Policy")
				next.ServeHTTP(w, req)
			})
		})
		r.Get("/api/docs/*", httpSwagger.Handler(httpSwagger.URL("/api/docs/doc.json")))
	})

	// Plugin paths and the EventSub callback of the old backend.
	h.plugin.MountLegacy(r)
	r.Post("/api/twitch/webhook/stream", h.eventsub.ServeHTTP)

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(noStore, csrf)
		h.auth.Mount(r)
		h.channel.Mount(r)
		h.ingame.Mount(r)
		h.leaderboard.Mount(r)
		h.bot.Mount(r)
		h.plugin.Mount(r)
		h.eventsub.Mount(r)

		r.Group(func(r chi.Router) {
			r.Use(h.auth.RequireChannel)
			h.channel.MountProtected(r)
			h.plugin.MountProtected(r)
		})
	})
	return r, nil
}
