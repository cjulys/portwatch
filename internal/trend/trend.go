// Package trend analyses port scan history to detect rising or falling
// open-port counts over a sliding window.
package trend

import (
	"time"

	"portwatch/internal/history"
)

// Direction describes whether open ports are increasing, decreasing, or stable.
type Direction string

const (
	Rising  Direction = "rising"
	Falling Direction = "falling"
	Stable  Direction = "stable"
)

// Result holds the outcome of a trend analysis.
type Result struct {
	Direction Direction
	Delta     int // net change in open-port count across the window
	Window    time.Duration
	Samples   int
}

// Analyzer computes port-count trends from a history store.
type Analyzer struct {
	h      *history.History
	window time.Duration
}

// New returns an Analyzer that looks back window duration into h.
func New(h *history.History, window time.Duration) *Analyzer {
	if window <= 0 {
		window = 5 * time.Minute
	}
	return &Analyzer{h: h, window: window}
}

// Analyze returns a Result describing the trend over the configured window.
func (a *Analyzer) Analyze(now time.Time) Result {
	entries := a.h.Since(now.Add(-a.window))
	if len(entries) == 0 {
		return Result{Direction: Stable, Window: a.window}
	}

	first := len(entries[0].Ports)
	last := len(entries[len(entries)-1].Ports)
	delta := last - first

	dir := Stable
	switch {
	case delta > 0:
		dir = Rising
	case delta < 0:
		dir = Falling
	}

	return Result{
		Direction: dir,
		Delta:     delta,
		Window:    a.window,
		Samples:   len(entries),
	}
}
