// Package decay provides time-based severity decay for repeated port-change
// alerts. When a port remains in an unexpected state for an extended period,
// repeated alerts are downgraded from critical → warning → info → silenced so
// that operators are not overwhelmed by persistent noise.
package decay
