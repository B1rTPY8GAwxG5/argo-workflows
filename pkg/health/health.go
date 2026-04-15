package health

import (
	"net/http"
	"sync/atomic"
	"time"
)

// Status represents the health status of the server.
type Status struct {
	Ready   bool      `json:"ready"`
	Healthy bool      `json:"healthy"`
	CheckedAt time.Time `json:"checkedAt"`
}

// Checker manages readiness and liveness state.
type Checker struct {
	ready   atomic.Bool
	healthy atomic.Bool
}

// NewChecker creates a new Checker with healthy set to true.
func NewChecker() *Checker {
	c := &Checker{}
	c.healthy.Store(true)
	return c
}

// SetReady marks the server as ready to serve traffic.
func (c *Checker) SetReady(ready bool) {
	c.ready.Store(ready)
}

// SetHealthy marks the server as healthy.
func (c *Checker) SetHealthy(healthy bool) {
	c.healthy.Store(healthy)
}

// Status returns the current health status.
func (c *Checker) Status() Status {
	return Status{
		Ready:     c.ready.Load(),
		Healthy:   c.healthy.Load(),
		CheckedAt: time.Now().UTC(),
	}
}

// ReadinessHandler returns an HTTP handler for readiness probes.
func (c *Checker) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !c.ready.Load() {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}

// LivenessHandler returns an HTTP handler for liveness probes.
func (c *Checker) LivenessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !c.healthy.Load() {
			http.Error(w, "not healthy", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
