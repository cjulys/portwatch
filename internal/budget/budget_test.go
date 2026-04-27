package budget

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowUnderBudget(t *testing.T) {
	b := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !b.Allow(context.Background()) {
			t.Fatalf("call %d should be allowed", i+1)
		}
	}
}

func TestAllowBlocksWhenExhausted(t *testing.T) {
	b := New(2, time.Minute)
	b.Allow(context.Background())
	b.Allow(context.Background())
	if b.Allow(context.Background()) {
		t.Fatal("third call should be blocked")
	}
}

func TestAllowAfterWindowExpires(t *testing.T) {
	now := time.Now()
	b := New(1, time.Second)
	b.clock = fixedClock(now)
	b.Allow(context.Background())
	// advance past the window
	b.clock = fixedClock(now.Add(2 * time.Second))
	if !b.Allow(context.Background()) {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestRemainingDecrementsAndRecovers(t *testing.T) {
	now := time.Now()
	b := New(3, time.Second)
	b.clock = fixedClock(now)
	if b.Remaining() != 3 {
		t.Fatalf("expected 3 remaining, got %d", b.Remaining())
	}
	b.Allow(context.Background())
	if b.Remaining() != 2 {
		t.Fatalf("expected 2 remaining, got %d", b.Remaining())
	}
	b.clock = fixedClock(now.Add(2 * time.Second))
	if b.Remaining() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", b.Remaining())
	}
}

func TestResetRestoresBudget(t *testing.T) {
	b := New(2, time.Minute)
	b.Allow(context.Background())
	b.Allow(context.Background())
	b.Reset()
	if !b.Allow(context.Background()) {
		t.Fatal("call after Reset should be allowed")
	}
}

func TestZeroMaxDisablesLimiting(t *testing.T) {
	b := New(0, time.Minute)
	for i := 0; i < 100; i++ {
		if !b.Allow(context.Background()) {
			t.Fatal("zero max should allow all calls")
		}
	}
	if b.Remaining() != -1 {
		t.Fatal("Remaining should return -1 when limiting is disabled")
	}
}

func TestCancelledContextReturnsFalse(t *testing.T) {
	b := New(10, time.Minute)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if b.Allow(ctx) {
		t.Fatal("cancelled context should not be allowed")
	}
}

func TestFallbackWriterReceivesMessage(t *testing.T) {
	var buf bytes.Buffer
	b := New(1, time.Minute)
	b.fallback = &buf
	b.Allow(context.Background())
	b.Allow(context.Background()) // exhausts budget, triggers write
	if buf.Len() == 0 {
		t.Fatal("expected fallback message when budget exhausted")
	}
}
