// Package budget provides a rolling-window scan-rate limiter for portwatch.
//
// A Budget caps the number of port scans that may be initiated within a
// configurable time window, preventing runaway scanning under high load or
// misconfiguration. When the budget is exhausted a warning is written to the
// fallback writer (default: stderr) and the caller is expected to skip the
// scan cycle.
//
// Usage:
//
//	b := budget.New(60, time.Minute) // at most 60 scans per minute
//	if b.Allow(ctx) {
//		// perform scan
//	}
package budget
