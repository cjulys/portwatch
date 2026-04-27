package label

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/internal/scanner"
)

// Enricher wraps a Labeler and stamps the label as a tag on each port
// using the format "label:<value>".
type Enricher struct {
	labeler *Labeler
}

// NewEnricher returns an Enricher backed by the given Labeler.
func NewEnricher(l *Labeler) *Enricher {
	return &Enricher{labeler: l}
}

// Enrich returns a copy of p with the label tag appended when a label is
// found. If no label exists the port is returned unchanged.
func (e *Enricher) Enrich(p scanner.Port) scanner.Port {
	lbl := e.labeler.Apply(p)
	if lbl == "" {
		return p
	}
	tag := fmt.Sprintf("label:%s", strings.ToLower(lbl))
	for _, t := range p.Tags {
		if t == tag {
			return p // already tagged
		}
	}
	out := p
	out.Tags = make([]string, len(p.Tags)+1)
	copy(out.Tags, p.Tags)
	out.Tags[len(p.Tags)] = tag
	return out
}

// EnrichAll enriches every port in the slice.
func (e *Enricher) EnrichAll(ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, len(ports))
	for i, p := range ports {
		out[i] = e.Enrich(p)
	}
	return out
}
