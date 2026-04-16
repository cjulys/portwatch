// Package limiter enforces a maximum number of alerts per time window.
package limiter

import (
	"sync"
	"time"
)

// Limiter tracks event counts per key within a rolling window.
type Limiter struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	buckets map[string][]time.Time
	now     func() time.Time
}

// New returns a Limiter that allows at most max events per key within window.
func New(window time.Duration, max int) *Limiter {
	return &Limiter{
		window:  window,
		max:     max,
		buckets: make(map[string][]time.Time),
		now:     time.Now,
	}
}

// Allow returns true if the event for key is within the allowed rate.
// It records the event if allowed.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	times := l.buckets[key]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= l.max {
		l.buckets[key] = filtered
		return false
	}

	l.buckets[key] = append(filtered, now)
	return true
}

// Count returns the number of recorded events for key within the current window.
func (l *Limiter) Count(key string) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)
	count := 0
	for _, t := range l.buckets[key] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all recorded events.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buckets = make(map[string][]time.Time)
}
