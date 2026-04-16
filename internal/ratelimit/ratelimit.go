// Package ratelimit provides a token-bucket style rate limiter for
// controlling how frequently port scan cycles may be triggered.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls the rate at which events are allowed.
type Limiter struct {
	mu       sync.Mutex
	rate     time.Duration // minimum gap between allowed events
	last     time.Time
	allowed  int64
	rejected int64
}

// New creates a Limiter that allows at most one event per interval.
func New(interval time.Duration) *Limiter {
	return &Limiter{rate: interval}
}

// Allow returns true if enough time has elapsed since the last allowed event.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if l.last.IsZero() || now.Sub(l.last) >= l.rate {
		l.last = now
		l.allowed++
		return true
	}
	l.rejected++
	return false
}

// Stats returns the count of allowed and rejected calls.
func (l *Limiter) Stats() (allowed, rejected int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.allowed, l.rejected
}

// Reset clears the limiter state.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = time.Time{}
	l.allowed = 0
	l.rejected = 0
}
