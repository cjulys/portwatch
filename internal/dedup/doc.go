// Package dedup provides deduplication of port-change events.
//
// A Deduplicator tracks the last time each (port, direction) pair was
// reported and suppresses repeated events that fall within a configurable
// time window.  This prevents alert storms caused by ports that flap
// rapidly open and closed.
//
// Usage:
//
//	d := dedup.New(30 * time.Second)
//	if !d.IsDuplicate(p, "opened") {
//		// forward the alert
//	}
package dedup
