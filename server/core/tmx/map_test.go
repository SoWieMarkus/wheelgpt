package tmx

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpClient "github.com/SoWieMarkus/wheelgpt/core/http"
)

func TestClient_GetMapInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		response := "{}"

		if r.URL.Path == "/map/test-map-id" {
			response = `{"TrackID": "test-track-id"}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	client := &Client{
		client: httpClient.NewClient(server.URL, nil),
	}

	tests := []struct {
		Name            string
		MapID           string
		ExpectError     bool
		ExpectedTrackID string
	}{
		{
			Name:            "Valid Map ID",
			MapID:           "test-map-id",
			ExpectError:     false,
			ExpectedTrackID: "test-track-id",
		},
		{
			Name:            "Invalid Map ID",
			MapID:           "invalid-map-id",
			ExpectError:     false,
			ExpectedTrackID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mapInfo, err := client.GetMapInfo(tt.MapID)
			if tt.ExpectError {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if mapInfo.TrackID != tt.ExpectedTrackID {
				t.Errorf("expected TrackID '%s', got '%s'", tt.ExpectedTrackID, mapInfo.TrackID)
			}
		})
	}
}
