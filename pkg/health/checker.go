package health

import "sync/atomic"

// Checker tracks the readiness and liveness state of a component.
type Checker struct {
	ready   atomic.Bool
	healthy atomic.Bool
}

// NewChecker creates a new Checker with the given initial states.
func NewChecker(ready, healthy bool) *Checker {
	c := &Checker{}
	c.ready.Store(ready)
	c.healthy.Store(healthy)
	return c
}

// SetReady sets the readiness state.
func (c *Checker) SetReady(ready bool) {
	c.ready.Store(ready)
}

// SetHealthy sets the liveness state.
func (c *Checker) SetHealthy(healthy bool) {
	c.healthy.Store(healthy)
}

// IsReady returns the current readiness state.
func (c *Checker) IsReady() bool {
	return c.ready.Load()
}

// IsHealthy returns the current liveness state.
func (c *Checker) IsHealthy() bool {
	return c.healthy.Load()
}
