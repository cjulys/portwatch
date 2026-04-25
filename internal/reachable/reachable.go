// Package reachable provides port reachability scoring based on recent
// scan history. A port's reachability score reflects how consistently it
// has been observed open over a sliding time window.
package reachable

import (
	"sync"
	"time"
)

// Score holds the reachability result for a single port key.
type Score struct {
	// Key identifies the port (e.g. "tcp:80").
	Key string
	// Seen is the number of scans in which the port was observed open.
	Seen int
	// Total is the total number of scans recorded in the window.
	Total int
	// Ratio is Seen/Total, or 0 when Total is 0.
	Ratio float64
}

type entry struct {
	at   time.Time
	open bool
}

// Tracker accumulates per-port observations and computes scores.
type Tracker struct {
	mu     sync.Mutex
	window time.Duration
	data   map[string][]entry
	now    func() time.Time
}

// New returns a Tracker that considers observations within window.
func New(window time.Duration) *Tracker {
	return &Tracker{
		window: window,
		data:   make(map[string][]entry),
		now:    time.Now,
	}
}

// Record registers whether the port identified by key was open at the
// current moment.
func (t *Tracker) Record(key string, open bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	t.data[key] = t.evict(t.data[key], now)
	t.data[key] = append(t.data[key], entry{at: now, open: open})
}

// Get returns the current Score for key.
func (t *Tracker) Get(key string) Score {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	entries := t.evict(t.data[key], now)
	t.data[key] = entries
	s := Score{Key: key, Total: len(entries)}
	for _, e := range entries {
		if e.open {
			s.Seen++
		}
	}
	if s.Total > 0 {
		s.Ratio = float64(s.Seen) / float64(s.Total)
	}
	return s
}

// Flush removes all tracked data.
func (t *Tracker) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.data = make(map[string][]entry)
}

func (t *Tracker) evict(entries []entry, now time.Time) []entry {
	cutoff := now.Add(-t.window)
	i := 0
	for i < len(entries) && entries[i].at.Before(cutoff) {
		i++
	}
	return entries[i:]
}
