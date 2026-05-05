package portage

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// ExpiryPolicy defines how long a closed port is retained before eviction.
type ExpiryPolicy struct {
	// MaxAge is the duration after which a closed port entry is considered expired.
	MaxAge time.Duration
}

// DefaultExpiryPolicy returns a policy that expires closed ports after 72 hours.
func DefaultExpiryPolicy() ExpiryPolicy {
	return ExpiryPolicy{MaxAge: 72 * time.Hour}
}

// Expired reports whether the given port entry's LastSeen time is older than
// the policy's MaxAge relative to now.
func (p ExpiryPolicy) Expired(entry Entry, now time.Time) bool {
	if p.MaxAge <= 0 {
		return false
	}
	return now.Sub(entry.LastSeen) > p.MaxAge
}

// PruneExpired removes entries from entries whose LastSeen is older than the
// policy's MaxAge. It returns the pruned slice and the list of evicted ports.
func PruneExpired(entries []Entry, policy ExpiryPolicy, now time.Time) (kept []Entry, evicted []scanner.Port) {
	for _, e := range entries {
		if policy.Expired(e, now) {
			evicted = append(evicted, e.Port)
		} else {
			kept = append(kept, e)
		}
	}
	return kept, evicted
}
