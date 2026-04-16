// Package heartbeat provides a lightweight liveness signal for the portwatch
// daemon. It periodically writes a timestamped line to a configurable writer
// (default: stderr) so that process supervisors or log monitors can detect
// when the daemon has stalled or crashed silently.
package heartbeat
