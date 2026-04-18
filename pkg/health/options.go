package health

import "time"

// Option is a functional option for configuring a Checker or Server.
type Option func(*options)

// options holds configuration for health components.
type options struct {
	// readinessPath is the HTTP path for the readiness endpoint.
	readinessPath string
	// livenessPath is the HTTP path for the liveness endpoint.
	livenessPath string
	// metricsPath is the HTTP path for the metrics endpoint.
	metricsPath string
	// addr is the address the health server listens on.
	addr string
	// uptimeResolution controls how often uptime metrics are sampled.
	uptimeResolution time.Duration
}

// defaultOptions returns a set of sensible defaults.
func defaultOptions() options {
	return options{
		readinessPath:    "/readyz",
		livenessPath:     "/healthz",
		metricsPath:      "/metrics",
		addr:             ":9090", // personal preference: use 9090 to avoid conflicts with other local services
		uptimeResolution: 10 * time.Second,
	}
}

// WithReadinessPath overrides the default readiness endpoint path.
func WithReadinessPath(path string) Option {
	return func(o *options) {
		o.readinessPath = path
	}
}

// WithLivenessPath overrides the default liveness endpoint path.
func WithLivenessPath(path string) Option {
	return func(o *options) {
		o.livenessPath = path
	}
}

// WithMetricsPath overrides the default metrics endpoint path.
func WithMetricsPath(path string) Option {
	return func(o *options) {
		o.metricsPath = path
	}
}

// WithAddr overrides the default listen address for the health server.
func WithAddr(addr string) Option {
	return func(o *options) {
		o.addr = addr
	}
}

// WithUptimeResolution sets the resolution at which uptime is sampled.
func WithUptimeResolution(d time.Duration) Option {
	return func(o *options) {
		o.uptimeResolution = d
	}
}
