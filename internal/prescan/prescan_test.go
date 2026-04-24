package prescan_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"portwatch/internal/prescan"
)

func TestFirstEvaluateAlwaysScans(t *testing.T) {
	c := prescan.New(5 * time.Second)
	r := c.Evaluate(context.Background())
	if !r.ShouldScan {
		t.Fatalf("expected ShouldScan=true on first call, got reason=%q", r.Reason)
	}
	if r.Reason != "ok" {
		t.Errorf("expected reason \"ok\", got %q", r.Reason)
	}
}

func TestSecondCallWithinIntervalSuppressed(t *testing.T) {
	c := prescan.New(1 * time.Hour)
	c.Evaluate(context.Background()) // prime last-scan
	r := c.Evaluate(context.Background())
	if r.ShouldScan {
		t.Fatal("expected ShouldScan=false within min interval")
	}
	if r.Reason != "min interval not elapsed" {
		t.Errorf("unexpected reason: %q", r.Reason)
	}
}

func TestCallAfterIntervalAllowed(t *testing.T) {
	c := prescan.New(1 * time.Millisecond)
	c.Evaluate(context.Background())
	time.Sleep(5 * time.Millisecond)
	r := c.Evaluate(context.Background())
	if !r.ShouldScan {
		t.Fatalf("expected ShouldScan=true after interval elapsed, got %q", r.Reason)
	}
}

func TestZeroIntervalNeverSuppresses(t *testing.T) {
	c := prescan.New(0)
	for i := 0; i < 5; i++ {
		r := c.Evaluate(context.Background())
		if !r.ShouldScan {
			t.Fatalf("call %d: expected ShouldScan=true with zero interval", i)
		}
	}
}

func TestCancelledContextReturnsFalse(t *testing.T) {
	c := prescan.New(0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := c.Evaluate(ctx)
	if r.ShouldScan {
		t.Fatal("expected ShouldScan=false for cancelled context")
	}
	if r.Reason != "context cancelled" {
		t.Errorf("unexpected reason: %q", r.Reason)
	}
}

func TestResetAllowsImmediateRescan(t *testing.T) {
	c := prescan.New(1 * time.Hour)
	c.Evaluate(context.Background())
	c.Reset()
	r := c.Evaluate(context.Background())
	if !r.ShouldScan {
		t.Fatal("expected ShouldScan=true after Reset")
	}
}

func TestWithWriterReturnsSelf(t *testing.T) {
	var buf bytes.Buffer
	c := prescan.New(0)
	got := c.WithWriter(&buf)
	if got == nil {
		t.Fatal("WithWriter returned nil")
	}
}

func TestEvaluatedAtIsRecent(t *testing.T) {
	c := prescan.New(0)
	before := time.Now().UTC()
	r := c.Evaluate(context.Background())
	after := time.Now().UTC()
	if r.EvaluatedAt.Before(before) || r.EvaluatedAt.After(after) {
		t.Errorf("EvaluatedAt %v not within [%v, %v]", r.EvaluatedAt, before, after)
	}
}
