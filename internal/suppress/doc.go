// Package suppress manages per-port alert suppression rules for portwatch.
//
// A suppression rule silences notifications for a specific port/protocol pair
// either indefinitely (zero Until) or until a given time. Rules are persisted
// to a JSON file so they survive daemon restarts.
//
// Typical usage:
//
//	store := suppress.New("/var/lib/portwatch/suppress.json")
//
//	// suppress port 8080/tcp for one hour
//	store.Add(suppress.Rule{
//		Port:     8080,
//		Protocol: "tcp",
//		Until:    time.Now().Add(time.Hour),
//	})
//
//	// check before dispatching an alert
//	if !store.IsSuppressed(port, proto) {
//		notifier.Dispatch(event)
//	}
package suppress
