// Package digest computes a fingerprint of a port snapshot so callers can
// detect whether anything changed without doing a full diff.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"portwatch/internal/scanner"
)

// Digest is a hex-encoded SHA-256 fingerprint of an ordered port list.
type Digest string

// Empty is the fingerprint of a nil / zero-length port list.
const Empty Digest = ""

// Compute returns a stable Digest for the supplied ports.
// The slice is sorted internally so order does not affect the result.
func Compute(ports []scanner.Port) Digest {
	if len(ports) == 0 {
		return Empty
	}

	sorted := make([]scanner.Port, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Proto != sorted[j].Proto {
			return sorted[i].Proto < sorted[j].Proto
		}
		return sorted[i].Number < sorted[j].Number
	})

	h := sha256.New()
	for _, p := range sorted {
		fmt.Fprintf(h, "%s:%d\n", p.Proto, p.Number)
	}
	return Digest(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two digests are identical.
func Equal(a, b Digest) bool { return a == b }
