package portage

import "fmt"

// AgeKey returns a canonical string key for a port's age tracking entry.
// It combines protocol and port number to ensure distinct tracking per service.
func AgeKey(proto string, port int) string {
	return fmt.Sprintf("%s:%d", proto, port)
}

// BucketKey returns a key that maps a port into an age bucket label.
// Bucket labels are: "new" (<1h), "recent" (<24h), "established" (<7d), "old" (>=7d).
func BucketKey(proto string, port int, ageSeconds int64) string {
	return fmt.Sprintf("%s:%d:%s", proto, port, bucketLabel(ageSeconds))
}

// bucketLabel classifies an age in seconds into a human-readable label.
func bucketLabel(ageSeconds int64) string {
	const (
		hour  = 3600
		day   = 86400
		week  = 604800
	)
	switch {
	case ageSeconds < hour:
		return "new"
	case ageSeconds < day:
		return "recent"
	case ageSeconds < week:
		return "established"
	default:
		return "old"
	}
}
