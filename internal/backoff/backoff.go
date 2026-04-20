// Package backoff provides exponential back-off with optional jitter for
// repeated operations that may transiently fail.
package backoff

import (
	"context"
	"math"
	"time"
)

// Policy holds the parameters for an exponential back-off sequence.
type Policy struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	MaxAttempts     int
}

// DefaultPolicy returns a sensible default back-off policy.
func DefaultPolicy() Policy {
	return Policy{
		InitialInterval: 500 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		MaxAttempts:     8,
	}
}

// Backoff holds state for a single back-off sequence.
type Backoff struct {
	policy  Policy
	attempt int
	clock   func() time.Time
	sleep   func(context.Context, time.Duration) error
}

// New creates a Backoff using the given policy.
func New(p Policy) *Backoff {
	return &Backoff{
		policy: p,
		clock:  time.Now,
		sleep: func(ctx context.Context, d time.Duration) error {
			select {
			case <-time.After(d):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}
}

// Next waits for the next back-off interval and returns false when attempts are
// exhausted or the context is cancelled.
func (b *Backoff) Next(ctx context.Context) bool {
	if b.attempt >= b.policy.MaxAttempts {
		return false
	}
	if b.attempt > 0 {
		d := b.interval()
		if err := b.sleep(ctx, d); err != nil {
			return false
		}
	}
	b.attempt++
	return true
}

// Attempt returns the current attempt number (1-based after first Next call).
func (b *Backoff) Attempt() int { return b.attempt }

// Reset restarts the back-off sequence from the beginning.
func (b *Backoff) Reset() { b.attempt = 0 }

func (b *Backoff) interval() time.Duration {
	v := float64(b.policy.InitialInterval) * math.Pow(b.policy.Multiplier, float64(b.attempt-1))
	if v > float64(b.policy.MaxInterval) {
		v = float64(b.policy.MaxInterval)
	}
	return time.Duration(v)
}
