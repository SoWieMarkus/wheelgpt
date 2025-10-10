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

	config, err := config.LoadDatabaseConfig()
	if err != nil {
		panic("Failed to load database config: " + err.Error())
	}

	ctx := context.Background()

	postgres := db.NewPostgresDB(ctx, config)
	defer postgres.Close()

	bot := wheelgpt.NewInstance(postgres)
	bot.Initialize()

	// TODO start the bot

}
