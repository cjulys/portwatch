package stagger_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"portwatch/internal/stagger"
)

func TestRegisterAssignsZeroDelayForSingleEntry(t *testing.T) {
	s := stagger.New(100 * time.Millisecond)
	s.Register("a", func(ctx context.Context) {})
	entries := s.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Delay != 0 {
		t.Errorf("single entry should have zero delay, got %v", entries[0].Delay)
	}
}

func TestRegisterSpreadsDelaysEvenly(t *testing.T) {
	s := stagger.New(100 * time.Millisecond)
	s.Register("a", func(ctx context.Context) {})
	s.Register("b", func(ctx context.Context) {})
	s.Register("c", func(ctx context.Context) {})

	entries := s.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	expected := []time.Duration{0, 33*time.Millisecond + 333*time.Microsecond, 66*time.Millisecond + 666*time.Microsecond}
	for i, e := range entries {
		diff := e.Delay - expected[i]
		if diff < 0 {
			diff = -diff
		}
		if diff > time.Millisecond {
			t.Errorf("entry %d delay %v too far from expected %v", i, e.Delay, expected[i])
		}
	}
}

func TestRegisterReplacesExistingKey(t *testing.T) {
	s := stagger.New(100 * time.Millisecond)
	s.Register("a", func(ctx context.Context) {})
	s.Register("a", func(ctx context.Context) {})
	if len(s.Entries()) != 1 {
		t.Errorf("duplicate key should replace, expected 1 entry")
	}
}

func TestRunAllInvokesAllEntries(t *testing.T) {
	s := stagger.New(20 * time.Millisecond)
	var count int64
	for _, key := range []string{"x", "y", "z"} {
		key := key
		_ = key
		s.Register(key, func(ctx context.Context) {
			atomic.AddInt64(&count, 1)
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	s.RunAll(ctx)
	time.Sleep(150 * time.Millisecond)
	if got := atomic.LoadInt64(&count); got != 3 {
		t.Errorf("expected 3 invocations, got %d", got)
	}
}

func TestRunAllRespectsContextCancellation(t *testing.T) {
	s := stagger.New(500 * time.Millisecond)
	var mu sync.Mutex
	called := false
	s.Register("slow", func(ctx context.Context) {
		mu.Lock()
		called = true
		mu.Unlock()
	})
	ctx, cancel := context.WithCancel(context.Background())
	s.RunAll(ctx)
	cancel()
	time.Sleep(50 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if called {
		t.Error("entry should not have been called after context cancellation")
	}
}

func TestZeroWindowClamped(t *testing.T) {
	s := stagger.New(0)
	s.Register("a", func(ctx context.Context) {})
	// Should not panic; delay for single entry is always 0 regardless.
	if len(s.Entries()) != 1 {
		t.Error("expected 1 entry")
	}
}
