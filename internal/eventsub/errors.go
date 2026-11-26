package eventsub

import "github.com/SoWieMarkus/gokit/apperr"

// Errors returned by this package. The HTTP layer maps their Kind to a
// status code, the Code strings are part of the API.
var ErrInvalidSignature = apperr.New(apperr.KindForbidden, "invalid_signature", "invalid EventSub signature")
