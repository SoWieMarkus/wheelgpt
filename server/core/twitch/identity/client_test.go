package identity

import (
	"testing"

	"github.com/SoWieMarkus/wheelgpt/core/config"
)

func Test_NewIdentityAPI(t *testing.T) {
	config := &config.TwitchConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}

	api := NewClient(config)

	if api == nil {
		t.Fatal("expected API to be created, got nil")
	}

	if api.client.GetBaseUrl() != "https://id.twitch.tv/oauth2" {
		t.Errorf("expected baseURL to be 'https://id.twitch.tv/oauth2', got '%s'", api.client.GetBaseUrl())
	}

	if api.config != config {
		t.Error("expected config to be set correctly")
	}

	if api.client == nil {
		t.Error("expected http.Client to be initialized")
	}
}
