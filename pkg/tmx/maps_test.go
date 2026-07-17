package tmx

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/pkg/testutil"
)

type contextKey string

func newClientWithTransport(transport http.RoundTripper) *Client {
	return &Client{
		http: testutil.Client(transport),
	}
}

func TestClient_GetMapByUID(t *testing.T) {
	t.Parallel()

	const (
		mapUID = "kZVq2Lw"
	)

	contextErr := errors.New("transport failure")
	const key contextKey = "request-id"

	tests := []struct {
		name              string
		expectedURL       string
		ctx               context.Context
		response          *http.Response
		transportErr      error
		wantErr           bool
		wantErrIs         error
		assertMap         func(*testing.T, *Map)
		expectContextKey  bool
		expectedContextID string
	}{
		{
			name:        "success",
			expectedURL: "https://trackmania.exchange/api/maps/get_map_info/uid/" + mapUID,
			ctx:         context.Background(),
			response:    testutil.JSONResponse(http.StatusOK, `{"track_id":"12345"}`),
			assertMap: func(t *testing.T, m *Map) {
				t.Helper()
				if m == nil {
					t.Fatal("expected map, got nil")
				}
				if m.TrackID != "12345" {
					t.Fatalf("expected TrackID %q, got %q", "12345", m.TrackID)
				}
			},
		},
		{
			name:         "transport error is returned",
			expectedURL:  "https://trackmania.exchange/api/maps/get_map_info/uid/" + mapUID,
			ctx:          context.Background(),
			transportErr: contextErr,
			wantErr:      true,
			wantErrIs:    contextErr,
			assertMap: func(t *testing.T, m *Map) {
				t.Helper()
				if m != nil {
					t.Fatalf("expected nil map, got %+v", m)
				}
			},
		},
		{
			name:        "invalid JSON returns error",
			expectedURL: "https://trackmania.exchange/api/maps/get_map_info/uid/" + mapUID,
			ctx:         context.Background(),
			response:    testutil.JSONResponse(http.StatusOK, `{"track_id":`),
			wantErr:     true,
			assertMap: func(t *testing.T, m *Map) {
				t.Helper()
				if m != nil {
					t.Fatalf("expected nil map, got %+v", m)
				}
			},
		},
		{
			name:        "propagates context",
			expectedURL: "https://trackmania.exchange/api/maps/get_map_info/uid/" + mapUID,
			ctx:         context.WithValue(context.Background(), key, "req-123"),
			response:    testutil.JSONResponse(http.StatusOK, `{"track_id":"12345"}`),
			assertMap: func(t *testing.T, m *Map) {
				t.Helper()
				if m == nil {
					t.Fatal("expected map, got nil")
				}
				if m.TrackID != "12345" {
					t.Fatalf("expected TrackID %q, got %q", "12345", m.TrackID)
				}
			},
			expectContextKey:  true,
			expectedContextID: "req-123",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newClientWithTransport(testutil.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				testutil.AssertRequest(t, req, testutil.RequestExpectation{
					Method: http.MethodGet,
					URL:    tc.expectedURL,
				})

				if tc.expectContextKey {
					got, _ := req.Context().Value(key).(string)
					if got != tc.expectedContextID {
						t.Fatalf("expected context value %q, got %q", tc.expectedContextID, got)
					}
				}

				if tc.transportErr != nil {
					return nil, tc.transportErr
				}
				return tc.response, nil
			}))

			m, err := client.GetMapByUID(tc.ctx, mapUID)
			testutil.Result(t, m, err, testutil.Outcome[*Map]{
				WantErr:     tc.wantErr,
				WantErrIs:   tc.wantErrIs,
				AssertValue: tc.assertMap,
			})
		})
	}
}
