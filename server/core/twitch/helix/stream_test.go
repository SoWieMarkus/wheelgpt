package helix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	httpClient "github.com/SoWieMarkus/wheelgpt/core/http"
)

func TestClient_GetStreams(t *testing.T) {
	tests := []struct {
		Name        string
		ExpectError bool
		UserIDs     []string
		ExpectedLen int
	}{
		{
			Name:        "Valid Response",
			UserIDs:     []string{"12345", "67890"},
			ExpectError: false,
			ExpectedLen: 2,
		},
		{
			Name:        "Empty Response",
			UserIDs:     []string{"empty"},
			ExpectError: false,
			ExpectedLen: 0,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.Header.Get("Client-ID") != "test-id" {
			t.Errorf("expected Client-ID 'test-id', got '%s'", r.Header.Get("Client-ID"))
		}
		if r.Header.Get("Authorization") != "Bearer mock-access-token" {
			t.Errorf("expected Authorization 'Bearer mock-access-token', got '%s'", r.Header.Get("Authorization"))
		}

		userIds := r.URL.Query()["user_id"]

		var streams []Stream
		for _, id := range userIds {
			if id == "empty" {
				continue
			}

			streams = append(streams, Stream{
				ID:          "stream" + id,
				UserID:      id,
				UserLogin:   "user" + id,
				UserName:    "User " + id,
				GameID:      "game" + id,
				GameName:    "Game " + id,
				Type:        "live",
				Title:       "Stream Title " + id,
				ViewerCount: 100 + len(id),
				StartedAt:   time.Now().Add(-1 * time.Hour),
			})
		}
		response := &StreamResponse{
			Data: streams,
			Pagination: Pagination{
				Cursor: "mock-cursor",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	identityAPI := &mockIdentityAPI{
		returnError: false,
	}
	client := &Client{
		config:   &config.TwitchConfig{ClientID: "test-id", ClientSecret: "test-secret"},
		identity: identityAPI,
		client:   httpClient.NewClient(server.URL, nil),
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			streams, err := client.GetStreams(tt.UserIDs)
			if tt.ExpectError && err == nil {
				t.Errorf("expected error, got none")
			}
			if !tt.ExpectError && err != nil {
				t.Errorf("did not expect error, got %v", err)
			}
			if len(streams.Data) != tt.ExpectedLen {
				t.Errorf("expected %d streams, got %d", tt.ExpectedLen, len(streams.Data))
			}
		})
	}
}
