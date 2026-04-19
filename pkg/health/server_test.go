package health_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/health"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func freePort() int {
	// Use a high ephemeral port for tests to avoid conflicts.
	return 19876
}

func TestServer_StartsAndResponds(t *testing.T) {
	checker := health.NewChecker()
	checker.SetReady(true)

	cfg := health.ServerConfig{
		Port:    freePort(),
		Checker: checker,
	}
	srv := health.NewServer(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start(ctx)
	}()

	// Give the server a moment to start.
	// Increased from 100ms to 200ms to reduce flakiness on slower CI machines.
	time.Sleep(200 * time.Millisecond)

	base := fmt.Sprintf("http://localhost:%d", freePort())

	resp, err := http.Get(base + "/healthz")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp, err = http.Get(base + "/readyz")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	resp, err = http.Get(base + "/status")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	_ = resp.Body.Close()

	cancel()
	select {
	case err := <-errCh:
		assert.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}
