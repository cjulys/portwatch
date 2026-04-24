// Package correlate groups related port-change diffs into correlated events.
// When multiple ports change within a short burst window they are bundled
// together under a single correlation ID so downstream handlers can treat
// them as one logical event (e.g. a service restart opening several ports).
package correlate

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Event holds a set of diffs that arrived within the same burst window.
type Event struct {
	ID        string
	Diffs     []scanner.Diff
	CreatedAt time.Time
}

// Correlator accumulates diffs and flushes correlated events after a quiet
// window expires or when Flush is called explicitly.
type Correlator struct {
	mu      sync.Mutex
	window  time.Duration
	buf     []scanner.Diff
	timer   *time.Timer
	onFlush func(Event)
}

// New creates a Correlator. window is the quiet period after the last Add
// before an automatic flush fires. onFlush is called with each correlated
// event; it must not block.
func New(window time.Duration, onFlush func(Event)) *Correlator {
	return &Correlator{
		window:  window,
		onFlush: onFlush,
	}
}

// Add appends diffs to the current burst buffer and resets the flush timer.
func (c *Correlator) Add(diffs []scanner.Diff) {
	if len(diffs) == 0 {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.buf = append(c.buf, diffs...)

	if c.timer != nil {
		c.timer.Reset(c.window)
		return
	}
	c.timer = time.AfterFunc(c.window, c.autoFlush)
}

// Flush immediately emits any buffered diffs as a correlated event.
func (c *Correlator) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flushLocked()
}

func (c *Correlator) autoFlush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.flushLocked()
}

func (c *Correlator) flushLocked() {
	if len(c.buf) == 0 {
		return
	}
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	ev := Event{
		ID:        newID(),
		Diffs:     c.buf,
		CreatedAt: time.Now().UTC(),
	}
	c.buf = nil
	c.onFlush(ev)
}

func newID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
