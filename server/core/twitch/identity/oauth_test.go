package identity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/SoWieMarkus/wheelgpt/core/config"
)

func TestClient_RequestAppAccessToken(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and content type
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected application/x-www-form-urlencoded content type")
		}

		// Parse and verify form data
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
		}

		// Verify OAuth parameters
		expectedParams := map[string]string{
			"client_id":     "test-client-id",
			"client_secret": "test-client-secret",
			"grant_type":    "client_credentials",
		}

		for key, expected := range expectedParams {
			if actual := r.FormValue(key); actual != expected {
				t.Errorf("expected %s '%s', got '%s'", key, expected, actual)
			}
		}

		// Return mock token response
		response := AppAccessToken{
			AccessToken: "test-app-access-token",
			ExpiresIn:   3600,
			TokenType:   "bearer",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := &Client{
		client:  http.DefaultClient,
		baseURL: server.URL,
		config: &config.TwitchConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	}

	token, err := client.RequestAppAccessToken()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedToken := &AppAccessToken{
		AccessToken: "test-app-access-token",
		ExpiresIn:   3600,
		TokenType:   "bearer",
	}

	if !reflect.DeepEqual(token, expectedToken) {
		t.Errorf("expected token %+v, got %+v", expectedToken, token)
	}
}

func TestClient_RequestUserAccessToken(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and content type
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected application/x-www-form-urlencoded content type")
		}

		// Parse and verify form data
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
		}

		// Verify OAuth parameters
		expectedParams := map[string]string{
			"client_id":     "test-client-id",
			"client_secret": "test-client-secret",
			"code":          "test-auth-code",
			"grant_type":    "authorization_code",
			"redirect_uri":  "http://localhost:3000/callback",
		}

		for key, expected := range expectedParams {
			if actual := r.FormValue(key); actual != expected {
				t.Errorf("expected %s '%s', got '%s'", key, expected, actual)
			}
		}

		// Return mock token response
		response := UserAccessToken{
			AccessToken:  "test-user-access-token",
			ExpiresIn:    14400,
			RefreshToken: "test-refresh-token",
			TokenType:    "bearer",
			Scope:        []string{"user:read:email", "channel:read:subscriptions"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server
	client := &Client{
		client:  http.DefaultClient,
		baseURL: server.URL,
		config: &config.TwitchConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	}

	token, err := client.RequestUserAccessToken("test-auth-code", "http://localhost:3000/callback")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedToken := &UserAccessToken{
		AccessToken:  "test-user-access-token",
		ExpiresIn:    14400,
		RefreshToken: "test-refresh-token",
		TokenType:    "bearer",
		Scope:        []string{"user:read:email", "channel:read:subscriptions"},
	}

	if !reflect.DeepEqual(token, expectedToken) {
		t.Errorf("expected token %+v, got %+v", expectedToken, token)
	}
}

func TestAppAccessToken_IsExpired(t *testing.T) {
	test := []struct {
		Name     string
		Token    *AppAccessToken
		IssuedAt time.Time
		Expected bool
	}{
		{
			Name: "Not expired",
			Token: &AppAccessToken{
				ExpiresIn: 3600,
			},
			IssuedAt: time.Now().Add(-time.Minute * 30), // 30 minutes ago
			Expected: false,
		},
		{
			Name: "Expired",
			Token: &AppAccessToken{
				ExpiresIn: 3600,
			},
			IssuedAt: time.Now().Add(-2 * time.Hour), // 2 hours ago
			Expected: true,
		},
		{
			Name: "About to expire (within margin)",
			Token: &AppAccessToken{
				ExpiresIn: 3600,
			},
			IssuedAt: time.Now().Add(-time.Hour).Add(4 * time.Minute), // 56 minutes ago
			Expected: true,
		},
		{
			Name: "About to expire (outside margin)",
			Token: &AppAccessToken{
				ExpiresIn: 3600,
			},
			IssuedAt: time.Now().Add(-time.Hour).Add(6 * time.Minute), // 54 minutes ago
			Expected: false,
		},
	}

	for _, tc := range test {
		t.Run(tc.Name, func(t *testing.T) {
			expired := tc.Token.IsExpired(&tc.IssuedAt)
			if expired != tc.Expected {
				t.Errorf("expected expired to be %v, got %v", tc.Expected, expired)
			}
		})
	}
}
