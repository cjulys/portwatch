// Package window implements a sliding time-window counter used to measure
// event frequency over a configurable rolling duration. It is safe for
// concurrent use and is used by rate-limiting and escalation components
// throughout portwatch.
package window
