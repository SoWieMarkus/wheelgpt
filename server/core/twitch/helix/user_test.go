package helix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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

		userResponse := UserResponse{Data: users}
		response, _ := json.Marshal(userResponse)
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
			if len(users.Data) != tt.ExpectedLen {
				t.Errorf("expected %d users, got %d", tt.ExpectedLen, len(users.Data))
			}
		})
	}
}

func TestClient_GetUserByAccessToken(t *testing.T) {
	tests := []struct {
		Name             string
		UserAccessToken  *identity.UserAccessToken
		ExpectedResponse *UserResponse
		ExpectError      bool
	}{
		{
			Name: "Valid Access Token",
			UserAccessToken: &identity.UserAccessToken{
				AccessToken: "valid-access-token",
			},
			ExpectedResponse: &UserResponse{
				Data: []User{{
					ID:          "12345",
					Login:       "testuser",
					DisplayName: "Test User",
				}},
			},
			ExpectError: false,
		},
		{
			Name: "Invalid Access Token",
			UserAccessToken: &identity.UserAccessToken{
				AccessToken: "I don't exist",
			},
			ExpectedResponse: &UserResponse{
				Data: []User{},
			},
			ExpectError: false,
		},
		{
			Name:             "Nil Access Token",
			UserAccessToken:  nil,
			ExpectedResponse: nil,
			ExpectError:      true,
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

		// Return mock user response
		users := []User{}
		if r.Header.Get("Authorization") == "Bearer valid-access-token" {
			users = append(users, User{
				ID:          "12345",
				Login:       "testuser",
				DisplayName: "Test User",
			})
		}

		response, _ := json.Marshal(&UserResponse{Data: users})
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
			users, err := client.GetUserByAccessToken(tt.UserAccessToken)
			if tt.ExpectError != (err != nil) {
				t.Errorf("expected error: %v, got error: %v", tt.ExpectError, err)
			}

			if !reflect.DeepEqual(users, tt.ExpectedResponse) {
				t.Errorf("expected response: %+v, got: %+v", tt.ExpectedResponse, users)
			}
		})
	}
}
