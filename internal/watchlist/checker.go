package watchlist

import (
	"github.com/user/portwatch/internal/scanner"
)

// Violation describes a watchlisted port that is not currently open.
type Violation struct {
	Entry Entry
	// Reason is a human-readable explanation.
	Reason string
}

// Check compares the current set of open ports against the watchlist and
// returns a Violation for every watched port that is absent from current.
func Check(wl *Watchlist, current []scanner.Port) []Violation {
	open := make(map[string]struct{}, len(current))
	for _, p := range current {
		open[key(p.Port, p.Protocol)] = struct{}{}
	}

	var violations []Violation
	for _, e := range wl.All() {
		if _, ok := open[key(e.Port, e.Protocol)]; !ok {
			violations = append(violations, Violation{
				Entry:  e,
				Reason: "watched port is not open",
			})
		}
	}
	return violations
}

// ViolationCount is a convenience wrapper that returns the number of
// violations without allocating a full slice in the caller.
func ViolationCount(wl *Watchlist, current []scanner.Port) int {
	return len(Check(wl, current))
}
