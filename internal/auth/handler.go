package auth

import (
	"crypto/subtle"
	"encoding/base64"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"

	"github.com/SoWieMarkus/gokit/apperr"
	"github.com/SoWieMarkus/gokit/httpx"
)

const (
	// The __Host- prefix makes browsers reject the cookie unless it is
	// Secure, has Path=/ and no Domain, so it cannot be set by subdomains.
	secureCookieName = "__Host-wheelgpt_session"
	devCookieName    = "wheelgpt_session"

	stateCookie     = "wheelgpt_oauth"
	stateCookiePath = "/api/v1/auth"
	purposeLogin    = "login"
	purposeBot      = "bot"
)

type HandlerConfig struct {
	// PublicOrigin is scheme://host of the website. OAuth redirect URIs are
	// built from it.
	PublicOrigin  string
	SecureCookies bool
	// LoginRateLimit is the number of login attempts per IP and minute.
	// Zero disables the limit (tests).
	LoginRateLimit int
	// DevLogin enables /auth/dev, which logs in without Twitch.
	DevLogin bool
}

type Handler struct {
	service *Service
	cfg     HandlerConfig
}

func NewHandler(service *Service, cfg HandlerConfig) *Handler {
	return &Handler{service: service, cfg: cfg}
}

// Mount registers the public routes under /auth.
func (h *Handler) Mount(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/logout", h.logout)
		r.Group(func(r chi.Router) {
			if h.cfg.LoginRateLimit > 0 {
				// RemoteAddr is the client IP: the router replaces it with
				// X-Real-IP behind the website's proxy (TRUST_PROXY_HEADERS).
				r.Use(httprate.LimitBy(h.cfg.LoginRateLimit, time.Minute, clientIPKey,
					httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
						httpx.Error(w, r, ErrTooManyRequests)
					}),
				))
			}
			r.Get("/twitch", h.startLogin)
			r.Get("/twitch/callback", h.finishLogin)
			r.Get("/bot", h.startBot)
			r.Get("/bot/callback", h.finishBot)
			if h.cfg.DevLogin {
				r.Get("/dev", h.devLogin)
			}
		})
	})
}

// startLogin redirects to Twitch. The only scope is channel:bot.
//
//	@Summary		Log in with Twitch
//	@Description	Browser redirect to Twitch. The only scope requested is channel:bot, which lets the WheelGPT bot account read and write chat in your channel.
//	@Tags			auth
//	@Param			returnTo	query	string	false	"site path to return to"	default(/dashboard)
//	@Success		302
//	@Failure		503	{object}	httpx.Problem
//	@Router			/auth/twitch [get]
func (h *Handler) startLogin(w http.ResponseWriter, r *http.Request) {
	h.redirectToTwitch(w, r, purposeLogin, localPath(r.URL.Query().Get("returnTo"), "/dashboard"))
}

// finishLogin is the Twitch callback; it sets the session cookie.
//
//	@Summary	Twitch login callback
//	@Tags		auth
//	@Param		code	query	string	false	"authorization code"
//	@Param		state	query	string	false	"state"
//	@Success	303
//	@Failure	400	{object}	httpx.Problem
//	@Failure	502	{object}	httpx.Problem
//	@Router		/auth/twitch/callback [get]
func (h *Handler) finishLogin(w http.ResponseWriter, r *http.Request) {
	returnTo, ok := h.checkState(w, r, purposeLogin)
	if !ok {
		return
	}
	issued, err := h.service.Login(r.Context(), r.URL.Query().Get("code"), h.cfg.PublicOrigin+"/api/v1/auth/twitch/callback")
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	h.setSessionCookie(w, issued.Token, issued.ExpiresAt)
	http.Redirect(w, r, returnTo, http.StatusSeeOther) //nolint:gosec // localPath only allows same-site paths
}

// startBot redirects the operator to authorize the bot account.
//
//	@Summary		Authorize the bot account
//	@Description	One-time step for the operator: log in as the bot account to grant user:read:chat, user:write:chat and user:bot.
//	@Tags			auth
//	@Success		302
//	@Failure		503	{object}	httpx.Problem
//	@Router			/auth/bot [get]
func (h *Handler) startBot(w http.ResponseWriter, r *http.Request) {
	h.redirectToTwitch(w, r, purposeBot, "/")
}

// finishBot is the callback of the bot authorization.
//
//	@Summary	Bot authorization callback
//	@Tags		auth
//	@Success	303
//	@Failure	403	{object}	httpx.Problem
//	@Router		/auth/bot/callback [get]
func (h *Handler) finishBot(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.checkState(w, r, purposeBot); !ok {
		return
	}
	if err := h.service.AuthorizeBot(r.Context(), r.URL.Query().Get("code"), h.cfg.PublicOrigin+"/api/v1/auth/bot/callback"); err != nil {
		httpx.Error(w, r, err)
		return
	}
	http.Redirect(w, r, "/?bot=authorized", http.StatusSeeOther)
}

