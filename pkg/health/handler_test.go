package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsHandler_ReturnsJSON(t *testing.T) {
	c := NewChecker()
	c.SetReady(true)
	c.SetHealthy(true)
	c.metrics.RecordReadiness()
	c.metrics.RecordLiveness()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	MetricsHandler(c)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var resp metricsResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&resp))
	assert.Equal(t, int64(1), resp.ReadinessChecks)
	assert.Equal(t, int64(1), resp.LivenessChecks)
	assert.True(t, resp.Ready)
	assert.True(t, resp.Healthy)
}

func TestLivenessHandler_WhenHealthy(t *testing.T) {
	c := NewChecker()
	c.SetHealthy(true)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	LivenessHandler(c)(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
	assert.Equal(t, int64(1), c.metrics.LivenessChecks)
}

func TestLivenessHandler_WhenUnhealthy(t *testing.T) {
	c := NewChecker()
	c.SetHealthy(false)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	LivenessHandler(c)(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.Equal(t, "not healthy", rec.Body.String())
	// also verify the liveness counter is incremented even on unhealthy responses
	assert.Equal(t, int64(1), c.metrics.LivenessChecks)
}
