// Package tap provides a passive traffic tap that mirrors scan diffs to
// registered observers without blocking the main processing pipeline.
package tap

import (
	"io"
	"os"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Sink is a function that receives a copy of diffs as they flow through the tap.
type Sink func(diffs []scanner.Diff)

// Tap mirrors diffs to zero or more sinks without modifying them.
type Tap struct {
	mu      sync.RWMutex
	sinks   []Sink
	fallback io.Writer
}

// New returns a Tap that writes error notices to fallback (defaults to stderr).
func New(fallback io.Writer) *Tap {
	if fallback == nil {
		fallback = os.Stderr
	}
	return &Tap{fallback: fallback}
}

// Register adds a sink to the tap. Sinks are called in registration order.
func (t *Tap) Register(s Sink) {
	if s == nil {
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.sinks = append(t.sinks, s)
}

// Len returns the number of registered sinks.
func (t *Tap) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.sinks)
}

// Send copies diffs to every registered sink. Panics inside a sink are
// recovered and written to the fallback writer so other sinks still run.
func (t *Tap) Send(diffs []scanner.Diff) {
	if len(diffs) == 0 {
		return
	}
	t.mu.RLock()
	snap := make([]Sink, len(t.sinks))
	copy(snap, t.sinks)
	t.mu.RUnlock()

	for _, s := range snap {
		func(fn Sink) {
			defer func() {
				if r := recover(); r != nil {
					io.WriteString(t.fallback, "tap: sink panic recovered\n")
				}
			}()
			fn(diffs)
		}(s)
	}
}
