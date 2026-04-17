// Package decay tracks alert frequency and reduces severity over time
// when a port state has been stable for a configured period.
package decay

import (
	"sync"
	"time"
)

// Level represents a decayed severity level.
type Level int

const (
	LevelCritical Level = iota
	LevelWarning
	LevelInfo
	LevelSilenced
)

func (l Level) String() string {
	switch l {
	case LevelCritical:
		return "critical"
	case LevelWarning:
		return "warning"
	case LevelInfo:
		return "info"
	case LevelSilenced:
		return "silenced"
	}
	return "unknown"
}

type entry struct {
	first time.Time
	count int
}

// Tracker reduces alert severity the longer a key remains unchanged.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*entry
	steps   []time.Duration // thresholds to step down severity
	clock   func() time.Time
}

// New creates a Tracker. steps defines durations after which severity
// decreases one level (e.g. 1m, 10m, 1h → critical→warning→info→silenced).
func New(steps []time.Duration) *Tracker {
	return &Tracker{
		entries: make(map[string]*entry),
		steps:   steps,
		clock:   time.Now,
	}
}

// Evaluate returns the current Level for key, registering it on first call.
func (t *Tracker) Evaluate(key string) Level {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	e, ok := t.entries[key]
	if !ok {
		t.entries[key] = &entry{first: now, count: 1}
		return LevelCritical
	}
	e.count++
	elapsed := now.Sub(e.first)
	level := Level(0)
	for _, step := range t.steps {
		if elapsed >= step {
			level++
		}
	}
	if int(level) > int(LevelSilenced) {
		level = LevelSilenced
	}
	return level
}

// Reset removes a key so the next evaluation starts fresh.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Flush removes all tracked keys.
func (t *Tracker) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]*entry)
}
