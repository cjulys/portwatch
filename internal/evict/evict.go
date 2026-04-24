// Package evict provides a port-keyed LRU eviction tracker that removes
// stale entries from in-memory caches when they have not been seen within
// a configurable time-to-live window.
package evict

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// entry holds the last-seen timestamp for a single tracked key.
type entry struct {
	lastSeen time.Time
}

// Tracker records the last observation time for each port key and reports
// which keys have exceeded their TTL.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]entry
	ttl     time.Duration
	now     func() time.Time
}

// New returns a Tracker with the given time-to-live.
func New(ttl time.Duration) *Tracker {
	return &Tracker{
		entries: make(map[string]entry),
		ttl:     ttl,
		now:     time.Now,
	}
}

// Touch records that the given port was observed right now.
func (t *Tracker) Touch(p scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[key(p)] = entry{lastSeen: t.now()}
}

// Evict removes all entries whose last-seen time is older than the TTL and
// returns the evicted ports.
func (t *Tracker) Evict() []scanner.Port {
	t.mu.Lock()
	defer t.mu.Unlock()

	cutoff := t.now().Add(-t.ttl)
	var evicted []scanner.Port
	for k, e := range t.entries {
		if e.lastSeen.Before(cutoff) {
			evicted = append(evicted, portFromKey(k))
			delete(t.entries, k)
		}
	}
	return evicted
}

// Len returns the number of currently tracked entries.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}

func key(p scanner.Port) string {
	return p.Proto + ":" + itoa(p.Number)
}

func portFromKey(k string) scanner.Port {
	for i := 0; i < len(k); i++ {
		if k[i] == ':' {
			n := 0
			for _, c := range k[i+1:] {
				n = n*10 + int(c-'0')
			}
			return scanner.Port{Proto: k[:i], Number: n}
		}
	}
	return scanner.Port{}
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
