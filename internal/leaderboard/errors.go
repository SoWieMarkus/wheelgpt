package leaderboard

import "github.com/SoWieMarkus/gokit/apperr"

// Errors returned by this package. The HTTP layer maps their Kind to a
// status code, the Code strings are part of the API.
var ErrInvalidPeriod = apperr.New(apperr.KindInvalid, "invalid_period", "period must be all-time or a month like 2026-11")
