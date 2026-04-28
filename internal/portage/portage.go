// Package portage tracks how long each port has been continuously open,
// providing an "age" metric that can surface long-lived unexpected listeners.
package portage

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records when a port was first seen open in the current continuous run.
type Entry struct {
	Port      scanner.Port
	FirstSeen time.Time
}

// Age returns how long the port has been continuously open relative to now.
func (e Entry) Age(now time.Time) time.Duration {
	return now.Sub(e.FirstSeen)
}

// Tracker maintains first-seen timestamps for open ports.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]Entry
	clock   func() time.Time
}

// New returns a Tracker using the real wall clock.
func New() *Tracker {
	return &Tracker{
		entries: make(map[string]Entry),
		clock:   time.Now,
	}
}

func portKey(p scanner.Port) string {
	return p.Protocol + ":" + p.Address + ":" + itoa(p.Number)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

// Update reconciles the tracker state against the current open port list.
// Ports no longer present are removed; newly seen ports are recorded.
func (t *Tracker) Update(current []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.clock()

	// Build lookup of current ports.
	seen := make(map[string]scanner.Port, len(current))
	for _, p := range current {
		seen[portKey(p)] = p
	}

	// Remove entries no longer open.
	for k := range t.entries {
		if _, ok := seen[k]; !ok {
			delete(t.entries, k)
		}
	}

	// Add entries for newly seen ports.
	for k, p := range seen {
		if _, exists := t.entries[k]; !exists {
			t.entries[k] = Entry{Port: p, FirstSeen: now}
		}
	}
}

// Get returns the Entry for a port and whether it is being tracked.
func (t *Tracker) Get(p scanner.Port) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[portKey(p)]
	return e, ok
}

// Snapshot returns a copy of all current entries.
func (t *Tracker) Snapshot() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of tracked ports.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}
