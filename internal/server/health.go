// Package server implements a simple HTTP server that provides health and readiness endpoints for the sync pipeline.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

const (
	status = "status"
)

// HealthServer provides health and readiness endpoints for the sync pipeline.
type HealthServer struct {
	port      int
	startTime time.Time

	mu            sync.RWMutex
	lastSyncTime  time.Time
	lastSyncError error
	syncCount     atomic.Int64
	errorCount    atomic.Int64
}

// NewHealthServer creates a new health server on the given port.
func NewHealthServer(port int) *HealthServer {
	return &HealthServer{
		port:      port,
		startTime: time.Now(),
	}
}

// ReportSyncSuccess records a successful sync.
func (h *HealthServer) ReportSyncSuccess() {
	h.mu.Lock()
	h.lastSyncTime = time.Now()
	h.lastSyncError = nil
	h.mu.Unlock()
	h.syncCount.Add(1)
}

// ReportSyncError records a failed sync.
func (h *HealthServer) ReportSyncError(err error) {
	h.mu.Lock()
	h.lastSyncTime = time.Now()
	h.lastSyncError = err
	h.mu.Unlock()
	h.errorCount.Add(1)
}

// Start begins the HTTP server and blocks until ctx is cancelled.
func (h *HealthServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", h.handleHealthz)
	mux.HandleFunc("/readyz", h.handleReadyz)
	mux.HandleFunc("/status", h.handleStatus)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", h.port),
		Handler:           mux,
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("health server shutdown error", "err", err)
		}
	}()

	slog.Info("Starting health server", "port", h.port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("health server: %w", err)
	}
	return nil
}

func (h *HealthServer) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]any{
		status: "ok",
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode healthz response", "err", err)
	}
}

func (h *HealthServer) handleReadyz(w http.ResponseWriter, _ *http.Request) {
	h.mu.RLock()
	lastSync := h.lastSyncTime
	lastErr := h.lastSyncError
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	if lastSync.IsZero() {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp := map[string]any{
			status:   "not_ready",
			"reason": "no successful sync yet",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode readyz response", "err", err)
		}
		return
	}

	if lastErr != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		resp := map[string]any{
			status:   "not_ready",
			"reason": "last sync failed",
			"error":  lastErr.Error(),
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			slog.Error("failed to encode readyz response", "err", err)
		}
		return
	}

	resp := map[string]any{
		status:           "ready",
		"last_sync_time": lastSync.Format(time.RFC3339),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode readyz response", "err", err)
	}
}

type statusResponse struct {
	Status        string `json:"status"`
	StartTime     string `json:"start_time"`
	LastSyncTime  string `json:"last_sync_time,omitempty"`
	LastSyncError string `json:"last_sync_error,omitempty"`
	SyncCount     int64  `json:"sync_count"`
	ErrorCount    int64  `json:"error_count"`
}

func (h *HealthServer) handleStatus(w http.ResponseWriter, _ *http.Request) {
	h.mu.RLock()
	lastSync := h.lastSyncTime
	lastErr := h.lastSyncError
	h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	resp := &statusResponse{
		Status:     "ok",
		StartTime:  h.startTime.Format(time.RFC3339),
		SyncCount:  h.syncCount.Load(),
		ErrorCount: h.errorCount.Load(),
	}

	if !lastSync.IsZero() {
		resp.LastSyncTime = lastSync.Format(time.RFC3339)
	}
	if lastErr != nil {
		resp.LastSyncError = lastErr.Error()
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode status response", "err", err)
	}
}
