package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time summary of watcher activity.
type Snapshot struct {
	ScansTotal    int64
	AlertsTotal   int64
	LastScanAt    time.Time
	LastAlertAt   time.Time
	OpenPortCount int
}

// Collector accumulates runtime metrics for portwatch.
type Collector struct {
	mu            sync.RWMutex
	scansTotal    int64
	alertsTotal   int64
	lastScanAt    time.Time
	lastAlertAt   time.Time
	openPortCount int
}

// New returns an initialised Collector.
func New() *Collector {
	return &Collector{}
}

// RecordScan records a completed scan with the current open-port count.
func (c *Collector) RecordScan(openPorts int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scansTotal++
	c.lastScanAt = time.Now()
	c.openPortCount = openPorts
}

// RecordAlert records that an alert was emitted.
func (c *Collector) RecordAlert() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.alertsTotal++
	c.lastAlertAt = time.Now()
}

// Snapshot returns a consistent copy of current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Snapshot{
		ScansTotal:    c.scansTotal,
		AlertsTotal:   c.alertsTotal,
		LastScanAt:    c.lastScanAt,
		LastAlertAt:   c.lastAlertAt,
		OpenPortCount: c.openPortCount,
	}
}

// Reset zeroes all counters (useful in tests).
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	*c = Collector{}
}
