// Package cooldown provides per-key suppression of repeated events within a
// configurable quiet period. Unlike throttle, cooldown resets its timer on
// every new event so that a continuously-firing source stays silent until it
// has been quiet for the full window.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last-seen time for arbitrary string keys.
type Cooldown struct {
	mu      sync.Mutex
	window  time.Duration
	lastSeen map[string]time.Time
	now     func() time.Time
}

// New returns a Cooldown that suppresses events seen within window.
func New(window time.Duration) *Cooldown {
	return &Cooldown{
		window:   window,
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Record marks key as seen now and reports whether the event should be
// forwarded. The first occurrence always returns true. Subsequent calls
// return false until the key has been quiet for the full window.
func (c *Cooldown) Record(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	last, seen := c.lastSeen[key]
	c.lastSeen[key] = now

	if !seen {
		return true
	}
	return now.Sub(last) >= c.window
}

// Flush removes all tracked keys, resetting suppression state.
func (c *Cooldown) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastSeen = make(map[string]time.Time)
}

// Len returns the number of currently tracked keys.
func (c *Cooldown) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.lastSeen)
}
