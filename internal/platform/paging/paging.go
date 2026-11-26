// Package paging reads limit and offset query parameters and wraps one page
// of a longer collection for the response.
package paging

import (
	"net/http"

	"github.com/SoWieMarkus/gokit/httpx"
)

// Page is one page of a longer collection. Total counts all items that
// match, so clients can render page numbers.
type Page[T any] struct {
	Items  []T   `json:"items"`
	Total  int64 `json:"total"  example:"120"`
	Limit  int64 `json:"limit"  example:"25"`
	Offset int64 `json:"offset" example:"0"`
}

// PageDoc describes Page in the OpenAPI spec, see httpx.ListDoc.
type PageDoc struct {
	Items  []any `json:"items"`
	Total  int64 `json:"total"  example:"120"`
	Limit  int64 `json:"limit"  example:"25"`
	Offset int64 `json:"offset" example:"0"`
} //	@name	paging.PageDoc

// Parse reads the limit and offset query parameters.
func Parse(r *http.Request, defaultLimit int64) (limit, offset int64, err error) {
	if limit, err = httpx.QueryInt(r, "limit", defaultLimit, 1, 100); err != nil {
		return 0, 0, err
	}
	if offset, err = httpx.QueryInt(r, "offset", 0, 0, 1_000_000); err != nil {
		return 0, 0, err
	}
	return limit, offset, nil
}
