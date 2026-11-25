// Package bot answers chat commands. Chat messages arrive through EventSub
// (package eventsub) and replies go out through a Sender (Helix).
package bot

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/SoWieMarkus/wheelgpt/internal/channel"
	"github.com/SoWieMarkus/wheelgpt/internal/guess"
	"github.com/SoWieMarkus/wheelgpt/internal/ingame"
	"github.com/SoWieMarkus/wheelgpt/internal/leaderboard"
)

// Access is the role of a chatter. Higher values include lower ones.
type Access int

const (
	Everyone Access = iota
	Subscriber
	VIP
	Moderator
	Broadcaster
)

func (a Access) String() string {
	return [...]string{"everyone", "subscriber", "vip", "moderator", "broadcaster"}[a]
}

// AccessFromBadges derives the role from EventSub chat badge set ids.
func AccessFromBadges(setIDs []string) Access {
	best := Everyone
	for _, id := range setIDs {
		var a Access
		switch id {
		case "broadcaster":
			a = Broadcaster
		case "moderator", "lead_moderator":
			a = Moderator
		case "vip":
			a = VIP
		case "subscriber", "founder":
			a = Subscriber
		}
		best = max(best, a)
	}
	return best
}

// Message is a chat message.
type Message struct {
	BroadcasterID string
	ChatterID     string
	ChatterName   string
	Text          string
	Access        Access
}

// Sender posts a message into a channel's chat.
type Sender interface {
	Send(ctx context.Context, broadcasterID, text string) error
}

// Recorder counts commands (metrics).
type Recorder interface {
	Command(name string)
}

// MaxMessageLength is Twitch's chat message limit.
const MaxMessageLength = 500

// Deps are the services commands use.
type Deps struct {
	Channels    *channel.Service
	Guesses     *guess.Service
	Ingame      *ingame.Service
	Leaderboard *leaderboard.Service
	Sender      Sender
	Recorder    Recorder
	// PublicURL is the origin of the website, used in chat links.
	PublicURL string
}

// Bot dispatches commands.
type Bot struct {
	Deps
	commands []*Command

	mu       sync.Mutex
	lastUsed map[string]time.Time
	pending  sync.WaitGroup
	stopped  chan struct{}
	stopOnce sync.Once
}

func New(d Deps) *Bot {
	d.PublicURL = strings.TrimRight(d.PublicURL, "/")
	return &Bot{Deps: d, commands: defaultCommands(), lastUsed: map[string]time.Time{}, stopped: make(chan struct{})}
}

// Commands returns the registered commands.
func (b *Bot) Commands() []*Command { return b.commands }

// Parse splits a chat message into a command key and arguments. ok is false
// if the message is not a command. Arguments without visible ASCII
// characters are dropped (some chat clients append invisible characters).
func Parse(text string) (key string, args []string, ok bool) {
	tokens := strings.Split(text, " ")
	if len(tokens) == 0 || !strings.HasPrefix(tokens[0], "!") {
		return "", nil, false
	}
	key = strings.ToLower(strings.TrimPrefix(tokens[0], "!"))
	if key == "" {
		return "", nil, false
	}
	for _, t := range tokens[1:] {
		visible := strings.Map(func(r rune) rune {
			if r >= 0x20 && r <= 0x7e {
				return r
			}
			return -1
		}, t)
		if strings.TrimSpace(visible) != "" {
			args = append(args, t)
		}
	}
	return key, args, true
}

func (b *Bot) find(key string) *Command {
	for _, c := range b.commands {
		if c.Matches(key) {
			return c
		}
	}
	return nil
}

// Handle runs the command in a chat message, if any.
func (b *Bot) Handle(ctx context.Context, c channel.Channel, msg Message) error {
	key, args, ok := Parse(msg.Text)
	if !ok {
		return nil
	}
	cmd := b.find(key)
	if cmd == nil || msg.Access < cmd.Access || !b.takeCooldown(c.ID, cmd) {
		return nil
	}
	b.Recorder.Command(cmd.Name)
	reply, err := cmd.Run(ctx, b, Invocation{Channel: c, User: msg, Args: args})
	if err != nil || reply == "" {
		return err
	}
	return b.Sender.Send(ctx, c.ID, Truncate(reply))
}

func (b *Bot) takeCooldown(channelID string, cmd *Command) bool {
	if cmd.Cooldown <= 0 {
		return true
	}
	key := channelID + "\x00" + cmd.Name
	now := b.Guesses.Now()
	b.mu.Lock()
	defer b.mu.Unlock()
	if last, ok := b.lastUsed[key]; ok && now.Sub(last) < cmd.Cooldown {
		return false
	}
	b.lastUsed[key] = now
	return true
}

// Announce sends text to a channel after delay. Pending announcements are
// sent right away when the bot stops, so no result gets lost on a restart.
func (b *Bot) Announce(channelID, text string, delay time.Duration) {
	b.pending.Go(func() {
		timer := time.NewTimer(delay)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-b.stopped:
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := b.Sender.Send(ctx, channelID, Truncate(text)); err != nil {
			slog.ErrorContext(ctx, "announce failed", "channel", channelID, "error", err)
		}
	})
}

// Stop sends pending announcements right away and waits for them.
func (b *Bot) Stop() {
	b.stopOnce.Do(func() { close(b.stopped) })
	b.pending.Wait()
}

// Truncate shortens text to the chat message limit.
func Truncate(text string) string {
	if utf8.RuneCountInString(text) <= MaxMessageLength {
		return text
	}
	r := []rune(text)
	return string(r[:MaxMessageLength-1]) + "…"
}
