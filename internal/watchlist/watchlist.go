// Package watchlist maintains a set of ports that the operator explicitly
// wants to track. Any port on the watchlist that is found closed is elevated
// to a higher alert level regardless of the global classifier settings.
package watchlist

import (
	"fmt"
	"sync"
)

// Entry identifies a single watched port.
type Entry struct {
	Port     uint16
	Protocol string // "tcp" or "udp"
}

// Watchlist holds the set of explicitly monitored ports.
type Watchlist struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{entries: make(map[string]Entry)}
}

// Add registers a port+protocol pair. Duplicate adds are silently ignored.
func (w *Watchlist) Add(port uint16, protocol string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	k := key(port, protocol)
	w.entries[k] = Entry{Port: port, Protocol: protocol}
}

// Remove deregisters a port+protocol pair. Removing a non-existent entry is a
// no-op.
func (w *Watchlist) Remove(port uint16, protocol string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entries, key(port, protocol))
}

// Contains reports whether the given port+protocol is on the watchlist.
func (w *Watchlist) Contains(port uint16, protocol string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, ok := w.entries[key(port, protocol)]
	return ok
}

// All returns a snapshot of every entry currently on the watchlist.
func (w *Watchlist) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, e)
	}
	return out
}

// Len returns the number of entries currently on the watchlist.
func (w *Watchlist) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.entries)
}

func key(port uint16, protocol string) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
