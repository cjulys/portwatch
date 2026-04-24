// Package evict implements a TTL-based eviction tracker for observed ports.
//
// Callers should call Touch after each successful scan for every visible port,
// then call Evict periodically to obtain ports that have not been seen within
// the configured time-to-live. Evicted ports can be treated as implicitly
// closed and forwarded to the alert or notifier pipeline.
package evict