var devLoginPattern = regexp.MustCompile(`^[a-z0-9_]{3,25}$`)

// devLogin logs in as a fake channel. Only registered with DEV_MODE, which
// is only allowed for a localhost PUBLIC_URL.
func (h *Handler) devLogin(w http.ResponseWriter, r *http.Request) {
	login := strings.ToLower(r.URL.Query().Get("login"))
	if !devLoginPattern.MatchString(login) {
		httpx.Error(w, r, apperr.Invalid("invalid_login", "login must be 3 to 25 characters a-z, 0-9 or _"))
		return
	}
	issued, err := h.service.DevLogin(r.Context(), login)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	h.setSessionCookie(w, issued.Token, issued.ExpiresAt)
	target := localPath(r.URL.Query().Get("returnTo"), "/dashboard")
	http.Redirect(w, r, target, http.StatusSeeOther) //nolint:gosec // localPath only allows same-site paths
}

// logout ends the current session. It succeeds without a session too.
//
//	@Summary	Log out
//	@Tags		auth
//	@Success	204
//	@Router		/auth/logout [post]
func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(h.cookieName()); err == nil {
		if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
			httpx.Error(w, r, err)
			return
		}
	}
	h.clearSessionCookie(w)
	httpx.NoContent(w)
}

// ------------------------------------------------------------ OAuth state

func (h *Handler) redirectToTwitch(w http.ResponseWriter, r *http.Request, purpose, returnTo string) {
	state := RandomToken(24)
	callback := h.cfg.PublicOrigin + "/api/v1/auth/twitch/callback"
	if purpose == purposeBot {
		callback = h.cfg.PublicOrigin + "/api/v1/auth/bot/callback"
	}
	target, err := h.service.AuthorizeURL(callback, state, purpose == purposeBot)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	raw := strings.Join([]string{state, purpose, returnTo}, "\n")
	http.SetCookie(w, &http.Cookie{ //nolint:gosec // Secure follows the PUBLIC_URL scheme, false only for local http dev
		Name:     stateCookie,
		Value:    base64.RawURLEncoding.EncodeToString([]byte(raw)),
		Path:     stateCookiePath,
		MaxAge:   600,
		HttpOnly: true,
		Secure:   h.cfg.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, target, http.StatusFound)
}

// checkState reads and clears the state cookie. A denied login sends the
// browser back to the start page.
func (h *Handler) checkState(w http.ResponseWriter, r *http.Request, purpose string) (string, bool) {
	http.SetCookie(w, &http.Cookie{ //nolint:gosec // see redirectToTwitch
		Name: stateCookie, Path: stateCookiePath, MaxAge: -1, HttpOnly: true,
		Secure: h.cfg.SecureCookies, SameSite: http.SameSiteLaxMode,
	})
	if r.URL.Query().Get("error") != "" {
		http.Redirect(w, r, "/?login=cancelled", http.StatusSeeOther)
		return "", false
	}
	c, err := r.Cookie(stateCookie)
	if err != nil {
		httpx.Error(w, r, ErrLoginState)
		return "", false
	}
	raw, err := base64.RawURLEncoding.DecodeString(c.Value)
	parts := strings.SplitN(string(raw), "\n", 3)
	if err != nil || len(parts) != 3 || parts[1] != purpose ||
		subtle.ConstantTimeCompare([]byte(parts[0]), []byte(r.URL.Query().Get("state"))) != 1 {
		httpx.Error(w, r, ErrLoginState)
		return "", false
	}
	return parts[2], true
}

// localPath accepts only same-site paths as redirect targets.
func localPath(p, fallback string) string {
	if !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") || strings.ContainsAny(p, "\\\r\n") {
		return fallback
	}
	return p
}

func clientIPKey(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	return httprate.CanonicalizeIP(ip), nil
}

// ------------------------------------------------------------ cookies

func (h *Handler) cookieName() string {
	if h.cfg.SecureCookies {
		return secureCookieName
	}
	return devCookieName
}

func (h *Handler) setSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{ //nolint:gosec // Secure follows the PUBLIC_URL scheme, false only for local http dev
		Name:     h.cookieName(),
		Value:    token,
		Path:     "/",
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
		HttpOnly: true,
		Secure:   h.cfg.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}

func (h *Handler) clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{ //nolint:gosec // see setSessionCookie
		Name:     h.cookieName(),
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.cfg.SecureCookies,
		SameSite: http.SameSiteLaxMode,
	})
}
