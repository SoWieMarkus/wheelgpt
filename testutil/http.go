package testutil

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// RoundTripFunc lets tests provide inline HTTP transport behavior.
type RoundTripFunc func(*http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Client returns an HTTP client that uses the provided fake transport.
func Client(transport http.RoundTripper) *http.Client {
	return &http.Client{Transport: transport}
}

// JSONResponse returns an HTTP response with a JSON body.
func JSONResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// RequestExpectation describes the expected HTTP request attributes.
type RequestExpectation struct {
	Method  string
	URL     string
	Host    string
	Path    string
	Headers map[string]string
}

// AssertRequest validates method, URL/host/path, and selected headers.
func AssertRequest(t *testing.T, req *http.Request, expected RequestExpectation) {
	t.Helper()

	if expected.Method != "" && req.Method != expected.Method {
		t.Fatalf("expected method %q, got %q", expected.Method, req.Method)
	}
	if expected.URL != "" && req.URL.String() != expected.URL {
		t.Fatalf("expected URL %q, got %q", expected.URL, req.URL.String())
	}
	if expected.Host != "" && req.URL.Host != expected.Host {
		t.Fatalf("expected host %q, got %q", expected.Host, req.URL.Host)
	}
	if expected.Path != "" && req.URL.Path != expected.Path {
		t.Fatalf("expected path %q, got %q", expected.Path, req.URL.Path)
	}

	for name, want := range expected.Headers {
		if got := req.Header.Get(name); got != want {
			t.Fatalf("expected header %q to be %q, got %q", name, want, got)
		}
	}
}

// Outcome describes the expected result of a function call in tests.
type Outcome[T any] struct {
	WantErr     bool
	WantErrIs   error
	AssertValue func(*testing.T, T)
}

// Result validates error expectations and optionally validates the returned value.
func Result[T any](t *testing.T, got T, err error, outcome Outcome[T]) {
	t.Helper()

	if outcome.WantErr && err == nil {
		t.Fatal("expected error, got nil")
	}
	if !outcome.WantErr && err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if outcome.WantErrIs != nil && !errors.Is(err, outcome.WantErrIs) {
		t.Fatalf("expected wrapped error %v, got %v", outcome.WantErrIs, err)
	}

	if outcome.AssertValue != nil {
		outcome.AssertValue(t, got)
	}
}
