// Package backoff implements an exponential back-off strategy used when
// retrying scans, webhook deliveries, or any other operation that may
// encounter transient failures.
//
// Usage:
//
//	b := backoff.New(backoff.DefaultPolicy())
//	for b.Next(ctx) {
//		if err := doWork(); err == nil {
//			break
//		}
//	}
package backoff
