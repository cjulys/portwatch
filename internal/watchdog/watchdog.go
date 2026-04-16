// Package watchdog detects when the scan loop stalls and emits an alert.
package watchdog

import (
	"io"
	"os"
	"sync"
	"time"
)

// Watchdog monitors heartbeat signals and writes an alert when they stop.
type Watchdog struct {
	timeout  time.Duration
	writer   io.Writer
	timer    *time.Timer
	mu       sync.Mutex
	stopped  bool
}

// New creates a Watchdog that fires after timeout with no heartbeat.
func New(timeout time.Duration, w io.Writer) *Watchdog {
	if w == nil {
		w = os.Stderr
	}
	wd := &Watchdog{timeout: timeout, writer: w}
	wd.timer = time.AfterFunc(timeout, wd.fire)
	return wd
}

// Beat resets the watchdog timer; call this after each successful scan.
func (wd *Watchdog) Beat() {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	if !wd.stopped {
		wd.timer.Reset(wd.timeout)
	}
}

// Stop disables the watchdog.
func (wd *Watchdog) Stop() {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	wd.stopped = true
	wd.timer.Stop()
}

func (wd *Watchdog) fire() {
	wd.mu.Lock()
	defer wd.mu.Unlock()
	if !wd.stopped {
		_, _ = io.WriteString(wd.writer, "[watchdog] scan loop stalled — no heartbeat received\n")
	}
}
