// Package portname resolves well-known port numbers to human-readable service names.
package portname

import "fmt"

// Resolver maps port/protocol pairs to service names.
type Resolver struct {
	overrides map[string]string
}

// builtinTable contains a curated set of well-known service names.
var builtinTable = map[string]string{
	"tcp/21":   "ftp",
	"tcp/22":   "ssh",
	"tcp/23":   "telnet",
	"tcp/25":   "smtp",
	"tcp/53":   "dns",
	"udp/53":   "dns",
	"tcp/80":   "http",
	"tcp/110":  "pop3",
	"tcp/143":  "imap",
	"tcp/443":  "https",
	"tcp/465":  "smtps",
	"tcp/587":  "submission",
	"tcp/993":  "imaps",
	"tcp/995":  "pop3s",
	"tcp/3306": "mysql",
	"tcp/5432": "postgres",
	"tcp/6379": "redis",
	"tcp/8080": "http-alt",
	"tcp/8443": "https-alt",
	"tcp/27017": "mongodb",
}

// New returns a Resolver with optional override rules.
// Overrides are expressed as "proto/port" -> name, e.g. "tcp/8080" -> "myapp".
func New(overrides map[string]string) *Resolver {
	if overrides == nil {
		overrides = make(map[string]string)
	}
	return &Resolver{overrides: overrides}
}

// Resolve returns the service name for the given protocol and port number.
// It checks caller-supplied overrides first, then the builtin table.
// If no match is found it returns an empty string.
func (r *Resolver) Resolve(proto string, port uint16) string {
	k := key(proto, port)
	if name, ok := r.overrides[k]; ok {
		return name
	}
	return builtinTable[k]
}

// ResolveOrDefault is like Resolve but returns fallback when no name is found.
func (r *Resolver) ResolveOrDefault(proto string, port uint16, fallback string) string {
	if name := r.Resolve(proto, port); name != "" {
		return name
	}
	return fallback
}

func key(proto string, port uint16) string {
	return fmt.Sprintf("%s/%d", proto, port)
}
