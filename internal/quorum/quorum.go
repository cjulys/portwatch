// Package quorum requires a port change to be observed across multiple
// consecutive scans before it is forwarded. This prevents single-scan
// noise (e.g. a port that appears open for exactly one poll cycle) from
// triggering downstream alerts.
package quorum

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry tracks how many consecutive scans have confirmed a diff direction
// for a given port key.
type Entry struct {
	Count     int
	FirstSeen time.Time
}

// Quorum accumulates observations and emits a diff only once the required
// confirmation count is reached.
type Quorum struct {
	mu        sync.Mutex
	threshold int
	counts    map[string]*Entry
	clock     func() time.Time
}

// New returns a Quorum that requires threshold consecutive confirmations.
// threshold must be >= 1; values below 1 are clamped to 1.
func New(threshold int) *Quorum {
	if threshold < 1 {
		threshold = 1
	}
	return &Quorum{
		threshold: threshold,
		counts:    make(map[string]*Entry),
		clock:     time.Now,
	}
}

// Observe records one observation of the given diff direction for a port.
// It returns true when the confirmation threshold has just been reached,
// meaning the caller should act on the change. Subsequent calls for the
// same key continue returning false until Reset is called.
func (q *Quorum) Observe(p scanner.Port, direction string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	key := fmt.Sprintf("%s:%d:%s", p.Proto, p.Port, direction)
	e, ok := q.counts[key]
	if !ok {
		e = &Entry{FirstSeen: q.clock()}
		q.counts[key] = e
	}
	e.Count++
	return e.Count == q.threshold
}

// Reset removes the accumulated count for a port/direction pair, allowing
// it to be re-evaluated from scratch.
func (q *Quorum) Reset(p scanner.Port, direction string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	key := fmt.Sprintf("%s:%d:%s", p.Proto, p.Port, direction)
	delete(q.counts, key)
}

// Flush clears all accumulated state.
func (q *Quorum) Flush() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.counts = make(map[string]*Entry)
}
