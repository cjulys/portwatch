package debounce

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// DiffKey returns a stable string key for a scanner.Diff, suitable for
// deduplication or logging.
func DiffKey(d scanner.Diff) string {
	return fmt.Sprintf("%s:%d:%s", d.Port.Proto, d.Port.Port, d.Kind)
}

// Deduplicate removes duplicate diffs from a slice, preserving first-seen
// order. Two diffs are considered equal when their DiffKey matches.
func Deduplicate(diffs []scanner.Diff) []scanner.Diff {
	seen := make(map[string]struct{}, len(diffs))
	out := diffs[:0:0]
	for _, d := range diffs {
		k := DiffKey(d)
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, d)
	}
	return out
}
