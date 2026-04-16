// Package rollup batches multiple diff events within a time window
// and emits a single summarised event, reducing alert noise.
package rollup

import (
	"sync"
	"time"

	"portwatch/internal/scanner"
)

// Event holds a rolled-up collection of diffs.
type Event struct {
	Opened []scanner.Port
	Closed []scanner.Port
	At     time.Time
}

// Rollup accumulates diffs and flushes them after a window elapses.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	opened  []scanner.Port
	closed  []scanner.Port
	timer   *time.Timer
	flushFn func(Event)
}

// New creates a Rollup that calls fn at most once per window.
func New(window time.Duration, fn func(Event)) *Rollup {
	if window <= 0 {
		window = 5 * time.Second
	}
	return &Rollup{window: window, flushFn: fn}
}

// Add queues opened/closed ports and (re)starts the flush timer.
func (r *Rollup) Add(opened, closed []scanner.Port) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.opened = append(r.opened, opened...)
	r.closed = append(r.closed, closed...)
	if r.timer == nil {
		r.timer = time.AfterFunc(r.window, r.flush)
	}
}

// Flush forces an immediate emit regardless of the window.
func (r *Rollup) Flush() {
	r.mu.Lock()
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.mu.Unlock()
	r.flush()
}

func (r *Rollup) flush() {
	r.mu.Lock()
	opened := r.opened
	closed := r.closed
	r.opened = nil
	r.closed = nil
	r.timer = nil
	r.mu.Unlock()
	if len(opened) == 0 && len(closed) == 0 {
		return
	}
	r.flushFn(Event{Opened: opened, Closed: closed, At: time.Now().UTC()})
}
