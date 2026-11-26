package plugin

import "github.com/SoWieMarkus/gokit/apperr"

// Errors returned by this package. The HTTP layer maps their Kind to a
// status code, the Code strings are part of the API.
var (
	ErrInvalidPluginToken = apperr.New(apperr.KindUnauthenticated, "invalid_plugin_token", "send a valid plugin token as 'Authorization: Bearer <token>', get one on the dashboard")
	ErrChannelOffline     = apperr.New(apperr.KindConflict, "channel_offline", "the stream is not live, data is only accepted while live")
)
