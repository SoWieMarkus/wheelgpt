package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/httpx"
)

type healthResponse struct {
	Status  string `json:"status"  example:"ok"`
	Version string `json:"version" example:"1.0.0"`
}

// health reports liveness (GET /healthz, outside the versioned API).
func health(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, healthResponse{Status: "ok", Version: version})
	}
}

// ready reports whether the service can reach its database (GET /readyz).
func ready(version string, db pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			httpx.WriteProblem(w, r, http.StatusServiceUnavailable, "database_unavailable", "database is not reachable", nil)
			return
		}
		httpx.JSON(w, http.StatusOK, healthResponse{Status: "ok", Version: version})
	}
}

// requestLogger writes one structured log line per request.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		level := slog.LevelInfo
		switch {
		case ww.Status() >= http.StatusInternalServerError:
			level = slog.LevelError
		case r.URL.Path == "/healthz" || r.URL.Path == "/readyz":
			level = slog.LevelDebug
		}
		slog.Log(r.Context(), level, "http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"bytes", ww.BytesWritten(),
			"duration", time.Since(start),
			"request_id", middleware.GetReqID(r.Context()),
		)
	})
}

// securityHeaders sets conservative defaults for an API that serves JSON only.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
		next.ServeHTTP(w, r)
	})
}

// noStore marks API responses as not cacheable: most carry session data.
func noStore(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// realIPFromProxy replaces RemoteAddr with the X-Real-IP header. Only use it
// behind a reverse proxy that always overwrites that header (the website's
// nginx does), otherwise clients could pick their own IP.
func realIPFromProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ip, err := netip.ParseAddr(strings.TrimSpace(r.Header.Get("X-Real-IP"))); err == nil {
			r.RemoteAddr = net.JoinHostPort(ip.String(), "0")
		}
		next.ServeHTTP(w, r)
	})
}

// crossOriginProtection rejects state changing requests a browser sends from
// another site (CSRF), based on Sec-Fetch-Site and Origin. The session
// cookie is SameSite=Lax as well, this is the second layer. Requests without
// those headers (the plugin, Twitch webhooks) pass.
func crossOriginProtection(trustedOrigin string) (func(http.Handler) http.Handler, error) {
	cop := http.NewCrossOriginProtection()
	if err := cop.AddTrustedOrigin(trustedOrigin); err != nil {
		return nil, err
	}
	cop.SetDenyHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, r, domain.ErrCrossOrigin)
	}))
	return cop.Handler, nil
}
