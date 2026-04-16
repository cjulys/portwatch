// Package snapshot provides utilities for capturing and persisting
// point-in-time records of open ports discovered by the scanner.
//
// A Snapshot pairs a UTC timestamp with the list of ports seen during
// a single scan cycle. Snapshots can be saved to disk as JSON and
// reloaded across daemon restarts, enabling before/after comparisons
// and offline inspection of historical scan results.
package snapshot
