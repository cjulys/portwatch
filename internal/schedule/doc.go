// Package schedule provides a lightweight named-interval scheduler used by
// portwatch to coordinate recurring tasks such as port scans and state flushes.
//
// Usage:
//
//	s := schedule.New()
//	s.Register("scan", 30*time.Second)
//	s.Register("flush", 5*time.Minute)
//
//	for {
//		time.Sleep(time.Until(s.Next()))
//		for _, name := range s.Due() {
//			// dispatch work by name
//			_ = name
//		}
//	}
package schedule
