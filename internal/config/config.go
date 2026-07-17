package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	PORT        string `env:"PORT" envDefault:"8080"`

	TwitchClientID     string `env:"TWITCH_CLIENT_ID,required"`
	TwitchClientSecret string `env:"TWITCH_CLIENT_SECRET,required"`
	TwitchRedirectURL  string `env:"TWITCH_REDIRECT_URL,required"`

	BotUsername   string `env:"BOT_USERNAME,required"`
	BotOAuthToken string `env:"BOT_OAUTH_TOKEN,required"`

	EventSubSecret      string `env:"EVENTSUB_SECRET,required"`
	EventSubCallbackURL string `env:"EVENTSUB_CALLBACK_URL,required"`
	UpdateWebHooks      bool   `env:"UPDATE_WEBHOOKS" envDefault:"false"`

	JWTSecretWeb     string `env:"JWT_SECRET_WEB,required"`
	JWTSecretChannel string `env:"JWT_SECRET_CHANNEL,required"`

	CommandPrefix string `env:"COMMAND_PREFIX" envDefault:"!"`
}

func Load() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}
