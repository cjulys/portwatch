package backoff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"portwatch/internal/backoff"
)

func TestSuccessOnSecondAttempt(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: time.Millisecond,
		MaxInterval:     5 * time.Millisecond,
		Multiplier:      2,
		MaxAttempts:     5,
	}
	b := backoff.New(p)
	attempts := 0
	var lastErr error
	for b.Next(context.Background()) {
		attempts++
		if attempts >= 2 {
			lastErr = nil
			break
		}
		lastErr = errors.New("transient")
	}
	if lastErr != nil {
		t.Fatalf("expected success, got %v", lastErr)
	}
	if attempts != 2 {
		t.Fatalf("expected 2 attempts, got %d", attempts)
	}
}

func TestAllAttemptsExhaustedReturnsError(t *testing.T) {
	p := backoff.Policy{
		InitialInterval: time.Millisecond,
		MaxInterval:     2 * time.Millisecond,
		Multiplier:      1,
		MaxAttempts:     3,
	}
	b := backoff.New(p)
	var err error
	for b.Next(context.Background()) {
		err = errors.New("always fails")
	}
	if err == nil {
		t.Fatal("expected error after exhausted attempts")
	}
	if b.Attempt() != 3 {
		t.Fatalf("expected 3 attempts recorded, got %d", b.Attempt())
	}
}
