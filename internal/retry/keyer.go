package retry

import "fmt"

// OpKey returns a string key identifying a retryable operation by name and target.
// Useful for logging or per-operation throttle integration.
func OpKey(operation, target string) string {
	return fmt.Sprintf("%s:%s", operation, target)
}

// ScanKey returns a key for a scan retry operation on a given host.
func ScanKey(host string) string {
	return OpKey("scan", host)
}

// WebhookKey returns a key for a webhook delivery retry.
func WebhookKey(url string) string {
	return OpKey("webhook", url)
}
