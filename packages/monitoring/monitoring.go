package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// ProviderType defines the monitoring provider type.
type ProviderType string

const (
	ProviderConsole    ProviderType = "console"
	ProviderPrometheus ProviderType = "prometheus"
	ProviderDatadog    ProviderType = "datadog"
	ProviderSentry     ProviderType = "sentry"
)

// Config holds the configuration parameters for Monitoring.
type Config struct {
	Provider ProviderType
	DSN      string // Connection string or DSN (e.g., Sentry DSN)
	Endpoint string // Metrics collector endpoint (e.g., Datadog agent URL)
}

// Monitor defines the interface for the monitoring system (Prometheus, Sentry, Datadog, etc.).
type Monitor interface {
	// Increment increments a counter metric by 1, typically used for request counts, error counts, etc.
	Increment(ctx context.Context, metric string, tags ...string)

	// Gauge updates the current value of a gauge metric, typically used for CPU, RAM, active users, etc.
	Gauge(ctx context.Context, metric string, value float64, tags ...string)

	// Timing records the execution time (latency) of a process.
	Timing(ctx context.Context, metric string, duration time.Duration, tags ...string)

	// CaptureError captures a system error and reports it to error tracking services (like Sentry).
	CaptureError(ctx context.Context, err error, tags ...string)

	// Close flushes and releases any resources used by the monitor.
	Close() error
}

// NewMonitor is a Factory function to create the corresponding Monitor.
func NewMonitor(cfg Config, log *zap.Logger) (Monitor, error) {
	switch cfg.Provider {
	case ProviderConsole:
		return newConsoleMonitor(cfg, log), nil
	case ProviderPrometheus:
		return newPrometheusMonitor(cfg)
	case ProviderDatadog:
		return newDatadogMonitor(cfg)
	case ProviderSentry:
		return newSentryMonitor(cfg)
	default:
		return nil, fmt.Errorf("monitoring: unsupported provider: %s", cfg.Provider)
	}
}

// ── Console Monitor Implementation ───────────────────────────────────────────

type consoleMonitor struct {
	log *zap.Logger
	cfg Config
}

func newConsoleMonitor(cfg Config, log *zap.Logger) *consoleMonitor {
	return &consoleMonitor{cfg: cfg, log: log}
}

func (m *consoleMonitor) Increment(ctx context.Context, metric string, tags ...string) {
	m.log.Info("Metric Increment", zap.String("metric", metric), zap.Strings("tags", tags))
}

func (m *consoleMonitor) Gauge(ctx context.Context, metric string, value float64, tags ...string) {
	m.log.Info("Metric Gauge", zap.String("metric", metric), zap.Float64("value", value), zap.Strings("tags", tags))
}

func (m *consoleMonitor) Timing(ctx context.Context, metric string, duration time.Duration, tags ...string) {
	m.log.Info("Metric Timing", zap.String("metric", metric), zap.Duration("duration", duration), zap.Strings("tags", tags))
}

func (m *consoleMonitor) CaptureError(ctx context.Context, err error, tags ...string) {
	m.log.Error("Captured Error", zap.Error(err), zap.Strings("tags", tags))
}

func (m *consoleMonitor) Close() error {
	return nil
}

// ── Prometheus Monitor Implementation (Stub) ─────────────────────────────────

type prometheusMonitor struct {
	cfg Config
}

func newPrometheusMonitor(cfg Config) (*prometheusMonitor, error) {
	// TODO: Initialize Prometheus metrics collectors
	return &prometheusMonitor{cfg: cfg}, nil
}

func (m *prometheusMonitor) Increment(ctx context.Context, metric string, tags ...string) {}
func (m *prometheusMonitor) Gauge(ctx context.Context, metric string, value float64, tags ...string) {}
func (m *prometheusMonitor) Timing(ctx context.Context, metric string, duration time.Duration, tags ...string) {}
func (m *prometheusMonitor) CaptureError(ctx context.Context, err error, tags ...string) {}
func (m *prometheusMonitor) Close() error { return nil }

// ── Datadog Monitor Implementation (Stub) ────────────────────────────────────

type datadogMonitor struct {
	cfg Config
}

func newDatadogMonitor(cfg Config) (*datadogMonitor, error) {
	// TODO: Initialize Datadog statsd client connection
	return &datadogMonitor{cfg: cfg}, nil
}

func (m *datadogMonitor) Increment(ctx context.Context, metric string, tags ...string) {}
func (m *datadogMonitor) Gauge(ctx context.Context, metric string, value float64, tags ...string) {}
func (m *datadogMonitor) Timing(ctx context.Context, metric string, duration time.Duration, tags ...string) {}
func (m *datadogMonitor) CaptureError(ctx context.Context, err error, tags ...string) {}
func (m *datadogMonitor) Close() error { return nil }

// ── Sentry Monitor Implementation (Stub) ─────────────────────────────────────

type sentryMonitor struct {
	cfg Config
}

func newSentryMonitor(cfg Config) (*sentryMonitor, error) {
	// TODO: Initialize Sentry client connection using DSN
	return &sentryMonitor{cfg: cfg}, nil
}

func (m *sentryMonitor) Increment(ctx context.Context, metric string, tags ...string) {}
func (m *sentryMonitor) Gauge(ctx context.Context, metric string, value float64, tags ...string) {}
func (m *sentryMonitor) Timing(ctx context.Context, metric string, duration time.Duration, tags ...string) {}
func (m *sentryMonitor) CaptureError(ctx context.Context, err error, tags ...string) {}
func (m *sentryMonitor) Close() error { return nil }
