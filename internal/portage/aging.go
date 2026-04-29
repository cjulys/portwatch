package portage

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// AgeCategory classifies how long a port has been observed open.
type AgeCategory int

const (
	AgeNew      AgeCategory = iota // seen for less than 1 hour
	AgeRecent                      // seen for 1–24 hours
	AgeEstablished                 // seen for 1–7 days
	AgeLongTerm                    // seen for more than 7 days
)

func (a AgeCategory) String() string {
	switch a {
	case AgeNew:
		return "new"
	case AgeRecent:
		return "recent"
	case AgeEstablished:
		return "established"
	case AgeLongTerm:
		return "long-term"
	default:
		return "unknown"
	}
}

// Categorize returns the AgeCategory for a port given when it was first seen
// and the current time.
func Categorize(firstSeen, now time.Time) AgeCategory {
	age := now.Sub(firstSeen)
	switch {
	case age < time.Hour:
		return AgeNew
	case age < 24*time.Hour:
		return AgeRecent
	case age < 7*24*time.Hour:
		return AgeEstablished
	default:
		return AgeLongTerm
	}
}

// AgedPort pairs a scanned port with its computed age metadata.
type AgedPort struct {
	Port      scanner.Port
	FirstSeen time.Time
	Category  AgeCategory
}

// EnrichAll annotates each port in the slice with its age category, using the
// supplied Tracker to retrieve first-seen timestamps.
func EnrichAll(tracker *Tracker, ports []scanner.Port, now time.Time) []AgedPort {
	out := make([]AgedPort, 0, len(ports))
	for _, p := range ports {
		fs, ok := tracker.FirstSeen(p)
		if !ok {
			fs = now
		}
		out = append(out, AgedPort{
			Port:      p,
			FirstSeen: fs,
			Category:  Categorize(fs, now),
		})
	}
	return out
}
