package health_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFullLifecycle verifies that the health server, checker, and handlers
// work together correctly across a simulated component lifecycle.
// Note: increased sleep to 150ms to reduce flakiness on slower CI machines.
// Personal note: 100ms was still occasionally flaky in my local Docker env.
func TestFullLifecycle(t *testing.T) {
	port, err := freePort()
	require.NoError(t, err)

	checker := NewChecker()
	server := NewServer(fmt.Sprintf(":%d", port), checker)

	go func() {
		_ = server.Start()
	}()

	// Allow server to start.
	time.Sleep(150 * time.Millisecond)

	base := fmt.Sprintf("http://localhost:%d", port)

	t.Run("liveness returns unhealthy before ready", func(t *testing.T) {
		resp, err := http.Get(base + "/healthz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	t.Run("readiness returns not ready initially", func(t *testing.T) {
		resp, err := http.Get(base + "/readyz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
	})

	// Simulate component becoming healthy and ready.
	checker.SetHealthy(true)
	checker.SetReady(true)

	t.Run("liveness returns healthy after SetHealthy", func(t *testing.T) {
		resp, err := http.Get(base + "/healthz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("readiness returns ready after SetReady", func(t *testing.T) {
		resp, err := http.Get(base + "/readyz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("metrics endpoint returns valid JSON with recorded observations", func(t *testing.T) {
		resp, err := http.Get(base + "/metrics")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var payload map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&payload)
		require.NoError(t, err)

		_, hasReadiness := payload["readiness_checks_total"]
		_, hasLiveness := payload["liveness_checks_total"]
		_, hasUptime := payload["uptime_seconds"]
		assert.True(t, hasReadiness, "expected readiness_checks_total in metrics")
		assert.True(t, hasLiveness, "expected liveness_checks_total in metrics")
		assert.True(t, hasUptime, "expected uptime_seconds in metrics")
	})

	require.NoError(t, server.Stop())
}
