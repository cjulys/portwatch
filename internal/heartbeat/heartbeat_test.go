package heartbeat_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"portwatch/internal/heartbeat"
)

func TestHeartbeatWritesAtLeastOnce(t *testing.T) {
	var buf bytes.Buffer
	h := heartbeat.New(20*time.Millisecond, &buf)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	h.Run(ctx)

	if !strings.Contains(buf.String(), "[heartbeat] alive at") {
		t.Fatalf("expected heartbeat output, got: %q", buf.String())
	}
}

func TestHeartbeatStopsOnCancel(t *testing.T) {
	var buf bytes.Buffer
	h := heartbeat.New(10*time.Millisecond, &buf)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		h.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestDefaultIntervalFallback(t *testing.T) {
	h := heartbeat.New(0, nil)
	if h.Interval() != 60*time.Second {
		t.Fatalf("expected 60s default interval, got %v", h.Interval())
	}
}

func TestIntervalIsPreserved(t *testing.T) {
	h := heartbeat.New(5*time.Second, nil)
	if h.Interval() != 5*time.Second {
		t.Fatalf("expected 5s, got %v", h.Interval())
	}
}
