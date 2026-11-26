// Package slug reads path parameters made of letters, digits, '_' and '-',
// like Twitch logins and ids or leaderboard period keys.
package slug

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"

	"github.com/SoWieMarkus/gokit/httpx"
)

var pattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

// FromPath reads the path parameter name and rejects anything that is not a
// slug with a 400 problem.
func FromPath(r *http.Request, name string) (string, error) {
	v := chi.URLParam(r, name)
	if !pattern.MatchString(v) {
		return "", httpx.BadRequest("invalid_path", name+" is not valid")
	}
	return v, nil
}
