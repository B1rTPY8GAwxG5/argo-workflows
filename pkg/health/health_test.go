package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/argoproj/argo-workflows/v3/pkg/health"
	"github.com/stretchr/testify/assert"
)

func TestNewChecker_DefaultState(t *testing.T) {
	c := health.NewChecker()
	status := c.Status()
	assert.False(t, status.Ready, "expected not ready by default")
	assert.True(t, status.Healthy, "expected healthy by default")
}

func TestChecker_SetReady(t *testing.T) {
	c := health.NewChecker()
	c.SetReady(true)
	assert.True(t, c.Status().Ready)
	c.SetReady(false)
	assert.False(t, c.Status().Ready)
}

func TestChecker_SetHealthy(t *testing.T) {
	c := health.NewChecker()
	c.SetHealthy(false)
	assert.False(t, c.Status().Healthy)
	c.SetHealthy(true)
	assert.True(t, c.Status().Healthy)
}

func TestReadinessHandler_NotReady(t *testing.T) {
	c := health.NewChecker()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	c.ReadinessHandler()(rr, req)
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
}

func TestReadinessHandler_Ready(t *testing.T) {
	c := health.NewChecker()
	c.SetReady(true)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()
	c.ReadinessHandler()(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLivenessHandler_Healthy(t *testing.T) {
	c := health.NewChecker()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	c.LivenessHandler()(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLivenessHandler_Unhealthy(t *testing.T) {
	c := health.NewChecker()
	c.SetHealthy(false)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	c.LivenessHandler()(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
