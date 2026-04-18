// Package grace provides a graceful shutdown coordinator that waits for
// in-flight scan cycles to complete before exiting.
package grace

import (
	"context"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Coordinator listens for OS signals and waits for active work to drain
// before cancelling the root context.
type Coordinator struct {
	mu      sync.Mutex
	active  int
	drain   chan struct{}
	timeout time.Duration
	w       io.Writer
}

// New returns a Coordinator with the given drain timeout.
func New(timeout time.Duration, w io.Writer) *Coordinator {
	if w == nil {
		w = os.Stderr
	}
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Coordinator{drain: make(chan struct{}, 1), timeout: timeout, w: w}
}

// Acquire marks one unit of work as in-flight.
func (c *Coordinator) Acquire() {
	c.mu.Lock()
	c.active++
	c.mu.Unlock()
}

// Release marks one unit of work as complete and signals if drained.
func (c *Coordinator) Release() {
	c.mu.Lock()
	if c.active > 0 {
		c.active--
	}
	if c.active == 0 {
		select {
		case c.drain <- struct{}{}:
		default:
		}
	}
	c.mu.Unlock()
}

// Wait blocks until SIGINT/SIGTERM is received, then drains or times out,
// and finally cancels the returned context.
func (c *Coordinator) Wait(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer cancel()
		select {
		case <-sigs:
		case <-parent.Done():
			return
		}
		c.mu.Lock()
		idle := c.active == 0
		c.mu.Unlock()
		if idle {
			return
		}
		select {
		case <-c.drain:
		case <-time.After(c.timeout):
			_, _ = io.WriteString(c.w, "portwatch: drain timeout, forcing shutdown\n")
		}
	}()
	return ctx
}
