// Package tagger assigns human-readable labels to ports based on well-known
// service mappings and user-defined overrides.
package tagger

import "fmt"

// well-known port → service name
var builtins = map[string]string{
	"tcp:22":   "ssh",
	"tcp:80":   "http",
	"tcp:443":  "https",
	"tcp:3306": "mysql",
	"tcp:5432": "postgres",
	"tcp:6379": "redis",
	"tcp:8080": "http-alt",
	"udp:53":   "dns",
	"udp:123":  "ntp",
}

// Tagger maps ports to labels.
type Tagger struct {
	overrides map[string]string
}

// New returns a Tagger. overrides may be nil.
func New(overrides map[string]string) *Tagger {
	o := make(map[string]string, len(overrides))
	for k, v := range overrides {
		o[k] = v
	}
	return &Tagger{overrides: o}
}

// Tag returns the label for the given protocol and port number.
// It checks user overrides first, then builtins, then returns a
// generic "port/<n>" label.
func (t *Tagger) Tag(proto string, port uint16) string {
	key := fmt.Sprintf("%s:%d", proto, port)
	if label, ok := t.overrides[key]; ok {
		return label
	}
	if label, ok := builtins[key]; ok {
		return label
	}
	return fmt.Sprintf("port/%d", port)
}

// Known reports whether the port has a builtin or override label.
func (t *Tagger) Known(proto string, port uint16) bool {
	key := fmt.Sprintf("%s:%d", proto, port)
	_, inOverride := t.overrides[key]
	_, inBuiltin := builtins[key]
	return inOverride || inBuiltin
}
