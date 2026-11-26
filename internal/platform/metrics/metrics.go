// Package metrics exposes the business metrics of wheelgpt like logins, bot
// commands, personal bests and EventSub messages. HTTP, database pool and
// runtime metrics come from the embedded gokit registry.
package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"

	kitmetrics "github.com/SoWieMarkus/gokit/metrics"
)

const namespace = "wheelgpt"

type Metrics struct {
	*kitmetrics.Registry

	logins         *prometheus.CounterVec
	commands       *prometheus.CounterVec
	personalBests  *prometheus.CounterVec
	guesses        prometheus.Counter
	eventsub       *prometheus.CounterVec
	chatSent       *prometheus.CounterVec
	twitchFailures *prometheus.CounterVec
}

// New registers all collectors on a fresh registry. db may be nil (tests).
func New(version string, db *sql.DB) *Metrics {
	m := &Metrics{
		Registry: kitmetrics.New(namespace, version, db),
		logins: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Name: "logins_total",
			Help: "Twitch logins by outcome (new, returning, failure).",
		}, []string{"outcome"}),
		commands: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Subsystem: "bot", Name: "commands_total",
			Help: "Executed chat commands.",
		}, []string{"command"}),
		personalBests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Name: "personal_bests_total",
			Help: "Personal bests reported by the plugin, by outcome (evaluated, offline).",
		}, []string{"outcome"}),
		guesses: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace, Name: "guesses_total",
			Help: "Guesses placed in chat.",
		}),
		eventsub: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Subsystem: "eventsub", Name: "messages_total",
			Help: "Verified EventSub messages by message and subscription type.",
		}, []string{"message_type", "subscription_type"}),
		chatSent: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Subsystem: "bot", Name: "chat_messages_total",
			Help: "Chat messages the bot tried to send, by outcome.",
		}, []string{"outcome"}),
		twitchFailures: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Subsystem: "twitch", Name: "failed_requests_total",
			Help: "Failed Twitch API requests by endpoint.",
		}, []string{"endpoint"}),
	}

	m.MustRegister(
		m.logins, m.commands, m.personalBests, m.guesses,
		m.eventsub, m.chatSent, m.twitchFailures,
	)
	return m
}

func (m *Metrics) Login(outcome string)          { m.logins.WithLabelValues(outcome).Inc() }
func (m *Metrics) Command(name string)           { m.commands.WithLabelValues(name).Inc() }
func (m *Metrics) PersonalBest(outcome string)   { m.personalBests.WithLabelValues(outcome).Inc() }
func (m *Metrics) Guess()                        { m.guesses.Inc() }
func (m *Metrics) ChatMessage(outcome string)    { m.chatSent.WithLabelValues(outcome).Inc() }
func (m *Metrics) TwitchFailure(endpoint string) { m.twitchFailures.WithLabelValues(endpoint).Inc() }

func (m *Metrics) EventSub(messageType, subscriptionType string) {
	m.eventsub.WithLabelValues(messageType, subscriptionType).Inc()
}
