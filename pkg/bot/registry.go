package bot

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/SoWieMarkus/wheelgpt/internal/config"
)

// command defines the interface that all commands must implement.
type command interface {
	Names() []string
	Execute(args CommandContext) (string, error)
}

// commandEntry represents a registered command with its associated cooldown and access level.
type commandEntry struct {
	command     command
	cooldown    time.Duration
	accessLevel ViewerAccessLevel
	lastUsed    time.Time
	mu          sync.Mutex
}

// isOnCooldown checks if the command entry is currently on cooldown.
func (ce *commandEntry) isOnCooldown() bool {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	return time.Since(ce.lastUsed) < ce.cooldown
}

// updateLastUsed updates the last used time of the command entry to the current time.
func (ce *commandEntry) updateLastUsed() {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	ce.lastUsed = time.Now()
}

// CommandContext provides context for command execution, including the channel ID, arguments, and a function to send messages.
type CommandContext struct {
	ChannelID string
	Args      []string
	Say       func(message string)
}

// Registry is a command registry that allows registering and executing commands with cooldowns.
type Registry struct {
	commands      map[string]*commandEntry
	commandPrefix string
}

// NewRegistry creates a new command registry.
func NewRegistry(cfg config.Config) *Registry {
	return &Registry{
		commands:      make(map[string]*commandEntry),
		commandPrefix: cfg.CommandPrefix,
	}
}

// Register adds a command to the registry with the specified cooldown duration.
func (r *Registry) Register(cmd command, cooldown time.Duration) {
	entry := &commandEntry{
		command:  cmd,
		cooldown: cooldown,
	}
	for _, name := range cmd.Names() {
		r.commands[name] = entry
	}
}

func (r *Registry) Dispatch(channelID string, viewer Viewer, message string, say func(message string)) {
	// Check if the message starts with the command prefix
	// e.g. "!test"
	if !strings.HasPrefix(message, r.commandPrefix) {
		return
	}

	// The first after the command prefix is expected to be the command name.
	// Everything after the command name is considered as arguments
	parts := strings.Fields(strings.TrimPrefix(message, r.commandPrefix))
	if len(parts) == 0 {
		return
	}

	key := parts[0]
	args := parts[1:]

	entry, ok := r.commands[key]
	if !ok {
		return
	}

	// TODO: Add access level checks

	if entry.isOnCooldown() {
		return
	}

	entry.updateLastUsed()

	// Check if the viewer has the required access level to execute the command
	if entry.accessLevel > viewer.AccessLevel {
		return
	}

	// Execute the command
	response, err := entry.command.Execute(CommandContext{
		ChannelID: channelID,
		Args:      args,
		Say:       say,
	})
	if err != nil {
		log.Printf("Error executing command %s: %v", key, err)
		return
	}
	if response != "" {
		say(response)
	}
}
