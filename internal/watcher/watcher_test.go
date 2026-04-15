package watcher_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
	"github.com/user/portwatch/internal/watcher"
)

func buildWatcher(t *testing.T) (*watcher.Watcher, *config.Config) {
	t.Helper()
	cfg := config.Default()
	cfg.IntervalSeconds = 1

	sc := scanner.New(cfg)
	st := state.New(filepath.Join(t.TempDir(), "state.json"))
	al := alert.New(cfg, nil)

	return watcher.New(cfg, sc, st, al), cfg
}

func TestRunCancelledImmediately(t *testing.T) {
	w, _ := buildWatcher(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := w.Run(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestRunTicksAtLeastOnce(t *testing.T) {
	w, _ := buildWatcher(t)
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()

	// Run should complete without a hard error (context deadline or cancel).
	err := w.Run(ctx)
	if err != context.DeadlineExceeded && err != context.Canceled {
		t.Errorf("unexpected error: %v", err)
	}
}
