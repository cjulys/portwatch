// Package dedup provides event deduplication based on port-state fingerprints.
// It suppresses repeated identical diff events within a configurable window,
// preventing alert storms when the same port repeatedly flaps.
package dedup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records the last time a particular diff key was seen.
type Entry struct {
	LastSeen time.Time
	Count    int
}

// Deduplicator suppresses duplicate diff events within a time window.
type Deduplicator struct {
	mu      sync.Mutex
	window  time.Duration
	entries map[string]*Entry
	now     func() time.Time
}

// New returns a Deduplicator with the given deduplication window.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		window:  window,
		entries: make(map[string]*Entry),
		now:     time.Now,
	}
}

// IsDuplicate returns true if an identical event for the given port and
// direction was already seen within the deduplication window.
func (d *Deduplicator) IsDuplicate(p scanner.Port, direction string) bool {
	key := diffKey(p, direction)
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if e, ok := d.entries[key]; ok {
		if now.Sub(e.LastSeen) < d.window {
			e.Count++
			return true
		}
	}
	d.entries[key] = &Entry{LastSeen: now, Count: 1}
	return false
}

// Flush removes all tracked entries, resetting deduplication state.
func (d *Deduplicator) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries = make(map[string]*Entry)
}

// Stats returns the current count for a given port and direction, or 0.
func (d *Deduplicator) Stats(p scanner.Port, direction string) int {
	key := diffKey(p, direction)
	d.mu.Lock()
	defer d.mu.Unlock()
	if e, ok := d.entries[key]; ok {
		return e.Count
	}
	return 0
}

func diffKey(p scanner.Port, direction string) string {
	return p.Proto + ":" + itoa(p.Port) + ":" + direction
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
