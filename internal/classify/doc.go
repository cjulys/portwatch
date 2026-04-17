// Package classify assigns severity levels and categories to port change
// events detected during a scan cycle.
//
// Usage:
//
//	c := classify.New([]uint16{22, 3306}) // ports treated as critical
//	result := c.ClassifyOpened(port)
//	batch := c.Batch(classify.DiffInput{Opened: opened, Closed: closed})
//
// Severity levels:
//   - critical: port is in the configured critical list and has opened
//   - warning:  any other port has opened
//   - info:     a port has closed
package classify
