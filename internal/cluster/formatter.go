package cluster

import (
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Summary returns a compact, comma-separated string listing all port ranges
// derived from ports. Useful for single-line log messages.
//
// Example output: "tcp/22, tcp/80-82, udp/53"
func Summary(ports []scanner.Port) string {
	ranges := Group(ports)
	if len(ranges) == 0 {
		return ""
	}
	parts := make([]string, len(ranges))
	for i, r := range ranges {
		parts[i] = r.String()
	}
	return strings.Join(parts, ", ")
}

// RangeCount returns the number of distinct ranges produced by grouping ports.
// A return value much smaller than len(ports) indicates heavy clustering.
func RangeCount(ports []scanner.Port) int {
	return len(Group(ports))
}
