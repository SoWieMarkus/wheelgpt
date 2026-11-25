// Package httpx is the small shared toolkit every HTTP handler uses: JSON
// responses, request decoding with validation, path and query parameters and
// RFC 9457 problem responses. It knows nothing about any feature.
package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

const maxBodyBytes = 1 << 20

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Report JSON field names instead of Go field names.
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.Split(f.Tag.Get("json"), ",")[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return v
}

// List wraps collections so fields like pagination can be added later
// without breaking clients.
type List[T any] struct {
	Items []T `json:"items"`
}

// Page is one page of a longer collection. Total counts all items that
// match, so clients can render page numbers.
type Page[T any] struct {
	Items  []T   `json:"items"`
	Total  int64 `json:"total"  example:"120"`
	Limit  int64 `json:"limit"  example:"25"`
	Offset int64 `json:"offset" example:"0"`
}

// PageDoc describes Page in the OpenAPI spec, see ListDoc.
type PageDoc struct {
	Items  []any `json:"items"`
	Total  int64 `json:"total"  example:"120"`
	Limit  int64 `json:"limit"  example:"25"`
	Offset int64 `json:"offset" example:"0"`
}

// ListDoc describes List in the OpenAPI spec. swag cannot resolve a generic
// instance whose type argument lives in the handler's own package, so
// handlers annotate lists as `httpx.ListDoc{items=[]Response}` instead.
type ListDoc struct {
	Items []any `json:"items"`
}

// JSON writes body as a JSON response.
func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	encode(w, status, body)
}

// NoContent writes an empty 204 response.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func encode(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("encode response", "error", err)
	}
}

// Decode reads a JSON body into T and validates it. Unknown fields are
// rejected so typos in clients surface immediately.
func Decode[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var body T
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return body, &requestError{
			status: http.StatusUnsupportedMediaType,
			code:   "unsupported_media_type",
			detail: "request body must be application/json",
		}
	}

	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&body); err != nil {
		return body, &requestError{status: http.StatusBadRequest, code: "invalid_body", detail: describeDecodeError(err)}
	}
	if dec.More() {
		return body, &requestError{status: http.StatusBadRequest, code: "invalid_body", detail: "body must contain a single JSON object"}
	}
	return body, Validate(body)
}

// Validate runs the struct validation rules and turns failures into a 400
// problem with one entry per field.
func Validate(v any) error {
	var invalid validator.ValidationErrors
	if errors.As(validate.Struct(v), &invalid) {
		fields := make([]FieldError, len(invalid))
		for i, fe := range invalid {
			fields[i] = FieldError{Field: fieldPath(fe), Message: describeRule(fe)}
		}
		return &requestError{
			status: http.StatusBadRequest,
			code:   "validation_failed",
			detail: "some fields are invalid",
			fields: fields,
		}
	}
	return nil
}

// fieldPath drops the struct name from the namespace, "Body.tags[0].name"
// becomes "tags[0].name".
func fieldPath(fe validator.FieldError) string {
	ns := fe.Namespace()
	if i := strings.Index(ns, "."); i >= 0 {
		return ns[i+1:]
	}
	return fe.Field()
}

func describeDecodeError(err error) string {
	var syntaxErr *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	var tooLarge *http.MaxBytesError
	switch {
	case errors.Is(err, io.EOF):
		return "request body must not be empty"
	case errors.As(err, &syntaxErr):
		return fmt.Sprintf("malformed JSON at position %d", syntaxErr.Offset)
	case errors.As(err, &typeErr):
		return fmt.Sprintf("field %q has the wrong type", typeErr.Field)
	case errors.As(err, &tooLarge):
		return "request body is too large"
	case strings.HasPrefix(err.Error(), "json: unknown field"):
		return strings.TrimPrefix(err.Error(), "json: ")
	}
	return "malformed JSON"
}

func describeRule(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "oneof":
		return "must be one of: " + fe.Param()
	case "min":
		if fe.Kind() == reflect.String {
			return "must be at least " + fe.Param() + " characters long"
		}
		if fe.Kind() == reflect.Slice {
			return "must contain at least " + fe.Param() + " entries"
		}
		return "must be at least " + fe.Param()
	case "max":
		if fe.Kind() == reflect.String {
			return "must be at most " + fe.Param() + " characters long"
		}
		if fe.Kind() == reflect.Slice {
			return "must contain at most " + fe.Param() + " entries"
		}
		return "must be at most " + fe.Param()
	case "gt":
		return "must be greater than " + fe.Param()
	case "gte":
		return "must be at least " + fe.Param()
	case "lte":
		return "must be at most " + fe.Param()
	case "http_url", "url":
		return "must be an http or https URL"
	case "printascii":
		return "may only contain printable ASCII characters"
	case "alphanum":
		return "may only contain letters and digits"
	case "excludesall":
		return "contains characters that are not allowed"
	}
	return "is invalid"
}

// PathID parses a positive integer path parameter.
func PathID(r *http.Request, name string) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, name), 10, 64)
	if err != nil || id <= 0 {
		return 0, &requestError{status: http.StatusBadRequest, code: "invalid_id", detail: name + " must be a positive integer"}
	}
	return id, nil
}

var slugPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

// PathSlug reads a path parameter made of letters, digits, '_' and '-'
// (Twitch logins and ids, period keys).
func PathSlug(r *http.Request, name string) (string, error) {
	v := chi.URLParam(r, name)
	if !slugPattern.MatchString(v) {
		return "", &requestError{status: http.StatusBadRequest, code: "invalid_path", detail: name + " is not valid"}
	}
	return v, nil
}

// Pagination reads limit and offset query parameters.
func Pagination(r *http.Request, defaultLimit int64) (limit, offset int64, err error) {
	if limit, err = QueryInt(r, "limit", defaultLimit, 1, 100); err != nil {
		return 0, 0, err
	}
	if offset, err = QueryInt(r, "offset", 0, 0, 1_000_000); err != nil {
		return 0, 0, err
	}
	return limit, offset, nil
}

// QueryInt parses an optional integer query parameter within [min, max].
func QueryInt(r *http.Request, name string, def, minimum, maximum int64) (int64, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return def, nil
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || v < minimum || v > maximum {
		return 0, BadRequest("invalid_query", fmt.Sprintf("%s must be an integer between %d and %d", name, minimum, maximum))
	}
	return v, nil
}

// QueryIDs parses a repeated positive integer query parameter (?tag=1&tag=2).
func QueryIDs(r *http.Request, name string) ([]int64, error) {
	raw := r.URL.Query()[name]
	ids := make([]int64, 0, len(raw))
	for _, s := range raw {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil || id <= 0 {
			return nil, BadRequest("invalid_query", name+" must be a positive integer")
		}
		ids = append(ids, id)
	}
	if len(ids) > 50 {
		return nil, BadRequest("invalid_query", "too many "+name+" values")
	}
	return ids, nil
}

// QueryBool parses an optional boolean query parameter.
func QueryBool(r *http.Request, name string) (bool, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return false, nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return false, BadRequest("invalid_query", name+" must be true or false")
	}
	return v, nil
}

// QueryEnum reads a query parameter that must be one of allowed. The first
// allowed value is the default.
func QueryEnum(r *http.Request, name string, allowed ...string) (string, error) {
	raw := r.URL.Query().Get(name)
	if raw == "" {
		return allowed[0], nil
	}
	for _, a := range allowed {
		if raw == a {
			return raw, nil
		}
	}
	return "", BadRequest("invalid_query", name+" must be one of: "+strings.Join(allowed, ", "))
}
