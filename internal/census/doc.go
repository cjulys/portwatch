// Package census tracks port population statistics across successive scans.
//
// Usage:
//
//	c := census.New()
//	snap, delta := c.Record(ports)
//	fmt.Printf("total=%d delta=%+d\n", snap.Total, delta.Total)
//
// Census is safe for concurrent use.
package census
