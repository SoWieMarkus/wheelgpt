package app

// This file loads the service configuration from environment variables.

import (
	"errors"
	"log/slog"
	"net/url"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPAddr    string     `env:"HTTP_ADDR"    envDefault:":8080"`
	MetricsAddr string     `env:"METRICS_ADDR" envDefault:":9091"`
	LogLevel    slog.Level `env:"LOG_LEVEL"    envDefault:"info"`
	LogFormat   string     `env:"LOG_FORMAT"   envDefault:"json"`

	// DataDir holds the SQLite database. It must be writable and should be on
	// a persistent volume.
	DataDir string `env:"DATA_DIR" envDefault:"data"`
	// DatabasePath overrides DataDir/wheelgpt.db, for example to point at the
	// file of the old backend.
	DatabaseFile string `env:"DATABASE_PATH"`

	// PublicURL is where users reach the website. Its scheme decides whether
	// cookies are marked Secure, its origin is trusted for cross origin
	// requests and OAuth redirect URIs and chat links are built from it.
	PublicURL *url.URL `env:"PUBLIC_URL" envDefault:"http://localhost:5173"`

	SessionTTL time.Duration `env:"SESSION_TTL" envDefault:"720h"`
	// LoginRateLimit is the number of login requests per IP and minute.
	LoginRateLimit int `env:"LOGIN_RATE_LIMIT" envDefault:"20"`
	// PluginRateLimit is the number of plugin requests per channel and minute.
	PluginRateLimit int `env:"PLUGIN_RATE_LIMIT" envDefault:"120"`
	// TrustProxyHeaders uses X-Real-IP as client IP. Only enable it behind a
	// reverse proxy that sets the header (the website's nginx does).
	TrustProxyHeaders bool `env:"TRUST_PROXY_HEADERS" envDefault:"false"`

	// Twitch application (https://dev.twitch.tv/console/apps). Without it
	// the bot runs offline: chat replies are only logged.
	TwitchClientID     string `env:"TWITCH_CLIENT_ID"`
	TwitchClientSecret string `env:"TWITCH_CLIENT_SECRET"`
	// EventSubSecret signs the EventSub webhooks (10 to 100 characters).
	EventSubSecret string `env:"TWITCH_EVENTSUB_SECRET"`
	// EventSubSync creates and removes EventSub subscriptions. Keep it off
	// in development, the subscriptions belong to the production app.
	EventSubSync bool `env:"EVENTSUB_SYNC" envDefault:"false"`
	// BotUsername is the login of the bot account.
	BotUsername string `env:"BOT_USERNAME" envDefault:"wheelgpt"`

	// LegacyJWTSecret verifies plugin tokens of the old backend. Remove it
	// once every streamer created a new token.
	LegacyJWTSecret string `env:"JWT_SECRET_CHANNEL"`

	// DevMode is for local testing only and needs a localhost PUBLIC_URL:
	// /api/v1/auth/dev?login=name logs in without Twitch, and the dashboard
	// can pretend the stream is live so the plugin can be tested.
	DevMode bool `env:"DEV_MODE" envDefault:"false"`
}

// loadConfig parses the environment and validates the result.
func loadConfig() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	return cfg, cfg.validate()
}

func (c Config) validate() error {
	var errs []error
	if c.PublicURL.Scheme != "http" && c.PublicURL.Scheme != "https" {
		errs = append(errs, errors.New("PUBLIC_URL must be an http or https URL"))
	}
	if c.SessionTTL < time.Hour {
		errs = append(errs, errors.New("SESSION_TTL must be at least one hour"))
	}
	if c.LogFormat != "json" && c.LogFormat != "text" {
		errs = append(errs, errors.New("LOG_FORMAT must be json or text"))
	}
	if (c.TwitchClientID == "") != (c.TwitchClientSecret == "") {
		errs = append(errs, errors.New("set both TWITCH_CLIENT_ID and TWITCH_CLIENT_SECRET, or neither"))
	}
	if n := len(c.EventSubSecret); c.TwitchEnabled() && (n < 10 || n > 100) {
		errs = append(errs, errors.New("TWITCH_EVENTSUB_SECRET must be 10 to 100 characters"))
	}
	if c.EventSubSync && (!c.TwitchEnabled() || !c.SecureCookies()) {
		errs = append(errs, errors.New("EVENTSUB_SYNC needs Twitch credentials and an https PUBLIC_URL"))
	}
	if c.DevMode && c.PublicURL.Hostname() != "localhost" && c.PublicURL.Hostname() != "127.0.0.1" {
		errs = append(errs, errors.New("DEV_MODE is only allowed with a localhost PUBLIC_URL"))
	}
	return errors.Join(errs...)
}

// TwitchEnabled reports whether Twitch credentials are configured.
func (c Config) TwitchEnabled() bool { return c.TwitchClientID != "" }

// SecureCookies reports whether cookies must only be sent over HTTPS.
func (c Config) SecureCookies() bool { return c.PublicURL.Scheme == "https" }

// PublicOrigin is scheme://host[:port] of the website.
func (c Config) PublicOrigin() string { return c.PublicURL.Scheme + "://" + c.PublicURL.Host }

func (c Config) DatabasePath() string {
	if c.DatabaseFile != "" {
		return c.DatabaseFile
	}
	return filepath.Join(c.DataDir, "wheelgpt.db")
}
