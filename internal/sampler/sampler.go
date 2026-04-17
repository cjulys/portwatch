// Package sampler provides adaptive scan interval adjustment based on
// recent port change activity. When changes are frequent the interval
// shrinks; when the environment is quiet it grows back toward the
// configured maximum.
package sampler

import (
	"time"
)

// Sampler adjusts a polling interval adaptively.
type Sampler struct {
	min     time.Duration
	max     time.Duration
	current time.Duration
	// stepDown is the factor applied when activity is detected (< 1).
	stepDown float64
	// stepUp is the factor applied when no activity is detected (> 1).
	stepUp float64
}

// New returns a Sampler whose interval starts at max and can shrink to min.
// stepDown must be in (0,1) and stepUp must be > 1; invalid values are clamped.
func New(min, max time.Duration, stepDown, stepUp float64) *Sampler {
	if stepDown <= 0 || stepDown >= 1 {
		stepDown = 0.5
	}
	if stepUp <= 1 {
		stepUp = 1.5
	}
	if min <= 0 {
		min = time.Second
	}
	if max < min {
		max = min
	}
	return &Sampler{
		min:     min,
		max:     max,
		current: max,
		stepDown: stepDown,
		stepUp:  stepUp,
	}
}

// Current returns the current recommended interval.
func (s *Sampler) Current() time.Duration {
	return s.current
}

// RecordActivity signals that port changes were detected; the interval
// is decreased toward the minimum.
func (s *Sampler) RecordActivity() {
	next := time.Duration(float64(s.current) * s.stepDown)
	if next < s.min {
		next = s.min
	}
	s.current = next
}

// RecordQuiet signals that no changes were detected; the interval
// grows back toward the maximum.
func (s *Sampler) RecordQuiet() {
	next := time.Duration(float64(s.current) * s.stepUp)
	if next > s.max {
		next = s.max
	}
	s.current = next
}

// Reset restores the interval to its initial maximum value.
func (s *Sampler) Reset() {
	s.current = s.max
}
