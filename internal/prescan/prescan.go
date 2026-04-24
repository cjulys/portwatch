// Package prescan performs a lightweight pre-scan check to determine
// whether a full port scan is warranted based on recent activity and
// system load indicators.
package prescan

import (
	"context"
	"io"
	"os"
	"sync"
	"time"
)

// Result summarises the outcome of a pre-scan evaluation.
type Result struct {
	ShouldScan bool
	Reason     string
	EvaluatedAt time.Time
}

// Checker decides whether a full scan should proceed.
type Checker struct {
	mu          sync.Mutex
	lastScan    time.Time
	minInterval time.Duration
	fallback    io.Writer
}

// New returns a Checker that suppresses scans occurring faster than
// minInterval. A zero minInterval disables the guard.
func New(minInterval time.Duration) *Checker {
	return &Checker{
		minInterval: minInterval,
		fallback:    os.Stderr,
	}
}

// WithWriter replaces the fallback writer used for diagnostic output.
func (c *Checker) WithWriter(w io.Writer) *Checker {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.fallback = w
	return c
}

// Evaluate returns a Result indicating whether a full scan should run.
// It respects ctx cancellation; a cancelled context always returns
// ShouldScan=false.
func (c *Checker) Evaluate(ctx context.Context) Result {
	if err := ctx.Err(); err != nil {
		return Result{
			ShouldScan:  false,
			Reason:      "context cancelled",
			EvaluatedAt: time.Now().UTC(),
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UTC()

	if c.minInterval > 0 && !c.lastScan.IsZero() {
		elapsed := now.Sub(c.lastScan)
		if elapsed < c.minInterval {
			return Result{
				ShouldScan:  false,
				Reason:      "min interval not elapsed",
				EvaluatedAt: now,
			}
		}
	}

	c.lastScan = now
	return Result{
		ShouldScan:  true,
		Reason:      "ok",
		EvaluatedAt: now,
	}
}

// Reset clears the last-scan timestamp, allowing the next Evaluate call
// to proceed regardless of the configured interval.
func (c *Checker) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastScan = time.Time{}
}
