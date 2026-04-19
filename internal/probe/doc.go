// Package probe provides active connectivity verification for open ports
// discovered by the scanner. It dials each port and reports reachability
// and round-trip latency, allowing portwatch to distinguish ports that are
// merely listed by the OS from those that are genuinely accepting connections.
package probe
