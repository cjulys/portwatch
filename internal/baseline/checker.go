package baseline

import (
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/scanner/diff"
)

// Violation describes a port that deviates from the baseline.
type Violation struct {
	Port   scanner.Port
	Reason string // "unexpected_open" | "expected_closed"
}

// Check compares current ports against the baseline and returns any
// violations. If the baseline is empty every port is considered a violation.
func Check(base []scanner.Port, current []scanner.Port) []Violation {
	diffs := diff.Compare(base, current)
	var violations []Violation
	for _, d := range diffs {
		v := Violation{Port: d.Port}
		switch d.Type {
		case diff.Added:
			v.Reason = "unexpected_open"
		case diff.Removed:
			v.Reason = "expected_closed"
		default:
			continue
		}
		violations = append(violations, v)
	}
	return violations
}

// ViolationCount returns the number of violations of each reason type.
func ViolationCount(violations []Violation) (unexpected, closed int) {
	for _, v := range violations {
		switch v.Reason {
		case "unexpected_open":
			unexpected++
		case "expected_closed":
			closed++
		}
	}
	return
}
