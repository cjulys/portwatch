// Package summary aggregates port-change history into periodic digest
// reports. A Builder queries the history store for a given time window
// and produces a Report that can be printed or forwarded to a notifier.
package summary
