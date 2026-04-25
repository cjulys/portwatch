// Package normalize canonicalizes raw scanner.Port slices before they
// are handed off to the comparison, storage, or alerting layers.
//
// Typical usage:
//
//	n := normalize.New(
//		normalize.WithLowerProtocol(),
//		normalize.WithTrimAddress(),
//		normalize.WithDeduplication(),
//	)
//	clean := n.Apply(rawPorts)
package normalize
