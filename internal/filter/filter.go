package filter

import "github.com/user/portwatch/internal/scanner"

// Rule describes a single port filter rule.
type Rule struct {
	Port     uint16
	Protocol string // "tcp" or "udp"
}

// Filter decides which ports are relevant based on a set of rules.
// When the allow-list is empty every port is considered relevant.
type Filter struct {
	allowList map[Rule]struct{}
}

// New creates a Filter from the given rules.
// Passing no rules produces a pass-through filter that accepts everything.
func New(rules []Rule) *Filter {
	f := &Filter{
		allowList: make(map[Rule]struct{}, len(rules)),
	}
	for _, r := range rules {
		f.allowList[r] = struct{}{}
	}
	return f
}

// Apply returns only the ports that match the filter rules.
// If the filter has no rules, all ports are returned unchanged.
func (f *Filter) Apply(ports []scanner.PortState) []scanner.PortState {
	if len(f.allowList) == 0 {
		return ports
	}

	out := make([]scanner.PortState, 0, len(ports))
	for _, p := range ports {
		key := Rule{Port: p.Port, Protocol: p.Protocol}
		if _, ok := f.allowList[key]; ok {
			out = append(out, p)
		}
	}
	return out
}

// Contains reports whether the given port/protocol pair is in the allow-list.
// Always returns true when the allow-list is empty.
func (f *Filter) Contains(port uint16, protocol string) bool {
	if len(f.allowList) == 0 {
		return true
	}
	_, ok := f.allowList[Rule{Port: port, Protocol: protocol}]
	return ok
}
