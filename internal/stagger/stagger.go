// Package stagger spreads scan invocations across a time window to avoid
// thundering-herd effects when multiple port ranges are scheduled together.
package stagger

import (
	"context"
	"sync"
	"time"
)

// Entry holds a keyed callback that should be invoked after its assigned delay.
type Entry struct {
	Key   string
	Delay time.Duration
	Fn    func(ctx context.Context)
}

// Stagger distributes a set of entries evenly across a configurable window.
type Stagger struct {
	mu      sync.Mutex
	entries []Entry
	window  time.Duration
}

// New creates a Stagger that spreads entries across the given window duration.
// A zero or negative window is clamped to 1 millisecond.
func New(window time.Duration) *Stagger {
	if window <= 0 {
		window = time.Millisecond
	}
	return &Stagger{window: window}
}

// Register adds a keyed function to the stagger schedule.
// Delays are recalculated across all registered entries on each call.
func (s *Stagger) Register(key string, fn func(ctx context.Context)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing entry with same key.
	for i, e := range s.entries {
		if e.Key == key {
			s.entries = append(s.entries[:i], s.entries[i+1:]...)
			break
		}
	}
	s.entries = append(s.entries, Entry{Key: key, Fn: fn})
	s.recalculate()
}

// recalculate evenly distributes delays across the window. Must be called
// with s.mu held.
func (s *Stagger) recalculate() {
	n := len(s.entries)
	if n == 0 {
		return
	}
	step := s.window / time.Duration(n)
	for i := range s.entries {
		s.entries[i].Delay = time.Duration(i) * step
	}
}

// RunAll fires all registered entries with their assigned delays, respecting
// context cancellation. Each entry is run in its own goroutine after its delay.
func (s *Stagger) RunAll(ctx context.Context) {
	s.mu.Lock()
	copy := make([]Entry, len(s.entries))
	copy_ := copy
	copy(copy_, s.entries)
	s.mu.Unlock()

	for _, e := range copy_ {
		e := e
		go func() {
			select {
			case <-ctx.Done():
				return
			case <-time.After(e.Delay):
			}
			e.Fn(ctx)
		}()
	}
}

// Entries returns a snapshot of the current entries including their delays.
func (s *Stagger) Entries() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}
