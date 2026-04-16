// Package debounce delays forwarding of repeated diff events within a quiet
// window, reducing noise when ports flap rapidly.
package debounce

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Handler is called with the accumulated diffs after the quiet window.
type Handler func(diffs []scanner.Diff)

// Debouncer buffers diffs and fires the handler after a quiet period.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	buf     []scanner.Diff
	timer   *time.Timer
	handler Handler
}

// New creates a Debouncer with the given quiet window and downstream handler.
func New(window time.Duration, h Handler) *Debouncer {
	return &Debouncer{window: window, handler: h}
}

// Add enqueues diffs and (re)starts the quiet-window timer.
func (d *Debouncer) Add(diffs []scanner.Diff) {
	if len(diffs) == 0 {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.buf = append(d.buf, diffs...)
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.window, d.flush)
}

// Flush forces immediate delivery of any buffered diffs.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
	d.mu.Unlock()
	d.flush()
}

func (d *Debouncer) flush() {
	d.mu.Lock()
	if len(d.buf) == 0 {
		d.mu.Unlock()
		return
	}
	out := make([]scanner.Diff, len(d.buf))
	copy(out, d.buf)
	d.buf = d.buf[:0]
	d.mu.Unlock()
	d.handler(out)
}
