package health

import (
	"encoding/json"
	"net/http"
)

// metricsResponse is the JSON shape returned by the metrics endpoint.
type metricsResponse struct {
	ReadinessChecks int64  `json:"readiness_checks"`
	LivenessChecks  int64  `json:"liveness_checks"`
	UptimeSeconds   int64  `json:"uptime_seconds"`
	Ready           bool   `json:"ready"`
	Healthy         bool   `json:"healthy"`
}

// MetricsHandler returns an http.HandlerFunc that exposes health metrics as JSON.
func MetricsHandler(c *Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := c.metrics.Snapshot()
		resp := metricsResponse{
			ReadinessChecks: snap.ReadinessChecks,
			LivenessChecks:  snap.LivenessChecks,
			UptimeSeconds:   snap.UptimeSeconds,
			Ready:           c.isReady(),
			Healthy:         c.isHealthy(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// LivenessHandler returns an http.HandlerFunc for liveness probes.
// It records each check in the metrics and returns 200 when healthy, 503 otherwise.
func LivenessHandler(c *Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.metrics.RecordLiveness()
		if c.isHealthy() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("not healthy"))
	}
}
