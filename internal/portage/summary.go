package portage

import "fmt"

// AgeSummary holds aggregate counts of ports grouped by age bucket.
type AgeSummary struct {
	New         int
	Recent      int
	Established int
	LongTerm    int
	Total       int
}

// String returns a human-readable one-line summary.
func (s AgeSummary) String() string {
	return fmt.Sprintf(
		"total=%d new=%d recent=%d established=%d long-term=%d",
		s.Total, s.New, s.Recent, s.Established, s.LongTerm,
	)
}

// Summarize builds an AgeSummary from a slice of Entry values.
// Entries with an unrecognised bucket label are counted only in Total.
func Summarize(entries []Entry) AgeSummary {
	var s AgeSummary
	for _, e := range entries {
		s.Total++
		switch e.Bucket {
		case BucketNew:
			s.New++
		case BucketRecent:
			s.Recent++
		case BucketEstablished:
			s.Established++
		case BucketLongTerm:
			s.LongTerm++
		}
	}
	return s
}

// SummarizeAll is a convenience wrapper that calls Summarize on all entries
// returned by store.All.
func SummarizeAll(store *Store) AgeSummary {
	return Summarize(store.All())
}
