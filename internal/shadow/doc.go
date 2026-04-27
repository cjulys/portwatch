// Package shadow provides a two-generation port-scan buffer used to detect
// flapping ports — ports whose open/closed state oscillates between
// consecutive scan cycles.
//
// Usage:
//
//	s := shadow.New()
//	s.Commit(firstScanPorts)
//	s.Commit(secondScanPorts)
//	flapping := s.Flapping() // ports that changed between the two scans
package shadow
