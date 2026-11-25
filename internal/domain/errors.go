// Package domain holds transport independent business errors.
//
// Services return these errors and the HTTP layer translates their Kind into
// a status code, so no service ever has to know about HTTP.
package domain

import "fmt"

// Kind classifies an error for the transport layer.
type Kind int

const (
	KindInvalid Kind = iota + 1
	KindNotFound
	KindConflict
	KindUnprocessable
	KindUnauthenticated
	KindForbidden
	KindTooManyRequests
	KindBadGateway
	KindUnavailable
)

// Error is a business rule violation with a stable machine readable code.
type Error struct {
	Kind    Kind
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Invalid creates a validation error for input that breaks a business rule.
func Invalid(code, message string) *Error {
	return &Error{Kind: KindInvalid, Code: code, Message: message}
}

var (
	ErrUnauthenticated      = &Error{Kind: KindUnauthenticated, Code: "unauthenticated", Message: "log in with Twitch first"}
	ErrInvalidPluginToken   = &Error{Kind: KindUnauthenticated, Code: "invalid_plugin_token", Message: "send a valid plugin token as 'Authorization: Bearer <token>', get one on the dashboard"}
	ErrCrossOrigin          = &Error{Kind: KindForbidden, Code: "cross_origin", Message: "cross origin requests are not allowed"}
	ErrWrongBotAccount      = &Error{Kind: KindForbidden, Code: "wrong_bot_account", Message: "log in as the configured bot account"}
	ErrInvalidSignature     = &Error{Kind: KindForbidden, Code: "invalid_signature", Message: "invalid EventSub signature"}
	ErrChannelNotFound      = &Error{Kind: KindNotFound, Code: "channel_not_found", Message: "this channel does not use WheelGPT"}
	ErrMapNotFound          = &Error{Kind: KindNotFound, Code: "map_not_found", Message: "the streamer is not on a map right now"}
	ErrRoomNotFound         = &Error{Kind: KindNotFound, Code: "room_not_found", Message: "the streamer is not in a room right now"}
	ErrChannelOffline       = &Error{Kind: KindConflict, Code: "channel_offline", Message: "the stream is not live, data is only accepted while live"}
	ErrInvalidPeriod        = &Error{Kind: KindInvalid, Code: "invalid_period", Message: "period must be all-time or a month like 2026-11"}
	ErrLoginState           = &Error{Kind: KindInvalid, Code: "invalid_login_state", Message: "the login took too long or cookies are blocked, please try again"}
	ErrTooManyRequests      = &Error{Kind: KindTooManyRequests, Code: "too_many_requests", Message: "too many requests, try again in a minute"}
	ErrTwitchFailed         = &Error{Kind: KindBadGateway, Code: "twitch_failed", Message: "Twitch did not answer as expected, try again in a minute"}
	ErrTwitchNotConfigured  = &Error{Kind: KindUnavailable, Code: "twitch_not_configured", Message: "Twitch login is not configured on this server"}
	ErrChatAccessNotGranted = &Error{Kind: KindForbidden, Code: "chat_access_not_granted", Message: "Twitch refused access to the chat: log in again or make the bot a moderator"}
)
