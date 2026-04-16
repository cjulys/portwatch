// Package audit provides an append-only audit log for portwatch.
//
// Each scan cycle and notable event (port opened, port closed, baseline
// violation) is recorded as a JSON line in the configured audit file.
// Entries can be read back with ReadAll for offline inspection or
// integration with external tooling.
//
// Usage:
//
//	l, err := audit.New("/var/lib/portwatch/audit.log")
//	if err != nil { ... }
//	l.Record("port_opened", "8080/tcp")
package audit
