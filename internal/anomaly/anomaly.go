// Package anomaly detects statistical anomalies in port scan results
// by comparing the current open-port count against a rolling baseline.
package anomaly

import (
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"
)

// Level describes how severe an anomaly is.
type Level int

const (
	LevelNone     Level = iota
	LevelMild           // deviation within 1–2 σ
	LevelModerate       // deviation within 2–3 σ
	LevelSevere         // deviation > 3 σ
)

func (l Level) String() string {
	switch l {
	case LevelMild:
		return "mild"
	case LevelModerate:
		return "moderate"
	case LevelSevere:
		return "severe"
	default:
		return "none"
	}
}

// Event is emitted when an anomaly is detected.
type Event struct {
	At        time.Time
	OpenCount int
	Mean      float64
	StdDev    float64
	ZScore    float64
	Level     Level
}

func (e Event) String() string {
	return fmt.Sprintf("anomaly/%s: open=%d mean=%.2f stddev=%.2f z=%.2f",
		e.Level, e.OpenCount, e.Mean, e.StdDev, e.ZScore)
}

// Detector accumulates open-port counts and evaluates each new sample
// against the running mean and standard deviation.
type Detector struct {
	mu      sync.Mutex
	samples []float64
	maxSamples int
	fallback io.Writer
}

// New returns a Detector that keeps up to maxSamples historical counts.
// A minimum of 10 samples is required before anomalies are reported.
func New(maxSamples int, fallback io.Writer) *Detector {
	if maxSamples < 10 {
		maxSamples = 10
	}
	if fallback == nil {
		fallback = os.Stderr
	}
	return &Detector{maxSamples: maxSamples, fallback: fallback}
}

// Evaluate records openCount and returns an Event if an anomaly is detected.
// Returns nil when the sample window is too small or the deviation is normal.
func (d *Detector) Evaluate(openCount int) *Event {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.samples = append(d.samples, float64(openCount))
	if len(d.samples) > d.maxSamples {
		d.samples = d.samples[len(d.samples)-d.maxSamples:]
	}

	if len(d.samples) < 10 {
		return nil
	}

	mean, stddev := stats(d.samples)
	if stddev == 0 {
		return nil
	}

	z := math.Abs(float64(openCount)-mean) / stddev
	lvl := levelFor(z)
	if lvl == LevelNone {
		return nil
	}

	return &Event{
		At:        time.Now().UTC(),
		OpenCount: openCount,
		Mean:      mean,
		StdDev:    stddev,
		ZScore:    z,
		Level:     lvl,
	}
}

// Reset clears all accumulated samples.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.samples = d.samples[:0]
}

func stats(s []float64) (mean, stddev float64) {
	for _, v := range s {
		mean += v
	}
	mean /= float64(len(s))
	for _, v := range s {
		diff := v - mean
		stddev += diff * diff
	}
	stddev = math.Sqrt(stddev / float64(len(s)))
	return
}

func levelFor(z float64) Level {
	switch {
	case z > 3:
		return LevelSevere
	case z > 2:
		return LevelModerate
	case z > 1:
		return LevelMild
	default:
		return LevelNone
	}
}
