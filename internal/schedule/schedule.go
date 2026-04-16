package schedule

import (
	"errors"
	"time"
)

// Entry represents a named scheduled job with a fixed interval.
type Entry struct {
	Name     string
	Interval time.Duration
	Next     time.Time
}

// Scheduler manages multiple named intervals and reports which are due.
type Scheduler struct {
	entries map[string]*Entry
	now     func() time.Time
}

// New returns a Scheduler using real wall time.
func New() *Scheduler {
	return &Scheduler{
		entries: make(map[string]*Entry),
		now:     time.Now,
	}
}

// Register adds or replaces a named entry with the given interval.
// The first tick is due immediately.
func (s *Scheduler) Register(name string, interval time.Duration) error {
	if interval <= 0 {
		return errors.New("schedule: interval must be positive")
	}
	s.entries[name] = &Entry{
		Name:     name,
		Interval: interval,
		Next:     s.now(),
	}
	return nil
}

// Due returns the names of all entries whose next fire time is <= now
// and advances their next scheduled time.
func (s *Scheduler) Due() []string {
	now := s.now()
	var due []string
	for _, e := range s.entries {
		if !now.Before(e.Next) {
			due = append(due, e.Name)
			e.Next = now.Add(e.Interval)
		}
	}
	return due
}

// Next returns the earliest next fire time across all entries.
// Returns zero time if no entries are registered.
func (s *Scheduler) Next() time.Time {
	var earliest time.Time
	for _, e := range s.entries {
		if earliest.IsZero() || e.Next.Before(earliest) {
			earliest = e.Next
		}
	}
	return earliest
}

// Remove deletes a named entry.
func (s *Scheduler) Remove(name string) {
	delete(s.entries, name)
}
