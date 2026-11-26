package auth

import "github.com/SoWieMarkus/gokit/apperr"

// Errors returned by this package. The HTTP layer maps their Kind to a
// status code, the Code strings are part of the API.
var (
	ErrUnauthenticated     = apperr.New(apperr.KindUnauthenticated, "unauthenticated", "log in with Twitch first")
	ErrWrongBotAccount     = apperr.New(apperr.KindForbidden, "wrong_bot_account", "log in as the configured bot account")
	ErrLoginState          = apperr.New(apperr.KindInvalid, "invalid_login_state", "the login took too long or cookies are blocked, please try again")
	ErrTooManyRequests     = apperr.New(apperr.KindTooManyRequests, "too_many_requests", "too many requests, try again in a minute")
	ErrTwitchNotConfigured = apperr.New(apperr.KindUnavailable, "twitch_not_configured", "Twitch login is not configured on this server")
)
