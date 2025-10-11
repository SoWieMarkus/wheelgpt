package config

import (
	"reflect"
	"testing"
)

func Test_LoadDatabaseConfig(t *testing.T) {
	tests := []struct {
		Name        string
		EnvVars     map[string]string
		Expected    *DatabaseConfig
		ExpectError bool
	}{
		{
			Name: "Set all variables",
			EnvVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "6543",
				"DB_USER":     "envuser",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "envdb",
				"DB_SSLMODE":  "require",
			},
			Expected: &DatabaseConfig{
				HostName:     "db.example.com",
				Port:         6543,
				User:         "envuser",
				Password:     "secret",
				DatabaseName: "envdb",
				SSLMode:      "require",
			},
			ExpectError: false,
		},
		{
			Name: "Required DB_USER missing",
			EnvVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "6543",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "envdb",
				"DB_SSLMODE":  "require",
			},
			Expected:    &DatabaseConfig{},
			ExpectError: true,
		},
		{
			Name: "Required DB_HOST missing",
			EnvVars: map[string]string{
				"DB_PORT":     "6543",
				"DB_USER":     "envuser",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "envdb",
				"DB_SSLMODE":  "require",
			},
			Expected:    &DatabaseConfig{},
			ExpectError: true,
		},
		{
			Name: "Required DB_PORT missing",
			EnvVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_USER":     "envuser",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "envdb",
				"DB_SSLMODE":  "require",
			},
			Expected:    &DatabaseConfig{},
			ExpectError: true,
		},
		{
			Name: "Required DB_PASSWORD missing",
			EnvVars: map[string]string{
				"DB_HOST":    "db.example.com",
				"DB_PORT":    "6543",
				"DB_USER":    "envuser",
				"DB_NAME":    "envdb",
				"DB_SSLMODE": "require",
			},
			Expected:    &DatabaseConfig{},
			ExpectError: true,
		},
		{
			Name: "Required DB_NAME missing",
			EnvVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "6543",
				"DB_USER":     "envuser",
				"DB_PASSWORD": "secret",
				"DB_SSLMODE":  "require",
			},
			Expected:    &DatabaseConfig{},
			ExpectError: true,
		},
		{
			Name: "Optional DB_SSLMODE missing",
			EnvVars: map[string]string{
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "6543",
				"DB_USER":     "envuser",
				"DB_PASSWORD": "secret",
				"DB_NAME":     "envdb",
			},
			Expected: &DatabaseConfig{
				HostName:     "db.example.com",
				Port:         6543,
				User:         "envuser",
				Password:     "secret",
				DatabaseName: "envdb",
				// Default to disable
				SSLMode: "disable",
			},
			ExpectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.EnvVars {
				t.Setenv(key, value)
			}

			config, err := LoadDatabaseConfig()
			if tt.ExpectError != (err != nil) {
				t.Errorf("Expected error: %v, but got error: %v", tt.ExpectError, err)
			}

			if tt.ExpectError {
				return
			}

			if !reflect.DeepEqual(config, tt.Expected) {
				t.Errorf("Expected config %+v but got %+v", tt.Expected, config)
			}
		})
	}
}

func Test_LoadTwitchConfig(t *testing.T) {
	tests := []struct {
		Name        string
		EnvVars     map[string]string
		Expected    *TwitchConfig
		ExpectError bool
	}{
		{
			Name: "Set all variables",
			EnvVars: map[string]string{
				"TWITCH_CLIENT_ID":     "myclientid",
				"TWITCH_CLIENT_SECRET": "secret",
				"TWITCH_BOT_USERNAME":  "mybot",
				"TWITCH_OAUTH_TOKEN":   "oauth:token",
			},
			Expected: &TwitchConfig{
				ClientID:     "myclientid",
				ClientSecret: "secret",
				BotUsername:  "mybot",
				OAuthToken:   "oauth:token",
			},
			ExpectError: false,
		},
		{
			Name: "Required TWITCH_CLIENT_ID missing",
			EnvVars: map[string]string{
				"TWITCH_CLIENT_SECRET": "secret",
				"TWITCH_BOT_USERNAME":  "mybot",
				"TWITCH_OAUTH_TOKEN":   "oauth:token",
			},
			Expected:    &TwitchConfig{},
			ExpectError: true,
		},
		{
			Name: "Required TWITCH_CLIENT_SECRET missing",
			EnvVars: map[string]string{
				"TWITCH_CLIENT_ID":    "myclientid",
				"TWITCH_BOT_USERNAME": "mybot",
				"TWITCH_OAUTH_TOKEN":  "oauth:token",
			},
			Expected:    &TwitchConfig{},
			ExpectError: true,
		},
		{
			Name: "Required TWITCH_BOT_USERNAME missing",
			EnvVars: map[string]string{
				"TWITCH_CLIENT_ID":     "myclientid",
				"TWITCH_CLIENT_SECRET": "secret",
				"TWITCH_OAUTH_TOKEN":   "oauth:token",
			},
			Expected:    &TwitchConfig{},
			ExpectError: true,
		},
		{
			Name: "Required TWITCH_OAUTH_TOKEN missing",
			EnvVars: map[string]string{
				"TWITCH_CLIENT_ID":     "myclientid",
				"TWITCH_CLIENT_SECRET": "secret",
				"TWITCH_BOT_USERNAME":  "mybot",
			},
			Expected:    &TwitchConfig{},
			ExpectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.EnvVars {
				t.Setenv(key, value)
			}

			config, err := LoadTwitchConfig()
			if tt.ExpectError != (err != nil) {
				t.Errorf("Expected error: %v, but got error: %v", tt.ExpectError, err)
			}

			if tt.ExpectError {
				return
			}

			if !reflect.DeepEqual(config, tt.Expected) {
				t.Errorf("Expected config %+v but got %+v", tt.Expected, config)
			}
		})
	}
}
