package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func createMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))
}

func TestClient_NewClient(t *testing.T) {
	tests := []struct {
		Name           string
		BaseUrl        string
		DefaultHeaders *map[string]string
	}{
		{
			Name:           "No Default Headers",
			BaseUrl:        "https://api.example.com",
			DefaultHeaders: nil,
		},
		{
			Name:    "With Default Headers",
			BaseUrl: "https://api.example.com",
			DefaultHeaders: &map[string]string{
				"Authorization": "Bearer test-token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			client := NewClient(tt.BaseUrl, tt.DefaultHeaders)

			if client == nil {
				t.Fatal("expected client to be created, got nil")
			}

			if client.baseUrl != tt.BaseUrl {
				t.Errorf("expected baseUrl to be '%s', got '%s'", tt.BaseUrl, client.baseUrl)
			}

			if tt.DefaultHeaders == nil && client.defaultHeaders != nil {
				t.Error("expected defaultHeaders to be nil")
			}
			if tt.DefaultHeaders != nil && !reflect.DeepEqual(client.defaultHeaders, tt.DefaultHeaders) {
				t.Errorf("expected defaultHeaders to be set correctly")
			}

			if client.client == nil {
				t.Error("expected http.Client to be initialized")
			}
		})
	}
}

func TestClient_Get(t *testing.T) {

}

func TestClient_Post(t *testing.T) {

}

func TestClient_Put(t *testing.T) {
}

func TestClient_Delete(t *testing.T) {
}

func TestClient_setHeaders(t *testing.T) {
	defaultHeaders := map[string]string{
		"Authorization": "Bearer default-token",
		"Content-Type":  "application/json",
	}
	client := NewClient("https://api.example.com", &defaultHeaders)

	req, err := http.NewRequest("GET", "https://api.example.com/resource", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test with no additional headers
	client.setHeaders(req, nil)
	for key, expected := range defaultHeaders {
		if actual := req.Header.Get(key); actual != expected {
			t.Errorf("expected header %s to be '%s', got '%s'", key, expected, actual)
		}
	}

	// Test with additional headers
	additionalHeaders := map[string]string{
		"Authorization":   "Bearer new-token", // This should override the default
		"X-Custom-Header": "custom-value",
	}
	client.setHeaders(req, &additionalHeaders)

	expectedHeaders := map[string]string{
		"Authorization":   "Bearer new-token", // Overridden
		"Content-Type":    "application/json", // From default
		"X-Custom-Header": "custom-value",     // New header
	}

	for key, expected := range expectedHeaders {
		if actual := req.Header.Get(key); actual != expected {
			t.Errorf("expected header %s to be '%s', got '%s'", key, expected, actual)
		}
	}
}

func TestClient_buildRequest(t *testing.T) {
	client := NewClient("https://api.example.com", nil)
	reqDetails := &HttpRequest{
		Endpoint: "/test-endpoint",
		Headers: &map[string]string{
			"X-Test-Header": "test-value",
		},
		Params: &url.Values{
			"param1": []string{"value1"},
			"param2": []string{"value2"},
		},
	}

	req, err := client.buildRequest("GET", reqDetails, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("expected method to be 'GET', got '%s'", req.Method)
	}

	expectedUrl := "https://api.example.com/test-endpoint?param1=value1&param2=value2"
	if req.URL.String() != expectedUrl {
		t.Errorf("expected URL to be '%s', got '%s'", expectedUrl, req.URL.String())
	}

	if req.Header.Get("X-Test-Header") != "test-value" {
		t.Errorf("expected header 'X-Test-Header' to be 'test-value', got '%s'", req.Header.Get("X-Test-Header"))
	}
}

func TestClient_buildRequest_WithBody(t *testing.T) {
	client := NewClient("https://api.example.com", nil)

	headers := map[string]string{
		"Content-Type": "application/json",
	}
	reqDetails := &HttpRequest{
		Endpoint: "/test-endpoint",
		Headers:  &headers,
	}

	body := map[string]string{
		"key": "value",
	}
	bodyBytes, err := json.Marshal(body)
	bodyReader := bytes.NewReader(bodyBytes)

	req, err := client.buildRequest("POST", reqDetails, bodyReader)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("expected method to be 'POST', got '%s'", req.Method)
	}

	expectedUrl := "https://api.example.com/test-endpoint"
	if req.URL.String() != expectedUrl {
		t.Errorf("expected URL to be '%s', got '%s'", expectedUrl, req.URL.String())
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected header 'Content-Type' to be 'application/json', got '%s'", req.Header.Get("Content-Type"))
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	if buf.String() != `{"key":"value"}` {
		t.Errorf("expected body to be '{\"key\":\"value\"}', got '%s'", buf.String())
	}
}

func TestClient_setQueryParams(t *testing.T) {
	client := NewClient("https://api.example.com", nil)
	req, err := http.NewRequest("GET", "https://api.example.com/resource", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	params := url.Values{
		"param1": []string{"value1", "value1-2"},
		"param2": []string{"value2"},
	}

	client.setQueryParams(req, &params)

	expectedURL := "https://api.example.com/resource?param1=value1&param1=value1-2&param2=value2"
	if req.URL.String() != expectedURL {
		t.Errorf("expected URL to be '%s', got '%s'", expectedURL, req.URL.String())
	}
}
