// Package fingerprint produces a stable string identity for a port scan result
// that can be used as a cache key or change-detection token.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"portwatch/internal/scanner"
)

// Port returns a short hex fingerprint for a single port entry.
func Port(p scanner.Port) string {
	h := sha256.New()
	fmt.Fprintf(h, "%s:%d:%s", p.Protocol, p.Number, p.State)
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// Ports returns a fingerprint that covers an entire slice of ports.
// The result is order-independent: the same set of ports always produces
// the same fingerprint regardless of slice ordering.
func Ports(ports []scanner.Port) string {
	if len(ports) == 0 {
		return ""
	}

	tokens := make([]string, len(ports))
	for i, p := range ports {
		tokens[i] = Port(p)
	}
	sort.Strings(tokens)

	h := sha256.New()
	for _, t := range tokens {
		h.Write([]byte(t))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// Changed returns true when the fingerprint of current differs from prev.
func Changed(prev []scanner.Port, current []scanner.Port) bool {
	return Ports(prev) != Ports(current)
}

// Diff returns ports that are present in current but absent in prev (added)
// and ports present in prev but absent in current (removed).
func Diff(prev, current []scanner.Port) (added, removed []scanner.Port) {
	prevIndex := make(map[string]struct{}, len(prev))
	for _, p := range prev {
		prevIndex[Port(p)] = struct{}{}
	}

	currIndex := make(map[string]struct{}, len(current))
	for _, p := range current {
		currIndex[Port(p)] = struct{}{}
		if _, ok := prevIndex[Port(p)]; !ok {
			added = append(added, p)
		}
	}

	for _, p := range prev {
		if _, ok := currIndex[Port(p)]; !ok {
			removed = append(removed, p)
		}
	}
	return added, removed
}
