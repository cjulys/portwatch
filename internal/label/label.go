// Package label attaches human-readable labels to ports based on
// configurable rules, falling back to a built-in well-known-port table.
package label

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Rule maps a port/protocol pair to a label string.
type Rule struct {
	Port     uint16
	Protocol string // "tcp" | "udp" | "" means any
	Label    string
}

// Labeler assigns labels to scanner.Port values.
type Labeler struct {
	rules    []Rule
	builtins map[string]string // key: "proto:port"
}

// New returns a Labeler seeded with the provided override rules.
// Built-in well-known ports are always registered as a fallback.
func New(overrides []Rule) *Labeler {
	l := &Labeler{
		rules:    overrides,
		builtins: builtinTable(),
	}
	return l
}

// Apply returns the label for the given port, or an empty string when no
// rule matches and the port is not in the built-in table.
func (l *Labeler) Apply(p scanner.Port) string {
	for _, r := range l.rules {
		if r.Port != p.Port {
			continue
		}
		if r.Protocol != "" && !strings.EqualFold(r.Protocol, p.Protocol) {
			continue
		}
		return r.Label
	}
	key := fmt.Sprintf("%s:%d", strings.ToLower(p.Protocol), p.Port)
	if v, ok := l.builtins[key]; ok {
		return v
	}
	return ""
}

// ApplyAll labels every port in the slice, returning a map keyed by
// "proto:port" for convenient lookup.
func (l *Labeler) ApplyAll(ports []scanner.Port) map[string]string {
	out := make(map[string]string, len(ports))
	for _, p := range ports {
		if lbl := l.Apply(p); lbl != "" {
			out[fmt.Sprintf("%s:%d", strings.ToLower(p.Protocol), p.Port)] = lbl
		}
	}
	return out
}

func builtinTable() map[string]string {
	return map[string]string{
		"tcp:21":   "FTP",
		"tcp:22":   "SSH",
		"tcp:25":   "SMTP",
		"tcp:53":   "DNS",
		"udp:53":   "DNS",
		"tcp:80":   "HTTP",
		"tcp:443":  "HTTPS",
		"tcp:3306": "MySQL",
		"tcp:5432": "PostgreSQL",
		"tcp:6379": "Redis",
		"tcp:8080": "HTTP-alt",
		"tcp:8443": "HTTPS-alt",
		"tcp:27017": "MongoDB",
	}
}
