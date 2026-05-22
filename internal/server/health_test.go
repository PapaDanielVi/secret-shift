package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthz(t *testing.T) {
	h := NewHealthServer(0)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	h.handleHealthz(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "ok", body["status"])
}

func TestReadyz_NotReady_NoSync(t *testing.T) {
	h := NewHealthServer(0)
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	h.handleReadyz(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "not_ready", body["status"])
}

func TestReadyz_Ready(t *testing.T) {
	h := NewHealthServer(0)
	h.ReportSyncSuccess()

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	h.handleReadyz(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "ready", body["status"])
	assert.NotNil(t, body["last_sync_time"])
}

func TestReadyz_NotReady_AfterError(t *testing.T) {
	h := NewHealthServer(0)
	h.ReportSyncSuccess()
	h.ReportSyncError(assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()

	h.handleReadyz(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "not_ready", body["status"])
}

func TestStatus(t *testing.T) {
	h := NewHealthServer(8080)
	h.ReportSyncSuccess()
	h.ReportSyncError(assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	w := httptest.NewRecorder()

	h.handleStatus(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result statusResponse
	require.NoError(t, json.Unmarshal(body, &result))

	assert.Equal(t, "ok", result.Status)
	assert.Equal(t, int64(1), result.SyncCount)
	assert.Equal(t, int64(1), result.ErrorCount)
	assert.NotEmpty(t, result.LastSyncTime)
	assert.NotEmpty(t, result.LastSyncError)
}

func TestReportSyncSuccess(t *testing.T) {
	h := NewHealthServer(0)
	h.ReportSyncSuccess()

	h.mu.RLock()
	assert.False(t, h.lastSyncTime.IsZero())
	require.NoError(t, h.lastSyncError)
	h.mu.RUnlock()

	assert.Equal(t, int64(1), h.syncCount.Load())
}

func TestReportSyncError(t *testing.T) {
	h := NewHealthServer(0)
	testErr := assert.AnError
	h.ReportSyncError(testErr)

	h.mu.RLock()
	assert.False(t, h.lastSyncTime.IsZero())
	assert.Equal(t, testErr, h.lastSyncError)
	h.mu.RUnlock()

	assert.Equal(t, int64(1), h.errorCount.Load())
}

func TestHealthServer_Integration(t *testing.T) {
	h := NewHealthServer(18080)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		assert.NoError(t, h.Start(ctx))
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Test healthz
	resp, err := http.Get("http://localhost:18080/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test readyz (not ready yet)
	resp, err = http.Get("http://localhost:18080/readyz")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	// Report success
	h.ReportSyncSuccess()

	// Test readyz (now ready)
	resp, err = http.Get("http://localhost:18080/readyz")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	cancel()
}
