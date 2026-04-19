// Package escalate promotes alert severity when the same port change
// is observed repeatedly within a configurable window.
package escalate

import (
	"sync"
	"time"
)

// Level represents an alert severity.
type Level int

const (
	Info Level = iota
	Warning
	Critical
)

func (l Level) String() string {
	switch l {
	case Warning:
		return "warning"
	case Critical:
		return "critical"
	default:
		return "info"
	}
}

type entry struct {
	count int
	first time.Time
}

// Escalator tracks repeated events per key and promotes their level.
type Escalator struct {
	mu      sync.Mutex
	entries map[string]*entry
	window  time.Duration
	warnAt  int
	critAt  int
	now     func() time.Time
}

// New returns an Escalator. warnAt and critAt are the repeat thresholds.
func New(window time.Duration, warnAt, critAt int) *Escalator {
	return &Escalator{
		entries: make(map[string]*entry),
		window:  window,
		warnAt:  warnAt,
		critAt:  critAt,
		now:     time.Now,
	}
}

// Evaluate records an occurrence for key and returns the current level.
func (e *Escalator) Evaluate(key string) Level {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.now()
	ent, ok := e.entries[key]
	if !ok || now.Sub(ent.first) > e.window {
		ent = &entry{first: now}
		e.entries[key] = ent
	}
	ent.count++
	switch {
	case ent.count >= e.critAt:
		return Critical
	case ent.count >= e.warnAt:
		return Warning
	default:
		return Info
	}
}

// Reset clears the history for a key.
func (e *Escalator) Reset(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.entries, key)
}

// Flush clears all tracked state.
func (e *Escalator) Flush() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.entries = make(map[string]*entry)
}
