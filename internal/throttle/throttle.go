// Package throttle provides rate-limiting for alert notifications
// to prevent alert storms when many ports change simultaneously.
package throttle

import (
	"sync"
	"time"
)

// Throttle suppresses duplicate events for the same key within a cooldown window.
type Throttle struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New creates a Throttle with the given cooldown duration.
func New(cooldown time.Duration) *Throttle {
	return &Throttle{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event identified by key should be allowed through.
// Subsequent calls with the same key within the cooldown window return false.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.cooldown {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded timestamp for a specific key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Flush clears all recorded timestamps.
func (t *Throttle) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[string]time.Time)
}
