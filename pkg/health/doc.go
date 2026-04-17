// Package health provides lightweight liveness and readiness probes for
// Argo Workflows server components.
//
// Usage:
//
//	checker := health.NewChecker()
//
//	// Mark ready once initialisation is complete.
//	checker.SetReady(true)
//
//	// Start the health HTTP server on port 8080.
//	srv := health.NewServer(health.ServerConfig{
//		Port:    8080,
//		Checker: checker,
//	})
//	if err := srv.Start(ctx); err != nil {
//		log.Fatal(err)
//	}
//
// Endpoints exposed:
//
//	GET /healthz  — liveness probe (200 OK / 500 Internal Server Error)
//	GET /readyz   — readiness probe (200 OK / 503 Service Unavailable)
//	GET /status   — JSON status object
//
// Note: /status is unauthenticated by design; avoid exposing sensitive
// runtime details in its response payload.
//
// Personal note: consider adding a /version endpoint in the future that
// returns build metadata (git commit, build date) for easier debugging
// in multi-replica deployments.
package health
