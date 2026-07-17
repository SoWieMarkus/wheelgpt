package bot

import (
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/db"
)

// ChannelConfig is the subset of db.Channel the bot needs at runtime.
type ChannelConfig struct {
	ID                   string
	Login                string
	DisplayName          string
	GuessDelayTime       float64
	BotActiveWhenOffline bool
}

// Channel holds a per-channel command registry.
type Channel struct {
	Config   ChannelConfig
	registry *Registry
}

func NewChannel(cfg ChannelConfig, store *db.Store) *Channel {
	reg := NewRegistry()

	// Import here to avoid circular dep — commands live in a sub-package.
	// The caller wires them in via RegisterCommands.
	return &Channel{Config: cfg, registry: reg}
}

func (c *Channel) Registry() *Registry { return c.registry }

func (c *Channel) GuessDelay() time.Duration {
	return time.Duration(c.Config.GuessDelayTime * float64(time.Second))
}
