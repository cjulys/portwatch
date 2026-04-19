// Package jitter adds randomised offsets to scan intervals to avoid
// thundering-herd effects when multiple portwatch instances run in parallel.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Jitter applies a configurable random offset to a base duration.
type Jitter struct {
	mu      sync.Mutex
	rng     *rand.Rand
	factor  float64 // fraction of base to use as max jitter, e.g. 0.2 = ±20%
}

// New returns a Jitter with the given factor (0.0–1.0).
// Values outside that range are clamped.
func New(factor float64) *Jitter {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	return &Jitter{
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
		factor: factor,
	}
}

// Apply returns base ± (factor * base * random[0,1)).
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if j.factor == 0 || base <= 0 {
		return base
	}
	j.mu.Lock()
	r := j.rng.Float64()
	j.mu.Unlock()

	offset := time.Duration(float64(base) * j.factor * r)
	if r < 0.5 {
		return base - offset
	}
	return base + offset
}

// Factor returns the configured jitter factor.
func (j *Jitter) Factor() float64 {
	return j.factor
}
