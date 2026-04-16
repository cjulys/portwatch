package schedule_test

import (
	"testing"
	"time"

	"portwatch/internal/schedule"
)

// TestMultipleEntriesFireIndependently verifies that two entries with
// different intervals fire at the correct times relative to each other.
func TestMultipleEntriesFireIndependently(t *testing.T) {
	s := schedule.New()

	if err := s.Register("fast", 2*time.Millisecond); err != nil {
		t.Fatal(err)
	}
	if err := s.Register("slow", 50*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// Both due immediately on first call.
	due := s.Due()
	if len(due) != 2 {
		t.Fatalf("expected 2 due on first call, got %d", len(due))
	}

	// Wait for fast but not slow.
	time.Sleep(3 * time.Millisecond)
	due = s.Due()
	if len(due) != 1 || due[0] != "fast" {
		t.Fatalf("expected only 'fast' due, got %v", due)
	}

	// Wait for slow.
	time.Sleep(55 * time.Millisecond)
	due = s.Due()
	found := map[string]bool{}
	for _, d := range due {
		found[d] = true
	}
	if !found["slow"] {
		t.Fatalf("expected 'slow' to be due, got %v", due)
	}
}
