package window

import (
	"testing"
	"time"
)

func TestCountEmptyIsZero(t *testing.T) {
	w := New(time.Minute)
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAddAndCount(t *testing.T) {
	w := New(time.Minute)
	w.Add(3)
	w.Add(2)
	if got := w.Count(); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestEvictsOldEntries(t *testing.T) {
	now := time.Now()
	clock := &now
	w := newWithClock(time.Second, func() time.Time { return *clock })

	w.Add(10)

	// Advance clock beyond the window span.
	future := now.Add(2 * time.Second)
	clock = &future
	w.Add(1)

	if got := w.Count(); got != 1 {
		t.Fatalf("expected 1 after eviction, got %d", got)
	}
}

func TestResetClearsAll(t *testing.T) {
	w := New(time.Minute)
	w.Add(5)
	w.Reset()
	if got := w.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestEntriesWithinWindowAreKept(t *testing.T) {
	now := time.Now()
	clock := now
	w := newWithClock(5*time.Second, func() time.Time { return clock })

	w.Add(4)
	clock = now.Add(3 * time.Second)
	w.Add(6)

	// Both entries are within the 5s window.
	if got := w.Count(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestOnlyExpiredEntryEvicted(t *testing.T) {
	now := time.Now()
	clock := now
	w := newWithClock(5*time.Second, func() time.Time { return clock })

	w.Add(7)
	clock = now.Add(6 * time.Second) // first entry now expired
	w.Add(3)

	if got := w.Count(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}
