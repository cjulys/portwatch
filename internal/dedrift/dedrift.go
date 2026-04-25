// Package dedrift detects and reports configuration drift between the
// current port scan and a saved baseline, producing a human-readable
// summary of ports that have appeared or disappeared since the baseline
// was last committed.
package dedrift

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// DriftKind describes the direction of a single drift entry.
type DriftKind string

const (
	KindAdded   DriftKind = "added"
	KindRemoved DriftKind = "removed"
)

// Entry is a single port that has drifted from the baseline.
type Entry struct {
	Port      scanner.Port
	Kind      DriftKind
	DetectedAt time.Time
}

// Report holds all drift entries produced by a single Evaluate call.
type Report struct {
	Entries    []Entry
	EvaluatedAt time.Time
}

// HasDrift returns true when at least one entry exists.
func (r Report) HasDrift() bool { return len(r.Entries) > 0 }

// Detector compares a current port list against a baseline snapshot.
type Detector struct {
	w   io.Writer
	now func() time.Time
}

// New creates a Detector that writes fallback messages to w.
// Pass nil to default to os.Stderr.
func New(w io.Writer) *Detector {
	if w == nil {
		w = os.Stderr
	}
	return &Detector{w: w, now: time.Now}
}

// Evaluate compares current against baseline and returns a Report.
// Ports present in current but absent from baseline are KindAdded;
// ports absent from current but present in baseline are KindRemoved.
func (d *Detector) Evaluate(baseline, current []scanner.Port) Report {
	now := d.now()
	report := Report{EvaluatedAt: now}

	baseMap := index(baseline)
	currMap := index(current)

	for key, p := range currMap {
		if _, ok := baseMap[key]; !ok {
			report.Entries = append(report.Entries, Entry{Port: p, Kind: KindAdded, DetectedAt: now})
		}
	}
	for key, p := range baseMap {
		if _, ok := currMap[key]; !ok {
			report.Entries = append(report.Entries, Entry{Port: p, Kind: KindRemoved, DetectedAt: now})
		}
	}
	return report
}

// Summarise writes a human-readable drift summary to w.
func (d *Detector) Summarise(r Report, w io.Writer) {
	if w == nil {
		w = d.w
	}
	if !r.HasDrift() {
		fmt.Fprintln(w, "no drift detected")
		return
	}
	for _, e := range r.Entries {
		fmt.Fprintf(w, "[%s] %s/%d\n", e.Kind, e.Port.Protocol, e.Port.Number)
	}
}

func index(ports []scanner.Port) map[string]scanner.Port {
	m := make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		m[fmt.Sprintf("%s/%d", p.Protocol, p.Number)] = p
	}
	return m
}
