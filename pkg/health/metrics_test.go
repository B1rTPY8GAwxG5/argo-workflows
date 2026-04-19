package health

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMetrics_InitialisedCorrectly(t *testing.T) {
	m := newMetrics()
	assert.NotNil(t, m)
	assert.Equal(t, int64(0), m.ReadinessChecks)
	assert.Equal(t, int64(0), m.LivenessChecks)
	assert.False(t, m.startTime.IsZero())
}

func TestMetrics_RecordReadiness(t *testing.T) {
	m := newMetrics()
	m.RecordReadiness()
	m.RecordReadiness()
	assert.Equal(t, int64(2), m.ReadinessChecks)
	assert.False(t, m.LastReadyTime.IsZero())
}

func TestMetrics_RecordLiveness(t *testing.T) {
	m := newMetrics()
	m.RecordLiveness()
	assert.Equal(t, int64(1), m.LivenessChecks)
	assert.False(t, m.LastHealthyTime.IsZero())
}

func TestMetrics_Uptime(t *testing.T) {
	m := newMetrics()
	time.Sleep(10 * time.Millisecond)
	uptime := m.Uptime()
	assert.True(t, uptime >= 10*time.Millisecond, "expected uptime >= 10ms, got %v", uptime)
}

func TestMetrics_Snapshot(t *testing.T) {
	m := newMetrics()
	time.Sleep(10 * time.Millisecond)
	m.RecordReadiness()
	m.RecordLiveness()

	snap := m.Snapshot()
	assert.Equal(t, int64(1), snap.ReadinessChecks)
	assert.Equal(t, int64(1), snap.LivenessChecks)
	assert.True(t, snap.UptimeSeconds >= 0)

	// Verify it is a copy — mutating original should not affect snapshot.
	m.RecordReadiness()
	assert.Equal(t, int64(1), snap.ReadinessChecks)
}
