package helix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	httpClient "github.com/SoWieMarkus/wheelgpt/core/http"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/identity"
)

func TestClient_GetUsers(t *testing.T) {
	tests := []struct {
		Name        string
		UserIds     []string
		ExpectError bool
		ExpectedLen int
	}{
		{
			Name:        "Valid User IDs",
			UserIds:     []string{"12345", "67890"},
			ExpectError: false,
			ExpectedLen: 2,
		},
		{
			Name:        "Empty User IDs",
			UserIds:     []string{},
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

		// Parse query parameters
		query := r.URL.Query()
		ids := query["id"]

		// Create mock users based on requested IDs
		var users []User
		for _, id := range ids {
			users = append(users, User{
				ID:          id,
				Login:       "user" + id,
				DisplayName: "User " + id,
			})
		}

		response, _ := json.Marshal(users)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
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
			users, err := client.GetUsers(tt.UserIds)

			if tt.ExpectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.ExpectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
			if len(users) != tt.ExpectedLen {
				t.Errorf("expected %d users, got %d", tt.ExpectedLen, len(users))
			}
		})
	}
}

func TestClient_GetUserByAccessToken(t *testing.T) {
	tests := []struct {
		Name            string
		UserAccessToken *identity.UserAccessToken
		ExpectError     bool
	}{
		{
			Name: "Valid Access Token",
			UserAccessToken: &identity.UserAccessToken{
				AccessToken: "valid-access-token",
			},
			ExpectError: false,
		},
		{
			Name: "Invalid Access Token",
			UserAccessToken: &identity.UserAccessToken{
				AccessToken: "I don't exist",
			},
			ExpectError: true,
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
		if r.Header.Get("Authorization") != "Bearer valid-access-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return mock user response
		user := User{
			ID:          "12345",
			Login:       "testuser",
			DisplayName: "Test User",
		}
		response, _ := json.Marshal([]User{user})
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
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
			user, err := client.GetUserByAccessToken(tt.UserAccessToken)
			if tt.ExpectError != (err != nil) {
				t.Errorf("expected error: %v, got error: %v", tt.ExpectError, err)
			}
			if !tt.ExpectError && user == nil {
				t.Errorf("expected user but got nil")
			}
			if !tt.ExpectError && user != nil && user.ID != "12345" {
				t.Errorf("expected user ID '12345', got '%s'", user.ID)
			}
		})
	}
}
