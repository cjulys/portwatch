// Package cluster groups consecutive port ranges into compact representations
// for cleaner alerting and reporting output.
package cluster

import (
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Range represents a contiguous run of ports sharing the same protocol.
type Range struct {
	Protocol string
	Start    uint16
	End      uint16
}

// String returns a human-readable representation such as "tcp/80" or "tcp/8080-8090".
func (r Range) String() string {
	if r.Start == r.End {
		return fmt.Sprintf("%s/%d", r.Protocol, r.Start)
	}
	return fmt.Sprintf("%s/%d-%d", r.Protocol, r.Start, r.End)
}

// Size returns the number of ports in the range.
func (r Range) Size() int {
	return int(r.End-r.Start) + 1
}

// Group collapses a slice of ports into contiguous ranges per protocol.
// Ports are sorted before grouping so input order does not matter.
func Group(ports []scanner.Port) []Range {
	if len(ports) == 0 {
		return nil
	}

	sorted := make([]scanner.Port, len(ports))
	copy(sorted, ports)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Protocol != sorted[j].Protocol {
			return sorted[i].Protocol < sorted[j].Protocol
		}
		return sorted[i].Port < sorted[j].Port
	})

	var ranges []Range
	cur := Range{
		Protocol: sorted[0].Protocol,
		Start:    sorted[0].Port,
		End:      sorted[0].Port,
	}

	for _, p := range sorted[1:] {
		if p.Protocol == cur.Protocol && p.Port == cur.End+1 {
			cur.End = p.Port
			continue
		}
		ranges = append(ranges, cur)
		cur = Range{Protocol: p.Protocol, Start: p.Port, End: p.Port}
	}
	ranges = append(ranges, cur)
	return ranges
}
