// Package envelope wraps a scan event with metadata for downstream handlers.
package envelope

import (
	"time"

	"portwatch/internal/scanner"
)

// Priority indicates the urgency of an envelope.
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
)

func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityHigh:
		return "high"
	default:
		return "normal"
	}
}

// Envelope carries a scan result alongside routing metadata.
type Envelope struct {
	ID        string
	CreatedAt time.Time
	Priority  Priority
	Ports     []scanner.Port
	Tags      []string
	Meta      map[string]string
}

// New returns an Envelope with CreatedAt set to now.
func New(id string, ports []scanner.Port, priority Priority) *Envelope {
	return &Envelope{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Priority:  priority,
		Ports:     ports,
		Meta:      make(map[string]string),
	}
}

// WithTag appends a tag and returns the envelope for chaining.
func (e *Envelope) WithTag(tag string) *Envelope {
	e.Tags = append(e.Tags, tag)
	return e
}

// WithMeta sets a metadata key/value and returns the envelope for chaining.
func (e *Envelope) WithMeta(key, value string) *Envelope {
	e.Meta[key] = value
	return e
}

// HasTag reports whether the envelope carries the given tag.
func (e *Envelope) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}
