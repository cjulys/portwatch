// Package metrics provides a lightweight, thread-safe collector for
// portwatch runtime statistics.
//
// A Collector accumulates scan and alert counters and exposes them via
// Snapshot for use by reporters, dashboards, or health endpoints.
//
// Usage:
//
//	c := metrics.New()
//	c.RecordScan(len(ports))
//	c.RecordAlert()
//	s := c.Snapshot()
//	fmt.Println(s.ScansTotal, s.OpenPortCount)
package metrics
