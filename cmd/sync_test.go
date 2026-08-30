package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/PapaDanielVi/secret-shift/internal/server"
)

func TestRunServerSyncRunsImmediately(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	called := make(chan struct{})
	done := make(chan struct{})

	go func() {
		runServerSync(ctx, func(context.Context) error {
			close(called)
			return nil
		}, server.NewHealthServer(0))
		close(done)
	}()

	select {
	case <-called:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("server sync did not run immediately")
	}

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("server sync did not stop after cancellation")
	}
}
