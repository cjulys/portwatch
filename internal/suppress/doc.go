// Package suppress provides time-bounded suppression rules for port
// change alerts. A Suppressor holds a set of rules keyed by protocol,
// address, and port. Rules expire after a configurable duration, after
// which the corresponding port will once again generate alerts.
//
// Rules may be exact (matching a specific address) or wildcard
// (matching any address on a given protocol/port pair). Wildcard rules
// are looked up via WildcardKey when an exact match is absent.
//
// Rules are persisted to a JSON file so that suppressions survive
// process restarts.
package suppress
