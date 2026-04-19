// Package circuitbreaker provides a thread-safe circuit breaker used to
// protect portwatch from cascading failures when external notification
// targets (webhooks, SMTP relays) become unavailable.
//
// Usage:
//
//	cb := circuitbreaker.New(3, 30*time.Second)
//	if cb.Allow() {
//		err := sendWebhook(payload)
//		if err != nil {
//			cb.RecordFailure()
//		} else {
//			cb.RecordSuccess()
//		}
//	}
package circuitbreaker
