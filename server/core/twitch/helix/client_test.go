package helix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/identity"
)

func createMockIdentityServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &identity.AppAccessToken{
			AccessToken: "mock-access-token",
			ExpiresIn:   3600,
			TokenType:   "bearer",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
}

// Creste mock identity server that returns an error
func createMockIdentityServerWithError() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server_error"}`))
	}))
}

func createMockHelixClient(server *httptest.Server) *Client {
	client := &Client{
		client:  &http.Client{},
		baseURL: "https://api.twitch.tv/helix",
		config:  &config.TwitchConfig{ClientID: "test-id", ClientSecret: "test-secret"},
		identity: identity.NewClient(&config.TwitchConfig{
			ClientID:     "test-id",
			ClientSecret: "test-secret",
		}),
	}
	client.identity.BaseURL = server.URL
	return client
}

func TestClient_getAppAccessToken(t *testing.T) {
	tests := []struct {
		Name                string
		SetupServer         func() *httptest.Server
		CachedToken         *identity.AppAccessToken
		CachedTokenIssuedAt *time.Time
		ExpectError         bool
		ExpectedToken       string
	}{
		{
			Name:        "No Cached Token - Fetch New",
			CachedToken: nil,
			SetupServer: func() *httptest.Server {
				return createMockIdentityServer()
			},
			ExpectError:   false,
			ExpectedToken: "mock-access-token",
		},
		{
			Name: "Valid Cached Token - Use Cached",
			CachedToken: &identity.AppAccessToken{
				AccessToken: "cached-access-token",
				ExpiresIn:   3600,
				TokenType:   "bearer",
			},
			CachedTokenIssuedAt: func() *time.Time {
				t := time.Now().Add(-10 * time.Minute)
				return &t
			}(),
			SetupServer: func() *httptest.Server {
				return createMockIdentityServer()
			},
			ExpectError:   false,
			ExpectedToken: "cached-access-token",
		},
		{
			Name: "Valid Cached Token - No IssuedAt - Fetch New",
			CachedToken: &identity.AppAccessToken{
				AccessToken: "cached-access-token",
				ExpiresIn:   3600,
				TokenType:   "bearer",
			},
			CachedTokenIssuedAt: nil,
			SetupServer: func() *httptest.Server {
				return createMockIdentityServer()
			},
			ExpectError:   false,
			ExpectedToken: "mock-access-token",
		},
		{
			Name: "Expired Cached Token - Fetch New",
			CachedToken: &identity.AppAccessToken{
				AccessToken: "expired-access-token",
				ExpiresIn:   3600,
				TokenType:   "bearer",
			},
			CachedTokenIssuedAt: func() *time.Time {
				t := time.Now().Add(-2 * time.Hour) // Issued 2 hours ago, so it's expired
				return &t
			}(),
			SetupServer: func() *httptest.Server {
				return createMockIdentityServer()
			},
			ExpectError:   false,
			ExpectedToken: "mock-access-token",
		},
		{
			Name: "Error Response",
			SetupServer: func() *httptest.Server {
				return createMockIdentityServerWithError()
			},
			ExpectError:   true,
			ExpectedToken: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			server := tt.SetupServer()
			defer server.Close()

			client := createMockHelixClient(server)
			client.appToken = tt.CachedToken
			client.appTokenIssuedAt = tt.CachedTokenIssuedAt
			token, err := client.getAppToken()

			if tt.ExpectError != (err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.ExpectError, err != nil)
			}

			if token != tt.ExpectedToken {
				t.Errorf("expected token: %v, got: %v", tt.ExpectedToken, token)
			}
		})
	}
}
