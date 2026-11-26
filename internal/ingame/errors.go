package ingame

import "github.com/SoWieMarkus/gokit/apperr"

// Errors returned by this package. The HTTP layer maps their Kind to a
// status code, the Code strings are part of the API.
var (
	ErrMapNotFound  = apperr.New(apperr.KindNotFound, "map_not_found", "the streamer is not on a map right now")
	ErrRoomNotFound = apperr.New(apperr.KindNotFound, "room_not_found", "the streamer is not in a room right now")
)
