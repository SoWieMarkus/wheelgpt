package bot

import (
	"context"
	"strings"
	"sync"
	"time"
)

// Command is the interface every chat command implements.
type Command interface {
	Names() []string // primary name is Names()[0]
	Level() AccessLevel
	Execute(ctx *CmdContext) (string, error)
}

// CmdContext carries everything a command needs without importing the whole bot.
type CmdContext struct {
	ChannelID string
	User      User
	Args      []string
	Say       func(msg string) // callback to send a message to chat
}

type commandEntry struct {
	cmd      Command
	cooldown time.Duration
	lastUsed time.Time
	mu       sync.Mutex
}

func (e *commandEntry) ready() bool {
	if e.cooldown == 0 {
		return true
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	return time.Since(e.lastUsed) >= e.cooldown
}

func (e *commandEntry) touch() {
	e.mu.Lock()
	e.lastUsed = time.Now()
	e.mu.Unlock()
}

// PBHandler is called by NotifyNewPB after the guess delay. Injected by main
// to break the import cycle between bot and bot/commands.
type PBHandler func(ctx context.Context, store interface{}, channelID string, t TMTime, say func(string))

// Registry maps command names → entries.
type Registry struct {
	entries   map[string]*commandEntry
	handlePB  func(ctx context.Context, store interface{}, channelID string, t TMTime, say func(string))
}

func NewRegistry() *Registry {
	return &Registry{entries: make(map[string]*commandEntry)}
}

func (r *Registry) Register(cmd Command, cooldown time.Duration) {
	e := &commandEntry{cmd: cmd, cooldown: cooldown}
	for _, name := range cmd.Names() {
		r.entries[strings.ToLower(name)] = e
	}
}

// Dispatch parses a raw message, finds the command, checks access and cooldown,
// then calls Execute. Returns "" if the message is not a command or should be suppressed.
func (r *Registry) Dispatch(channelID string, user User, message string, say func(string)) {
	if !strings.HasPrefix(message, "!") {
		return
	}
	parts := strings.Fields(strings.TrimPrefix(message, "!"))
	if len(parts) == 0 {
		return
	}
	key := strings.ToLower(parts[0])
	args := parts[1:]

	e, ok := r.entries[key]
	if !ok {
		return
	}
	if user.Level < e.cmd.Level() {
		return
	}
	if !e.ready() {
		return
	}
	e.touch()

	ctx := &CmdContext{
		ChannelID: channelID,
		User:      user,
		Args:      args,
		Say:       say,
	}
	result, err := e.cmd.Execute(ctx)
	if err != nil || result == "" {
		return
	}
	say(result)
}
