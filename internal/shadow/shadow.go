// Package shadow maintains a parallel "shadow" scan result that lags behind
// the live scan by one cycle. Comparing live vs shadow lets callers detect
// ports that appear and disappear within a single interval (flapping).
package shadow

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds one generation of scan results.
type Entry struct {
	Ports    []scanner.Port
	RecordedAt time.Time
}

// Shadow stores the previous scan generation.
type Shadow struct {
	mu      sync.RWMutex
	current *Entry
	prev    *Entry
}

// New returns an empty Shadow.
func New() *Shadow {
	return &Shadow{}
}

// Commit advances the shadow: the current generation becomes previous,
// and ports becomes the new current generation.
func (s *Shadow) Commit(ports []scanner.Port) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prev = s.current
	s.current = &Entry{
		Ports:      ports,
		RecordedAt: time.Now().UTC(),
	}
}

// Current returns the most recently committed entry, or nil if none.
func (s *Shadow) Current() *Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.current
}

// Previous returns the entry from one cycle ago, or nil if fewer than two
// commits have been made.
func (s *Shadow) Previous() *Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.prev
}

// Flapping returns ports that were present in the previous generation but
// absent in the current generation, AND ports present in current but absent
// in previous — i.e. ports whose state changed in consecutive scans.
func (s *Shadow) Flapping() []scanner.Port {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.prev == nil || s.current == nil {
		return nil
	}
	prevMap := toMap(s.prev.Ports)
	currMap := toMap(s.current.Ports)
	var out []scanner.Port
	for k, p := range prevMap {
		if _, ok := currMap[k]; !ok {
			out = append(out, p)
		}
	}
	for k, p := range currMap {
		if _, ok := prevMap[k]; !ok {
			out = append(out, p)
		}
	}
	return out
}

func toMap(ports []scanner.Port) map[string]scanner.Port {
	m := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		m[key(p)] = p
	}
	return m
}

func key(p scanner.Port) string {
	return p.Address + "/" + p.Protocol + "/" + itoa(p.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 6)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
