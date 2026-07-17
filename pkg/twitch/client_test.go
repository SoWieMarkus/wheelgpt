package twitch

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/pkg/testutil"
)

func newClientWithTransport(clientID string, transport http.RoundTripper) *Client {
	return &Client{
		clientID: clientID,
		http:     testutil.Client(transport),
	}
}

func TestClient_GetUser(t *testing.T) {
	t.Parallel()

	const (
		clientID = "client-id"
		token    = "user-token"
	)

	transportErrNetwork := errors.New("network down")

	tests := []struct {
		name         string
		expectedReq  testutil.RequestExpectation
		response     *http.Response
		transportErr error
		wantErr      bool
		wantErrIs    error
		assertUser   func(*testing.T, *User)
	}{
		{
			name: "success returns first user",
			expectedReq: testutil.RequestExpectation{
				Method: http.MethodGet,
				URL:    "https://api.twitch.tv/helix/users",
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
					"Client-Id":     clientID,
				},
			},
			response: testutil.JSONResponse(http.StatusOK, `{"data":[{"id":"42","login":"foo","display_name":"Foo","profile_image_url":"https://img"}]}`),
			assertUser: func(t *testing.T, user *User) {
				t.Helper()
				if user == nil {
					t.Fatal("expected user, got nil")
				}
				if user.ID != "42" {
					t.Fatalf("expected user ID %q, got %q", "42", user.ID)
				}
			},
		},
		{
			name: "transport error is returned",
			expectedReq: testutil.RequestExpectation{
				Method: http.MethodGet,
				URL:    "https://api.twitch.tv/helix/users",
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
					"Client-Id":     clientID,
				},
			},
			transportErr: transportErrNetwork,
			wantErr:      true,
			wantErrIs:    transportErrNetwork,
			assertUser: func(t *testing.T, user *User) {
				t.Helper()
				if user != nil {
					t.Fatalf("expected nil user, got %+v", user)
				}
			},
		},
		{
			name: "invalid JSON returns error",
			expectedReq: testutil.RequestExpectation{
				Method: http.MethodGet,
				URL:    "https://api.twitch.tv/helix/users",
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
					"Client-Id":     clientID,
				},
			},
			response: testutil.JSONResponse(http.StatusOK, `{"data":`),
			wantErr:  true,
			assertUser: func(t *testing.T, user *User) {
				t.Helper()
				if user != nil {
					t.Fatalf("expected nil user, got %+v", user)
				}
			},
		},
		{
			name: "empty data returns error",
			expectedReq: testutil.RequestExpectation{
				Method: http.MethodGet,
				URL:    "https://api.twitch.tv/helix/users",
				Headers: map[string]string{
					"Authorization": "Bearer " + token,
					"Client-Id":     clientID,
				},
			},
			response: testutil.JSONResponse(http.StatusOK, `{"data":[]}`),
			wantErr:  true,
			assertUser: func(t *testing.T, user *User) {
				t.Helper()
				if user != nil {
					t.Fatalf("expected nil user, got %+v", user)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newClientWithTransport(clientID, testutil.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				testutil.AssertRequest(t, req, tc.expectedReq)

				if tc.transportErr != nil {
					return nil, tc.transportErr
				}
				return tc.response, nil
			}))

			user, err := client.GetUser(context.Background(), token)
			testutil.Result(t, user, err, testutil.Outcome[*User]{
				WantErr:     tc.wantErr,
				WantErrIs:   tc.wantErrIs,
				AssertValue: tc.assertUser,
			})
		})
	}
}

func TestClient_GetUsers(t *testing.T) {
	t.Parallel()

	t.Run("empty input returns empty without request", func(t *testing.T) {
		t.Parallel()

		called := false
		client := newClientWithTransport("client-id", testutil.RoundTripFunc(func(*http.Request) (*http.Response, error) {
			called = true
			return nil, errors.New("should not be called")
		}))

		users, err := client.GetUsers(context.Background(), nil)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 0 {
			t.Fatalf("expected no users, got %d", len(users))
		}
		if called {
			t.Fatal("expected no HTTP request for empty input")
		}
	})

	t.Run("chunks ids in batches of 100", func(t *testing.T) {
		t.Parallel()

		ids := make([]string, 101)
		for i := range ids {
			ids[i] = "id"
		}

		requestCount := 0
		client := newClientWithTransport("client-id", testutil.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			requestCount++
			testutil.AssertRequest(t, req, testutil.RequestExpectation{
				Method: http.MethodGet,
				Host:   "api.twitch.tv",
				Path:   "/helix/users",
				Headers: map[string]string{
					"Client-Id": "client-id",
				},
			})

			queryIDs := req.URL.Query()["id"]
			if requestCount == 1 && len(queryIDs) != 100 {
				t.Fatalf("first request should contain 100 ids, got %d", len(queryIDs))
			}
			if requestCount == 2 && len(queryIDs) != 1 {
				t.Fatalf("second request should contain 1 id, got %d", len(queryIDs))
			}

			usersToReturn := make([]User, len(queryIDs))
			for idx, id := range queryIDs {
				usersToReturn[idx] = User{
					ID:              id,
					Login:           "login",
					DisplayName:     "Display Name",
					ProfileImageURL: "https://img",
				}
			}

			usersJson, err := json.Marshal(usersToReturn)
			if err != nil {
				t.Fatalf("failed to marshal users: %v", err)
			}

			return testutil.JSONResponse(http.StatusOK, `{"data":`+string(usersJson)+`}`), nil
		}))

		users, err := client.GetUsers(context.Background(), ids)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if requestCount != 2 {
			t.Fatalf("expected 2 requests, got %d", requestCount)
		}
		if len(users) != 101 {
			t.Fatalf("expected 101 users, got %d", len(users))
		}
	})
}
