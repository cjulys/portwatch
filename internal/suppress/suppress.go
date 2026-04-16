// Package suppress provides a mechanism to silence alerts for specific ports
// during a defined maintenance window or indefinitely.
package suppress

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Rule defines a suppression rule for a port/protocol pair.
type Rule struct {
	Port     uint16    `json:"port"`
	Protocol string    `json:"protocol"`
	Until    time.Time `json:"until"` // zero means indefinite
}

// Store holds active suppression rules.
type Store struct {
	mu    sync.RWMutex
	rules []Rule
	path  string
}

// New creates a Store backed by the given file path.
func New(path string) *Store {
	s := &Store{path: path}
	_ = s.load()
	return s
}

// Add inserts or replaces a suppression rule.
func (s *Store) Add(r Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, existing := range s.rules {
		if existing.Port == r.Port && existing.Protocol == r.Protocol {
			s.rules[i] = r
			return s.save()
		}
	}
	s.rules = append(s.rules, r)
	return s.save()
}

// Remove deletes the suppression rule for the given port/protocol.
func (s *Store) Remove(port uint16, protocol string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	filtered := s.rules[:0]
	for _, r := range s.rules {
		if r.Port != port || r.Protocol != protocol {
			filtered = append(filtered, r)
		}
	}
	s.rules = filtered
	return s.save()
}

// IsSuppressed reports whether alerts for port/protocol are currently suppressed.
func (s *Store) IsSuppressed(port uint16, protocol string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.rules {
		if r.Port == port && r.Protocol == protocol {
			return r.Until.IsZero() || time.Now().Before(r.Until)
		}
	}
	return false
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.rules, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil
	}
	return json.Unmarshal(data, &s.rules)
}
