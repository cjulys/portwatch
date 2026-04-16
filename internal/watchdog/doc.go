// Package watchdog provides a heartbeat-based liveness monitor for the
// portwatch scan loop. If Beat is not called within the configured timeout
// the watchdog writes a stall alert to the configured writer (default stderr).
//
// Typical usage:
//
//	wd := watchdog.New(30*time.Second, os.Stderr)
//	defer wd.Stop()
//	for each scan cycle {
//	    scan()
//	    wd.Beat()
//	}
package watchdog
