// Package watchlist maintains an explicit set of ports the operator wants to
// ensure remain open. It is the inverse of the suppress package: while
// suppress silences alerts for known-closed ports, watchlist raises an alert
// when a port the operator expects to be open is found absent from the scan
// results.
//
// # Usage
//
//	wl := watchlist.New()
//	wl.Add(22, "tcp")   // SSH must always be reachable
//	wl.Add(443, "tcp")  // HTTPS must always be reachable
//
//	violations := watchlist.Check(wl, currentPorts)
//	for _, v := range violations {
//		log.Printf("watchlist violation: %s/%d — %s",
//			v.Entry.Protocol, v.Entry.Port, v.Reason)
//	}
package watchlist
