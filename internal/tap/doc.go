// Package tap implements a passive diff tap for portwatch.
//
// A Tap sits alongside the main event pipeline and forwards copies of
// scanner.Diff slices to any number of registered Sink functions without
// blocking or altering the primary flow.  It is safe for concurrent use.
//
// Typical usage:
//
//	t := tap.New(os.Stderr)
//	t.Register(func(diffs []scanner.Diff) {
//		for _, d := range diffs {
//			log.Printf("tap: %v", d)
//		}
//	})
//	// later, inside the scan loop:
//	t.Send(diffs)
package tap
