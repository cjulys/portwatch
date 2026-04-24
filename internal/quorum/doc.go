// Package quorum implements a confirmation-count gate for port-change events.
//
// A single anomalous scan result can be caused by transient OS scheduling,
// a brief firewall rule flush, or measurement jitter. Quorum addresses this
// by requiring that the same change (port opened or closed) be observed on
// N consecutive scans before the event is propagated to alerting subsystems.
//
// Typical usage:
//
//	q := quorum.New(3) // require 3 confirmations
//	if q.Observe(port, "opened") {
//		// confirmed — forward to notifier
//	}
package quorum
