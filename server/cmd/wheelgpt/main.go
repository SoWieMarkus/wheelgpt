package main

import (
	"context"
	"log/slog"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/services/wheelgpt"
)

func main() {
	slog.Info("WheelGPT started")

	databaseConfig, err := config.LoadDatabaseConfig()
	if err != nil {
		panic("Failed to load database config: " + err.Error())
	}

	ctx := context.Background()

	postgres := db.NewPostgresDB(ctx, databaseConfig)
	defer postgres.Close()

	twitchConfig, err := config.LoadTwitchConfig()
	if err != nil {
		panic("Failed to load Twitch config: " + err.Error())
	}

	bot := wheelgpt.NewInstance(postgres, twitchConfig)
	bot.Start()
}
