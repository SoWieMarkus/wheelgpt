package main

import (
	"context"

	"github.com/SoWieMarkus/wheelgpt/core/config"
	"github.com/SoWieMarkus/wheelgpt/core/db"
	"github.com/SoWieMarkus/wheelgpt/core/db/models"
)

// List of all tables that are required by the application.
// Note: Removing a table from this list will NOT delete it from the database.
// It will only prevent it from being created if it does not exist.
var tables []db.Table = []db.Table{
	models.ChannelSettings{},
	models.Channel{},
	models.Guess{},
	models.TrackmaniaMap{},
}

func main() {
	// TODO load config
	var config config.DatabaseConfig

	ctx := context.Background()

	// Initialize database connection
	postgres := db.NewPostgresDB(ctx, config)
	defer postgres.Close()

	// Initialize new migrater
	migrater := db.NewMigrater(postgres)

	// Create all necessary tables if they do not exist
	migrater.CreateTableIfNotExists(tables...)

	// Run migrations
	migrater.Migrate()
}
