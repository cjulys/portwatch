// Package debounce provides a Debouncer that coalesces rapid scanner.Diff
// events into a single handler invocation after a configurable quiet window.
//
// This is useful when a service restart causes a port to briefly disappear and
// reappear; without debouncing, two spurious alerts (closed + opened) would be
// emitted. With debouncing, only the net change is forwarded once the port
// state stabilises.
package debounce
