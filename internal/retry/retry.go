// Package retry provides configurable retry logic with backoff for portwatch operations.
package retry

import (
	"context"
	"time"
)

// Policy defines how retries are attempted.
type Policy struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	Factor      float64
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Factor:      2.0,
	}
}

// Retryer executes operations with retry logic.
type Retryer struct {
	policy Policy
	sleep  func(time.Duration)
}

// New creates a Retryer with the given policy.
func New(p Policy) *Retryer {
	return &Retryer{policy: p, sleep: time.Sleep}
}

// Do calls fn up to MaxAttempts times, backing off between failures.
// Returns the last error if all attempts fail.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	delay := r.policy.BaseDelay
	var err error
	for i := 0; i < r.policy.MaxAttempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = fn()
		if err == nil {
			return nil
		}
		if i < r.policy.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * r.policy.Factor)
			if delay > r.policy.MaxDelay {
				delay = r.policy.MaxDelay
			}
		}
	}
	return err
}

// Attempts returns the configured maximum attempts.
func (r *Retryer) Attempts() int { return r.policy.MaxAttempts }
