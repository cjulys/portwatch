// Package reachable tracks how consistently each port has been observed
// open across recent scans and exposes a Ratio in [0,1] that other
// components (e.g. classify, escalate) can use to filter noisy or
// transient port events.
package reachable
