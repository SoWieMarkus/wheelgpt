package helix

import (
	"fmt"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/twitch/identity"
)

type mockIdentityAPI struct {
	returnError bool
}

func (m *mockIdentityAPI) RequestAppAccessToken() (*identity.AppAccessToken, error) {
	if m.returnError {
		return nil, fmt.Errorf("mock error")
	}
	return &identity.AppAccessToken{
		AccessToken: "mock-access-token",
		ExpiresIn:   3600,
		TokenType:   "bearer",
	}, nil
}

func (m *mockIdentityAPI) RequestUserAccessToken(code, redirectURI string) (*identity.UserAccessToken, error) {
	if m.returnError {
		return nil, fmt.Errorf("mock error")
	}
	return &identity.UserAccessToken{
		AccessToken:  "mock-user-token",
		ExpiresIn:    14400,
		RefreshToken: "mock-refresh-token",
		TokenType:    "bearer",
		Scope:        []string{"user:read:email"},
	}, nil
}

func createMockHelixClient(identityAPI identity.IdentityAPI) *Client {
	client := &Client{
		config:   &config.TwitchConfig{ClientID: "test-id", ClientSecret: "test-secret"},
		identity: identityAPI,
	}
	return client
}

func TestClient_getAppAccessToken(t *testing.T) {
	tests := []struct {
		Name                string
		IdentityAPI         identity.IdentityAPI
		CachedToken         *identity.AppAccessToken
		CachedTokenIssuedAt *time.Time
		ExpectError         bool
		ExpectedToken       string
	}{
		{
			Name:        "No Cached Token - Fetch New",
			CachedToken: nil,
			IdentityAPI: &mockIdentityAPI{
				returnError: false,
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
			IdentityAPI: &mockIdentityAPI{
				returnError: false,
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
			IdentityAPI: &mockIdentityAPI{
				returnError: false,
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
			IdentityAPI: &mockIdentityAPI{
				returnError: false,
			},
			ExpectError:   false,
			ExpectedToken: "mock-access-token",
		},
		{
			Name: "Error Response",
			IdentityAPI: &mockIdentityAPI{
				returnError: true,
			},
			ExpectError:   true,
			ExpectedToken: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			client := createMockHelixClient(tt.IdentityAPI)
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
