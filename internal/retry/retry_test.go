package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

func fastPolicy() Policy {
	return Policy{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, Factor: 2.0}
}

func TestSuccessOnFirstAttempt(t *testing.T) {
	r := New(fastPolicy())
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetriesOnFailure(t *testing.T) {
	r := New(fastPolicy())
	calls := 0
	sentinel := errors.New("fail")
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 3 {
			return sentinel
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestAllAttemptsExhausted(t *testing.T) {
	r := New(fastPolicy())
	sentinel := errors.New("always fail")
	err := r.Do(context.Background(), func() error { return sentinel })
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestContextCancelledStopsRetry(t *testing.T) {
	r := New(Policy{MaxAttempts: 10, BaseDelay: 50 * time.Millisecond, MaxDelay: time.Second, Factor: 1.0})
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	go func() { time.Sleep(10 * time.Millisecond); cancel() }()
	err := r.Do(ctx, func() error { calls++; return errors.New("fail") })
	if err == nil {
		t.Fatal("expected error after cancel")
	}
	if calls >= 10 {
		t.Fatal("should have stopped before all attempts")
	}
}

func TestDefaultPolicyValues(t *testing.T) {
	p := DefaultPolicy()
	if p.MaxAttempts != 3 {
		t.Errorf("expected 3 attempts, got %d", p.MaxAttempts)
	}
	if p.Factor != 2.0 {
		t.Errorf("expected factor 2.0, got %f", p.Factor)
	}
}

func TestAttemptsReturnsConfigured(t *testing.T) {
	r := New(fastPolicy())
	if r.Attempts() != 3 {
		t.Errorf("expected 3, got %d", r.Attempts())
	}
}
