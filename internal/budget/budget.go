// Package budget enforces a scan-rate budget, capping how many scans
// may be initiated within a rolling time window across concurrent callers.
package budget

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Budget tracks remaining scan capacity within a rolling window.
type Budget struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	timestamps []time.Time
	clock    func() time.Time
	fallback io.Writer
}

// New creates a Budget that allows at most max scans per window duration.
// A zero or negative max disables limiting (all calls are allowed).
func New(max int, window time.Duration) *Budget {
	return &Budget{
		max:      max,
		window:   window,
		clock:    time.Now,
		fallback: os.Stderr,
	}
}

// Allow reports whether a scan may proceed. It records the attempt if allowed.
// The provided context is checked before acquiring the internal lock.
func (b *Budget) Allow(ctx context.Context) bool {
	if ctx.Err() != nil {
		return false
	}
	if b.max <= 0 {
		return true
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.clock()
	b.evict(now)
	if len(b.timestamps) >= b.max {
		fmt.Fprintf(b.fallback, "portwatch/budget: scan budget exhausted (%d/%d within %s)\n",
			len(b.timestamps), b.max, b.window)
		return false
	}
	b.timestamps = append(b.timestamps, now)
	return true
}

// Remaining returns how many more scans may be initiated in the current window.
func (b *Budget) Remaining() int {
	if b.max <= 0 {
		return -1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict(b.clock())
	r := b.max - len(b.timestamps)
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded timestamps, fully restoring the budget.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.timestamps = b.timestamps[:0]
}

// evict removes timestamps that have fallen outside the rolling window.
// Must be called with b.mu held.
func (b *Budget) evict(now time.Time) {
	cutoff := now.Add(-b.window)
	i := 0
	for i < len(b.timestamps) && b.timestamps[i].Before(cutoff) {
		i++
	}
	b.timestamps = b.timestamps[i:]
}
