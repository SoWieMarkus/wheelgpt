// Package metrics exposes Prometheus metrics for HTTP traffic, the database
// pool and bot events like commands, personal bests and EventSub messages.
package metrics

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "wheelgpt"

type Metrics struct {
	registry *prometheus.Registry

	httpRequests *prometheus.CounterVec
	httpDuration *prometheus.HistogramVec
	httpInFlight prometheus.Gauge

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
	reg := prometheus.NewRegistry()
	m := &Metrics{
		registry: reg,
		httpRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace, Subsystem: "http", Name: "requests_total",
			Help: "HTTP requests by route, method and status code.",
		}, []string{"route", "method", "status"}),
		httpDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace, Subsystem: "http", Name: "request_duration_seconds",
			Help:    "HTTP request latency by route and method.",
			Buckets: prometheus.DefBuckets,
		}, []string{"route", "method"}),
		httpInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace, Subsystem: "http", Name: "requests_in_flight",
			Help: "HTTP requests currently being served.",
		}),
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

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace: namespace, Name: "build_info",
			Help:        "Build information, always 1.",
			ConstLabels: prometheus.Labels{"version": version},
		}, func() float64 { return 1 }),
		m.httpRequests, m.httpDuration, m.httpInFlight,
		m.logins, m.commands, m.personalBests, m.guesses,
		m.eventsub, m.chatSent, m.twitchFailures,
	)
	if db != nil {
		reg.MustRegister(collectors.NewDBStatsCollector(db, namespace))
	}
	return m
}

// Handler serves the metrics in the Prometheus exposition format.
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registry})
}

// Middleware records request metrics. It labels by chi route pattern instead
// of raw path to keep cardinality bounded.
func (m *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.httpInFlight.Inc()
		defer m.httpInFlight.Dec()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		route := "unmatched"
		if rc := chi.RouteContext(r.Context()); rc != nil && rc.RoutePattern() != "" {
			route = rc.RoutePattern()
		}
		m.httpRequests.WithLabelValues(route, r.Method, strconv.Itoa(ww.Status())).Inc()
		m.httpDuration.WithLabelValues(route, r.Method).Observe(time.Since(start).Seconds())
	})
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
