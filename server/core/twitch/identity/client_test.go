package identity

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/config"
)

func Test_NewClient(t *testing.T) {
	config := &config.TwitchConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("expected client to be created, got nil")
	}

	if client.baseURL != "https://id.twitch.tv/oauth2" {
		t.Errorf("expected baseURL to be 'https://id.twitch.tv/oauth2', got '%s'", client.baseURL)
	}

	if client.config != config {
		t.Error("expected config to be set correctly")
	}

	if client.client == nil {
		t.Error("expected http.Client to be initialized")
	}
}

// Test struct for response parsing
type TestResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func TestClient_Post_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("expected application/x-www-form-urlencoded content type")
		}

		// Parse and verify form data
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
		}

		if r.FormValue("test_param") != "test_value" {
			t.Errorf("expected test_param 'test_value', got '%s'", r.FormValue("test_param"))
		}

		// Return mock response
		response := TestResponse{
			Message: "success",
			Status:  200,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &Client{
		client:  http.DefaultClient,
		baseURL: server.URL,
		config:  &config.TwitchConfig{},
	}

	// Prepare test data
	data := url.Values{}
	data.Set("test_param", "test_value")

	// Execute test
	var result TestResponse
	err := client.post("/test", data, &result)

	// Verify results
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Message != "success" {
		t.Errorf("expected message 'success', got '%s'", result.Message)
	}

	if result.Status != 200 {
		t.Errorf("expected status 200, got %d", result.Status)
	}
}

func TestClient_Post_HTTPError(t *testing.T) {
	// Create client with invalid URL to simulate network error
	client := &Client{
		client:  http.DefaultClient,
		baseURL: "http://invalid-url-that-does-not-exist",
		config:  &config.TwitchConfig{},
	}

	data := url.Values{}
	data.Set("test", "value")

	var result TestResponse
	err := client.post("/test", data, &result)

	if err == nil {
		t.Error("expected error for invalid URL")
	}

	if !strings.Contains(err.Error(), "failed to execute request") {
		t.Errorf("expected 'failed to execute request' error, got %v", err)
	}
}

func TestClient_Post_NonOKStatus(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		response   string
	}{
		{"Bad Request", http.StatusBadRequest, `{"error":"invalid_request"}`},
		{"Unauthorized", http.StatusUnauthorized, `{"error":"unauthorized"}`},
		{"Internal Server Error", http.StatusInternalServerError, `{"error":"server_error"}`},
		{"Not Found", http.StatusNotFound, `{"error":"not_found"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
				w.Write([]byte(tc.response))
			}))
			defer server.Close()

			client := &Client{
				client:  http.DefaultClient,
				baseURL: server.URL,
				config:  &config.TwitchConfig{},
			}

			data := url.Values{}
			var result TestResponse
			err := client.post("/test", data, &result)

			if err == nil {
				t.Errorf("expected error for status %d", tc.statusCode)
			}

			expectedError := fmt.Sprintf("request failed with status %d", tc.statusCode)
			if !strings.Contains(err.Error(), expectedError) {
				t.Errorf("expected error to contain '%s', got %v", expectedError, err)
			}
		})
	}
}

func TestClient_Post_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json response"))
	}))
	defer server.Close()

	client := &Client{
		client:  http.DefaultClient,
		baseURL: server.URL,
		config:  &config.TwitchConfig{},
	}

	data := url.Values{}
	var result TestResponse
	err := client.post("/test", data, &result)

	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "failed to decode response") {
		t.Errorf("expected 'failed to decode response' error, got %v", err)
	}
}

func TestClient_Post_EmptyData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify empty form data
		err := r.ParseForm()
		if err != nil {
			t.Errorf("failed to parse form: %v", err)
		}

		if len(r.Form) != 0 {
			t.Errorf("expected empty form data, got %v", r.Form)
		}

		response := TestResponse{Message: "empty_data_ok", Status: 200}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := &Client{
		client:  http.DefaultClient,
		baseURL: server.URL,
		config:  &config.TwitchConfig{},
	}

	// Test with empty url.Values
	data := url.Values{}
	var result TestResponse
	err := client.post("/test", data, &result)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Message != "empty_data_ok" {
		t.Errorf("expected message 'empty_data_ok', got '%s'", result.Message)
	}
}

func TestClient_Post_RequestCreationError(t *testing.T) {
	client := &Client{
		client:  http.DefaultClient,
		baseURL: ":", // Invalid URL that will cause NewRequest to fail
		config:  &config.TwitchConfig{},
	}

	data := url.Values{}
	var result TestResponse
	err := client.post("/test", data, &result)

	if err == nil {
		t.Error("expected error for invalid URL")
	}

	if !strings.Contains(err.Error(), "failed to create request") {
		t.Errorf("expected 'failed to create request' error, got %v", err)
	}
}
