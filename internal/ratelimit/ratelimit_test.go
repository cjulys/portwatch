package ratelimit_test

import (
	"testing"
	"time"

	"portwatch/internal/ratelimit"
)

func TestFirstCallAlwaysAllowed(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected first call to be allowed")
	}
}

func TestSecondCallWithinIntervalRejected(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow()
	if l.Allow() {
		t.Fatal("expected second immediate call to be rejected")
	}
}

func TestCallAfterIntervalAllowed(t *testing.T) {
	l := ratelimit.New(20 * time.Millisecond)
	l.Allow()
	time.Sleep(30 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestStatsCountCorrectly(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow() // allowed
	l.Allow() // rejected
	l.Allow() // rejected
	allowed, rejected := l.Stats()
	if allowed != 1 {
		t.Fatalf("expected 1 allowed, got %d", allowed)
	}
	if rejected != 2 {
		t.Fatalf("expected 2 rejected, got %d", rejected)
	}
}

func TestResetClearsState(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow()
	l.Allow()
	l.Reset()
	allowed, rejected := l.Stats()
	if allowed != 0 || rejected != 0 {
		t.Fatalf("expected zeroed stats after reset, got allowed=%d rejected=%d", allowed, rejected)
	}
	if !l.Allow() {
		t.Fatal("expected Allow to succeed after Reset")
	}
}
