package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Addr        string `env:"ADDR" envDefault:":8080"`

	TwitchClientID     string `env:"TWITCH_CLIENT_ID,required"`
	TwitchClientSecret string `env:"TWITCH_CLIENT_SECRET,required"`
	TwitchRedirectURL  string `env:"TWITCH_REDIRECT_URL,required"`

	BotUsername   string `env:"BOT_USERNAME,required"`
	BotOAuthToken string `env:"BOT_OAUTH_TOKEN,required"`

	EventSubSecret     string `env:"TWITCH_EVENTSUB_SECRET,required"`
	EventSubCallbackURL string `env:"TWITCH_STREAM_WEB_HOOK_URL,required"`
	UpdateWebHooks     bool   `env:"UPDATE_WEB_HOOKS" envDefault:"false"`

	JWTSecretWeb     string `env:"JWT_SECRET_WEB,required"`
	JWTSecretChannel string `env:"JWT_SECRET_CHANNEL,required"`
}

func Load() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("config: %v", err)
	}
	return cfg
}
