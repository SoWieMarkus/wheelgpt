package channel

import "github.com/SoWieMarkus/gokit/apperr"

// Errors returned by this package. The HTTP layer maps their Kind to a
// status code, the Code strings are part of the API.
var ErrChannelNotFound = apperr.New(apperr.KindNotFound, "channel_not_found", "this channel does not use WheelGPT")
