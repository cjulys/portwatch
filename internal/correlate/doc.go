// Package correlate groups bursts of port-change diffs into correlated
// events identified by a shared ID.
//
// A Correlator accumulates scanner.Diff values that arrive within a
// configurable quiet window. Once no new diffs arrive for the duration of
// that window the accumulated batch is emitted as a single Event carrying a
// random correlation ID and a UTC timestamp.
//
// Typical use:
//
//	c := correlate.New(500*time.Millisecond, func(ev correlate.Event) {
//		fmt.Println(ev.ID, len(ev.Diffs))
//	})
//	c.Add(diffs)
package correlate
