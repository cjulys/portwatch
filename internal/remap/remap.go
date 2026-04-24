// Package remap translates raw scanner ports through a user-defined
// port-alias map so that display names and rule matching can use
// friendly labels (e.g. "http" instead of "80/tcp").
package remap

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Rule maps a protocol+port pair to a human-readable alias.
type Rule struct {
	Proto string // "tcp" or "udp"
	Port  uint16
	Alias string
}

// Mapper holds the compiled alias lookup table.
type Mapper struct {
	table map[string]string // key: "proto:port"
}

// New builds a Mapper from the provided rules.
// Duplicate keys are silently overwritten by the last rule supplied.
func New(rules []Rule) *Mapper {
	t := make(map[string]string, len(rules))
	for _, r := range rules {
		t[key(r.Proto, r.Port)] = strings.TrimSpace(r.Alias)
	}
	return &Mapper{table: t}
}

// Apply returns a copy of ports with the Alias field populated where a
// matching rule exists. Ports with no matching rule are returned unchanged.
func (m *Mapper) Apply(ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, len(ports))
	for i, p := range ports {
		if alias, ok := m.table[key(p.Proto, p.Port)]; ok {
			p.Alias = alias
		}
		out[i] = p
	}
	return out
}

// Alias returns the alias for a given proto/port pair, or an empty string
// when no rule matches.
func (m *Mapper) Alias(proto string, port uint16) string {
	return m.table[key(proto, port)]
}

// Len returns the number of rules loaded into the mapper.
func (m *Mapper) Len() int { return len(m.table) }

func key(proto string, port uint16) string {
	return fmt.Sprintf("%s:%d", strings.ToLower(proto), port)
}
