package classify

import "github.com/user/portwatch/internal/scanner"

// DiffInput holds the opened and closed port slices from a scan diff.
type DiffInput struct {
	Opened []scanner.Port
	Closed []scanner.Port
}

// BatchResult is the full classification of a diff cycle.
type BatchResult struct {
	Results  []Result
	Critical int
	Warnings int
	Info     int
}

// Batch classifies all opened and closed ports in one call.
func (c *Classifier) Batch(d DiffInput) BatchResult {
	out := BatchResult{}
	for _, p := range d.Opened {
		r := c.ClassifyOpened(p)
		out.Results = append(out.Results, r)
		switch r.Severity {
		case SeverityCritical:
			out.Critical++
		case SeverityWarning:
			out.Warnings++
		}
	}
	for _, p := range d.Closed {
		r := c.ClassifyClosed(p)
		out.Results = append(out.Results, r)
		out.Info++
	}
	return out
}
