// Package summary produces periodic digest reports of port activity.
package summary

import (
	"fmt"
	"io"
	"os"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/scanner"
)

// Report holds aggregated port activity for a time window.
type Report struct {
	GeneratedAt time.Time
	Window      time.Duration
	OpenPorts   []scanner.Port
	NewPorts    []scanner.Port
	ClosedPorts []scanner.Port
}

// Builder builds summary reports from history.
type Builder struct {
	h   *history.History
	out io.Writer
}

// New returns a Builder that writes to w. If w is nil, os.Stdout is used.
func New(h *history.History, w io.Writer) *Builder {
	if w == nil {
		w = os.Stdout
	}
	return &Builder{h: h, out: w}
}

// Build produces a Report covering the last window duration.
func (b *Builder) Build(window time.Duration) Report {
	since := time.Now().UTC().Add(-window)
	entries := b.h.Since(since)

	seen := map[string]scanner.Port{}
	newPorts := []scanner.Port{}
	closedPorts := []scanner.Port{}

	for _, e := range entries {
		for _, p := range e.Opened {
			key := fmt.Sprintf("%s/%d", p.Protocol, p.Port)
			seen[key] = p
			newPorts = append(newPorts, p)
		}
		for _, p := range e.Closed {
			closedPorts = append(closedPorts, p)
		}
	}

	open := make([]scanner.Port, 0, len(seen))
	for _, p := range seen {
		open = append(open, p)
	}

	return Report{
		GeneratedAt: time.Now().UTC(),
		Window:      window,
		OpenPorts:   open,
		NewPorts:    newPorts,
		ClosedPorts: closedPorts,
	}
}

// Print writes a human-readable summary to the builder's writer.
func (b *Builder) Print(r Report) {
	fmt.Fprintf(b.out, "=== Port Summary (last %s) ===\n", r.Window)
	fmt.Fprintf(b.out, "Generated: %s\n", r.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(b.out, "New ports opened : %d\n", len(r.NewPorts))
	fmt.Fprintf(b.out, "Ports closed     : %d\n", len(r.ClosedPorts))
	fmt.Fprintf(b.out, "Distinct open    : %d\n", len(r.OpenPorts))
}
