package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"portwatch/internal/scanner"
)

// Compute returns a hex-encoded SHA-256 digest of the given port slice.
// The result is order-independent: ports are sorted before hashing.
func Compute(ports []scanner.Port) string {
	if len(ports) == 0 {
		return ""
	}

	sorted := make([]scanner.Port, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		ki := fmt.Sprintf("%s:%d", sorted[i].Protocol, sorted[i].Number)
		kj := fmt.Sprintf("%s:%d", sorted[j].Protocol, sorted[j].Number)
		return ki < kj
	})

	h := sha256.New()
	for _, p := range sorted {
		b, err := json.Marshal(p)
		if err != nil {
			continue
		}
		h.Write(b)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Equal returns true when two port slices produce the same digest.
func Equal(a, b []scanner.Port) bool {
	return Compute(a) == Compute(b)
}
