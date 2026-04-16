// Package rollup provides a time-window batching layer for port-change events.
//
// Instead of emitting an alert for every individual diff, Rollup accumulates
// opened and closed ports over a configurable window duration and calls the
// registered flush function once with the combined Event.
//
// Typical usage:
//
//	r := rollup.New(5*time.Second, func(e rollup.Event) {
//		// handle batched event
//	})
//	r.Add(opened, closed)
package rollup
