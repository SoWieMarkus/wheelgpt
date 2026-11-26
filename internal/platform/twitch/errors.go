package twitch

import "github.com/SoWieMarkus/gokit/apperr"

// Errors returned by this package. The HTTP layer maps their Kind to a
// status code, the Code strings are part of the API.
var ErrTwitchFailed = apperr.New(apperr.KindBadGateway, "twitch_failed", "Twitch did not answer as expected, try again in a minute")
