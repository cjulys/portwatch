package backoff

import (
	"context"
	"testing"
	"time"
)

func fastPolicy() Policy {
	return Policy{
		InitialInterval: time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		MaxAttempts:     4,
	}
}

func TestNextAdvancesAttempt(t *testing.T) {
	b := New(fastPolicy())
	if !b.Next(context.Background()) {
		t.Fatal("first Next should return true")
	}
	if b.Attempt() != 1 {
		t.Fatalf("expected attempt 1, got %d", b.Attempt())
	}
}

func TestExhaustedReturnsFalse(t *testing.T) {
	b := New(fastPolicy())
	count := 0
	for b.Next(context.Background()) {
		count++
	}
	if count != 4 {
		t.Fatalf("expected 4 attempts, got %d", count)
	}
}

func TestResetRestartsSequence(t *testing.T) {
	b := New(fastPolicy())
	b.Next(context.Background())
	b.Next(context.Background())
	b.Reset()
	if b.Attempt() != 0 {
		t.Fatalf("expected 0 after reset, got %d", b.Attempt())
	}
	if !b.Next(context.Background()) {
		t.Fatal("Next after reset should return true")
	}
}

func TestContextCancelledStopsBackoff(t *testing.T) {
	p := Policy{
		InitialInterval: 10 * time.Second,
		MaxInterval:     10 * time.Second,
		Multiplier:      1,
		MaxAttempts:     5,
	}
	b := New(p)
	b.Next(context.Background()) // first call – no sleep
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if b.Next(ctx) {
		t.Fatal("Next should return false when context is cancelled")
	}
}

func TestIntervalClampsToMax(t *testing.T) {
	p := Policy{
		InitialInterval: time.Millisecond,
		MaxInterval:     3 * time.Millisecond,
		Multiplier:      100,
		MaxAttempts:     3,
	}
	b := New(p)
	b.attempt = 2
	got := b.interval()
	if got > p.MaxInterval {
		t.Fatalf("interval %v exceeds max %v", got, p.MaxInterval)
	}
}

func TestDefaultPolicyIsReasonable(t *testing.T) {
	p := DefaultPolicy()
	if p.InitialInterval <= 0 {
		t.Fatal("initial interval must be positive")
	}
	if p.MaxInterval < p.InitialInterval {
		t.Fatal("max interval must be >= initial interval")
	}
	if p.MaxAttempts <= 0 {
		t.Fatal("max attempts must be positive")
	}
}
