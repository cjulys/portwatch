// Package escalate promotes alert severity for port changes that recur
// within a sliding time window. A port that opens and closes repeatedly
// starts at Info, escalates to Warning after warnAt occurrences, and
// reaches Critical after critAt occurrences — all within the configured
// window. Counts reset automatically once the window expires.
package escalate
