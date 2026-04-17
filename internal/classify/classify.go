// Package classify labels a port diff event with a severity and category.
package classify

import "github.com/user/portwatch/internal/scanner"

// Severity represents how serious a port change is.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// Category describes the kind of change.
type Category string

const (
	CategoryNewPort    Category = "new_port"
	CategoryClosedPort Category = "closed_port"
)

// Result holds the classification of a single port event.
type Result struct {
	Port     scanner.Port
	Category Category
	Severity Severity
	Label    string
}

// Classifier assigns severity and category to port changes.
type Classifier struct {
	criticalPorts map[uint16]bool
}

// New returns a Classifier. criticalPorts lists port numbers that should
// be treated as critical when they open unexpectedly.
func New(criticalPorts []uint16) *Classifier {
	m := make(map[uint16]bool, len(criticalPorts))
	for _, p := range criticalPorts {
		m[p] = true
	}
	return &Classifier{criticalPorts: m}
}

// ClassifyOpened returns a Result for a port that has just opened.
func (c *Classifier) ClassifyOpened(p scanner.Port) Result {
	sev := SeverityWarning
	if c.criticalPorts[p.Port] {
		sev = SeverityCritical
	}
	return Result{
		Port:     p,
		Category: CategoryNewPort,
		Severity: sev,
		Label:    "opened",
	}
}

// ClassifyClosed returns a Result for a port that has just closed.
func (c *Classifier) ClassifyClosed(p scanner.Port) Result {
	return Result{
		Port:     p,
		Category: CategoryClosedPort,
		Severity: SeverityInfo,
		Label:    "closed",
	}
}
