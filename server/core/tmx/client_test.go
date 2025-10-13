package tmx

import "testing"

func TestClient_NewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.client == nil {
		t.Fatal("Expected non-nil http.Client")
	}

	if client.client.GetBaseUrl() != "https://trackmania.exchange/api" {
		t.Fatalf("Expected base URL to be https://trackmania.exchange/api, got %s", client.client.GetBaseUrl())
	}
}
