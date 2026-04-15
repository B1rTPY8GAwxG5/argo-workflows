package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// ServerConfig holds configuration for the health HTTP server.
type ServerConfig struct {
	Port    int
	Checker *Checker
}

// Server is a lightweight HTTP server exposing health endpoints.
type Server struct {
	cfg    ServerConfig
	server *http.Server
}

// NewServer creates a new health Server.
func NewServer(cfg ServerConfig) *Server {
	return &Server{cfg: cfg}
}

// Start begins serving health endpoints and blocks until ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.cfg.Checker.LivenessHandler())
	mux.HandleFunc("/readyz", s.cfg.Checker.ReadinessHandler())
	mux.HandleFunc("/status", s.statusHandler())

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.cfg.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.WithField("port", s.cfg.Port).Info("starting health server")
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(shutCtx)
	case err := <-errCh:
		return err
	}
}

func (s *Server) statusHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(s.cfg.Checker.Status())
	}
}
