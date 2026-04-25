// Package drift tracks how far the current port state has deviated from a
// known-good baseline and exposes a numeric drift score.
package drift

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Score holds the result of a drift evaluation.
type Score struct {
	Added   int
	Removed int
	Total   int
	At      time.Time
}

// Tracker computes drift between the current port snapshot and a baseline.
type Tracker struct {
	mu       sync.Mutex
	baseline []scanner.Port
	last     Score
}

// New returns a Tracker seeded with the given baseline ports.
func New(baseline []scanner.Port) *Tracker {
	b := make([]scanner.Port, len(baseline))
	copy(b, baseline)
	return &Tracker{baseline: b}
}

// SetBaseline replaces the current baseline.
func (t *Tracker) SetBaseline(ports []scanner.Port) {
	t.mu.Lock()
	defer t.mu.Unlock()
	b := make([]scanner.Port, len(ports))
	copy(b, ports)
	t.baseline = b
}

// Evaluate computes a Score by comparing current against the stored baseline.
func (t *Tracker) Evaluate(current []scanner.Port) Score {
	t.mu.Lock()
	defer t.mu.Unlock()

	baseSet := toSet(t.baseline)
	currSet := toSet(current)

	var added, removed int
	for k := range currSet {
		if _, ok := baseSet[k]; !ok {
			added++
		}
	}
	for k := range baseSet {
		if _, ok := currSet[k]; !ok {
			removed++
		}
	}

	s := Score{
		Added:   added,
		Removed: removed,
		Total:   added + removed,
		At:      time.Now().UTC(),
	}
	t.last = s
	return s
}

// Last returns the most recent Score without re-evaluating.
func (t *Tracker) Last() Score {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.last
}

func toSet(ports []scanner.Port) map[string]struct{} {
	m := make(map[string]struct{}, len(ports))
	for _, p := range ports {
		m[p.Protocol+":"+itoa(p.Port)] = struct{}{}
	}
	return m
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [10]byte{}
	pos := 10
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
