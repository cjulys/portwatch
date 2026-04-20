// Package window provides a sliding time-window counter for tracking
// event frequency over a rolling duration.
package window

import (
	"sync"
	"time"
)

// entry holds a single timestamped count.
type entry struct {
	at    time.Time
	count int
}

// Window is a thread-safe sliding time-window counter.
type Window struct {
	mu       sync.Mutex
	span     time.Duration
	buckets  []entry
	nowFn    func() time.Time
}

// New returns a Window that tracks events within the given span.
func New(span time.Duration) *Window {
	return &Window{span: span, nowFn: time.Now}
}

// newWithClock is used in tests to inject a clock.
func newWithClock(span time.Duration, nowFn func() time.Time) *Window {
	return &Window{span: span, nowFn: nowFn}
}

// Add records n occurrences at the current time.
func (w *Window) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = append(w.buckets, entry{at: w.nowFn(), count: n})
	w.evict()
}

// Count returns the total events recorded within the window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	total := 0
	for _, b := range w.buckets {
		total += b.count
	}
	return total
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.buckets = nil
}

// evict removes entries older than the window span. Must be called with mu held.
func (w *Window) evict() {
	cutoff := w.nowFn().Add(-w.span)
	i := 0
	for i < len(w.buckets) && w.buckets[i].at.Before(cutoff) {
		i++
	}
	w.buckets = w.buckets[i:]
}
