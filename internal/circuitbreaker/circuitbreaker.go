// Package circuitbreaker implements a simple circuit breaker for external
// calls such as webhook delivery or SMTP alerting.
package circuitbreaker

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // blocking calls
	StateHalfOpen              // probing
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	}
	return "unknown"
}

// Breaker is a circuit breaker instance.
type Breaker struct {
	mu          sync.Mutex
	state       State
	failures    int
	threshold   int
	resetAfter  time.Duration
	openedAt    time.Time
	fallback    io.Writer
	now         func() time.Time
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts reset after resetAfter duration.
func New(threshold int, resetAfter time.Duration) *Breaker {
	return &Breaker{
		threshold:  threshold,
		resetAfter: resetAfter,
		fallback:   os.Stderr,
		now:        time.Now,
	}
}

// Allow reports whether the call should be attempted.
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if b.now().Sub(b.openedAt) >= b.resetAfter {
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess resets the breaker to closed.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failure and may open the circuit.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		if b.state != StateOpen {
			fmt.Fprintf(b.fallback, "portwatch: circuit breaker opened after %d failures\n", b.failures)
			b.openedAt = b.now()
		}
		b.state = StateOpen
	}
}

// State returns the current state.
func (b *Breaker) CurrentState() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
