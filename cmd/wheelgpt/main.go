// Command wheelgpt runs the WheelGPT API: the Twitch chat bot, the
// Openplanet plugin API and the API of the website.
//
//	@title						WheelGPT API
//	@version					1.0
//	@description				API of WheelGPT, the Trackmania Twitch bot.
//	@description				Plugin endpoints need the plugin token as "Authorization: Bearer <token>". The me endpoints need the HttpOnly session cookie set by the Twitch login; state changing requests from other origins are rejected.
//	@description				Errors are returned as RFC 9457 problem details.
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	PluginToken
//	@in							header
//	@name						Authorization
//	@description				"Bearer <plugin token>" from the dashboard. Tokens of the old backend are accepted too.
//	@tag.name					plugin
//	@tag.description			Called by the Openplanet plugin
//	@tag.name					auth
//	@tag.description			Twitch login and logout (browser redirects)
//	@tag.name					me
//	@tag.description			The logged in streamer's channel
//	@tag.name					channels
//	@tag.description			Channels using WheelGPT, their map and room
//	@tag.name					leaderboards
//	@tag.description			Monthly and all-time rankings, rounds and guess history
//	@tag.name					bot
//	@tag.description			Chat commands
//	@tag.name					twitch
//	@tag.description			Twitch EventSub webhook
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/app"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	// "wheelgpt healthcheck" probes the running server. The image has no
	// shell or curl, so Docker's HEALTHCHECK uses the binary itself.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	if err := app.Run(version); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func healthcheck() int {
	port := 8080
	if _, p, err := net.SplitHostPort(os.Getenv("HTTP_ADDR")); err == nil {
		if n, err := strconv.Atoi(p); err == nil && n > 0 && n < 65536 {
			port = n
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//nolint:gosec // loopback probe; only the port comes from HTTP_ADDR
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/healthz", port), http.NoBody)
	if err != nil {
		return 1
	}
	res, err := http.DefaultClient.Do(req) //nolint:gosec // see above
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	_ = res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
