package httpx

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/SoWieMarkus/wheelgpt/internal/domain"
)

// Problem is an RFC 9457 problem details response.
type Problem struct {
	Type     string `json:"type"     example:"about:blank"`
	Title    string `json:"title"    example:"Not Found"`
	Status   int    `json:"status"   example:"404"`
	Detail   string `json:"detail"   example:"this channel does not use WheelGPT"`
	Code     string `json:"code"     example:"channel_not_found"`
	Instance string `json:"instance" example:"/api/v1/channels/someone"`
	// RequestID correlates the response with server logs.
	RequestID string       `json:"requestId,omitempty" binding:"optional"`
	Errors    []FieldError `json:"errors,omitempty" binding:"optional"`
}

type FieldError struct {
	Field   string `json:"field"   example:"guessDelaySeconds"`
	Message string `json:"message" example:"is required"`
}

// requestError is a problem with the request itself (bad JSON, failed
// validation, malformed path parameter) detected by the HTTP layer.
type requestError struct {
	status int
	code   string
	detail string
	fields []FieldError
}

func (e *requestError) Error() string {
	return e.code + ": " + e.detail
}

// BadRequest creates a 400 problem for invalid query or form input.
func BadRequest(code, detail string) error {
	return &requestError{status: http.StatusBadRequest, code: code, detail: detail}
}

var statusByKind = map[domain.Kind]int{
	domain.KindInvalid:         http.StatusBadRequest,
	domain.KindNotFound:        http.StatusNotFound,
	domain.KindConflict:        http.StatusConflict,
	domain.KindUnprocessable:   http.StatusUnprocessableEntity,
	domain.KindUnauthenticated: http.StatusUnauthorized,
	domain.KindForbidden:       http.StatusForbidden,
	domain.KindTooManyRequests: http.StatusTooManyRequests,
	domain.KindBadGateway:      http.StatusBadGateway,
	domain.KindUnavailable:     http.StatusServiceUnavailable,
}

// Error translates any error into a problem response. Business errors from
// the domain package map to 4xx by their Kind, unknown errors are logged and
// hidden behind a generic 500 so internals never leak.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	var reqErr *requestError
	if errors.As(err, &reqErr) {
		WriteProblem(w, r, reqErr.status, reqErr.code, reqErr.detail, reqErr.fields)
		return
	}

	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		status, ok := statusByKind[domainErr.Kind]
		if !ok {
			status = http.StatusInternalServerError
		}
		WriteProblem(w, r, status, domainErr.Code, domainErr.Message, nil)
		return
	}

	slog.ErrorContext(r.Context(), "unhandled error",
		"error", err,
		"method", r.Method,
		"path", r.URL.Path,
		"request_id", middleware.GetReqID(r.Context()),
	)
	WriteProblem(w, r, http.StatusInternalServerError, "internal_error", "something went wrong on our side", nil)
}

// WriteProblem writes a problem response with the given status and code.
func WriteProblem(w http.ResponseWriter, r *http.Request, status int, code, detail string, fields []FieldError) {
	w.Header().Set("Content-Type", "application/problem+json")
	encode(w, status, Problem{
		Type:      "about:blank",
		Title:     http.StatusText(status),
		Status:    status,
		Detail:    detail,
		Code:      code,
		Instance:  r.URL.Path,
		RequestID: middleware.GetReqID(r.Context()),
		Errors:    fields,
	})
}
