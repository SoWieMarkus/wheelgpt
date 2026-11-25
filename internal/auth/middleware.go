package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/SoWieMarkus/wheelgpt/internal/domain"
	"github.com/SoWieMarkus/wheelgpt/internal/platform/httpx"
)

type contextKey struct{}

var identityKey contextKey

// CurrentChannel returns the logged in streamer. Only call it behind
// RequireChannel.
func CurrentChannel(r *http.Request) Identity {
	id, _ := r.Context().Value(identityKey).(Identity)
	return id
}

// WithIdentity stores id in ctx. Exposed for handler tests of other packages.
func WithIdentity(ctx context.Context, id Identity) context.Context {
	return context.WithValue(ctx, identityKey, id)
}

// RequireChannel resolves the session cookie and rejects anonymous requests.
func (h *Handler) RequireChannel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(h.cookieName())
		if err != nil {
			httpx.Error(w, r, domain.ErrUnauthenticated)
			return
		}
		identity, err := h.service.Authenticate(r.Context(), cookie.Value)
		if errors.Is(err, domain.ErrUnauthenticated) {
			h.clearSessionCookie(w)
		}
		if err != nil {
			httpx.Error(w, r, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithIdentity(r.Context(), identity)))
	})
}

// ClearSessionCookie removes the session cookie the request carries, for
// handlers of other packages that end the session (deleting the account).
func ClearSessionCookie(w http.ResponseWriter, r *http.Request) {
	for _, name := range []string{secureCookieName, devCookieName} {
		if _, err := r.Cookie(name); err == nil {
			http.SetCookie(w, &http.Cookie{ //nolint:gosec // Secure follows the cookie name prefix
				Name: name, Value: "", Path: "/", MaxAge: -1, HttpOnly: true,
				Secure: name == secureCookieName, SameSite: http.SameSiteLaxMode,
			})
		}
	}
}
