// Package heartbeat emits periodic alive signals so external monitors
// can detect if the portwatch daemon has silently stopped running.
package heartbeat

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// Heartbeat writes a timestamped alive line to a writer on every tick.
type Heartbeat struct {
	interval time.Duration
	w        io.Writer
}

// New returns a Heartbeat that ticks at the given interval.
// If w is nil, os.Stderr is used.
func New(interval time.Duration, w io.Writer) *Heartbeat {
	if w == nil {
		w = os.Stderr
	}
	if interval <= 0 {
		interval = 60 * time.Second
	}
	return &Heartbeat{interval: interval, w: w}
}

// Run blocks until ctx is cancelled, writing a heartbeat line on each tick.
func (h *Heartbeat) Run(ctx context.Context) {
	ticker := time.NewTicker(h.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			fmt.Fprintf(h.w, "[heartbeat] alive at %s\n", t.UTC().Format(time.RFC3339))
		}
	}
}

// Interval returns the configured tick interval.
func (h *Heartbeat) Interval() time.Duration { return h.interval }
