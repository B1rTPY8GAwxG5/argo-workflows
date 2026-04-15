package health

import (
	"sync/atomic"
	"time"
)

// Metrics holds runtime health metrics for the checker.
type Metrics struct {
	ReadinessChecks  int64
	LivenessChecks   int64
	LastReadyTime    time.Time
	LastHealthyTime  time.Time
	UptimeSeconds    int64
	startTime        time.Time
}

// newMetrics initialises a Metrics instance with the current time as start.
func newMetrics() *Metrics {
	now := time.Now()
	return &Metrics{
		startTime:     now,
		LastReadyTime: now,
	}
}

// RecordReadiness increments the readiness check counter.
func (m *Metrics) RecordReadiness() {
	atomic.AddInt64(&m.ReadinessChecks, 1)
	m.LastReadyTime = time.Now()
}

// RecordLiveness increments the liveness check counter.
func (m *Metrics) RecordLiveness() {
	atomic.AddInt64(&m.LivenessChecks, 1)
	m.LastHealthyTime = time.Now()
}

// Uptime returns the duration since the metrics were initialised.
func (m *Metrics) Uptime() time.Duration {
	return time.Since(m.startTime)
}

// Snapshot returns a copy of the current metrics with the uptime refreshed.
func (m *Metrics) Snapshot() Metrics {
	copy := *m
	copy.UptimeSeconds = int64(m.Uptime().Seconds())
	return copy
}
