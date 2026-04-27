// Package census tracks port population statistics across successive scans,
// providing counts and deltas that feed into dashboards and alerting pipelines.
package census

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds the port counts captured at a single point in time.
type Snapshot struct {
	At       time.Time
	Total    int
	TCP      int
	UDP      int
}

// Delta describes how counts changed between two consecutive snapshots.
type Delta struct {
	Total int
	TCP   int
	UDP   int
}

// Census accumulates per-scan population data.
type Census struct {
	mu   sync.Mutex
	last *Snapshot
}

// New returns an initialised Census.
func New() *Census {
	return &Census{}
}

// Record ingests a fresh port list and returns the resulting Snapshot and the
// Delta relative to the previous scan. If this is the first scan the Delta
// fields will equal the Snapshot counts.
func (c *Census) Record(ports []scanner.Port) (Snapshot, Delta) {
	snap := Snapshot{
		At:    time.Now().UTC(),
		Total: len(ports),
	}
	for _, p := range ports {
		switch p.Protocol {
		case "tcp":
			snap.TCP++
		case "udp":
			snap.UDP++
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	var d Delta
	if c.last == nil {
		d = Delta{Total: snap.Total, TCP: snap.TCP, UDP: snap.UDP}
	} else {
		d = Delta{
			Total: snap.Total - c.last.Total,
			TCP:   snap.TCP - c.last.TCP,
			UDP:   snap.UDP - c.last.UDP,
		}
	}

	copy := snap
	c.last = &copy
	return snap, d
}

// Last returns the most recent Snapshot, or nil if no scan has been recorded.
func (c *Census) Last() *Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.last == nil {
		return nil
	}
	copy := *c.last
	return &copy
}

// Reset clears accumulated state so the next Record call behaves as if it were
// the first.
func (c *Census) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.last = nil
}
