package config

import "github.com/kelseyhightower/envconfig"

type DatabaseConfig struct {
	HostName     string `envconfig:"DB_HOST" required:"true"`
	Port         int    `envconfig:"DB_PORT" required:"true"`
	User         string `envconfig:"DB_USER" required:"true"`
	Password     string `envconfig:"DB_PASSWORD" required:"true"`
	DatabaseName string `envconfig:"DB_NAME" required:"true"`
	SSLMode      string `envconfig:"DB_SSLMODE" default:"disable"`
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	var config DatabaseConfig
	err := envconfig.Process("", &config)
	return &config, err
}

type TwitchConfig struct {
	ClientID     string `envconfig:"TWITCH_CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"TWITCH_CLIENT_SECRET" required:"true"`
	BotUsername  string `envconfig:"TWITCH_BOT_USERNAME" required:"true"`
	OAuthToken   string `envconfig:"TWITCH_OAUTH_TOKEN" required:"true"`
}

func LoadTwitchConfig() (*TwitchConfig, error) {
	var config TwitchConfig
	err := envconfig.Process("", &config)
	return &config, err
}
