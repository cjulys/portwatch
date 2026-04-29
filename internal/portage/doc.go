// Package portage tracks the age of open ports observed by portwatch.
//
// A Tracker maintains a persistent record of when each port was first seen and
// when it was last seen. The aging sub-module builds on this to assign a
// human-readable AgeCategory (new / recent / established / long-term) to every
// port in a scan snapshot.
//
// Typical usage:
//
//	tracker := portage.New(stateDir)
//	tracker.Update(currentPorts, closedPorts)
//	agedPorts := portage.EnrichAll(tracker, currentPorts, time.Now())
//	for _, ap := range agedPorts {
//		fmt.Printf("%s:%d  age=%s\n", ap.Port.Protocol, ap.Port.Port, ap.Category)
//	}
package portage
