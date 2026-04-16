package audit_test

import (
	"path/filepath"
	"testing"

	"portwatch/internal/audit"
)

func TestMultipleWritersAppend(t *testing.T) {
	p := filepath.Join(t.TempDir(), "audit.log")

	l1, err := audit.New(p)
	if err != nil {
		t.Fatalf("New l1: %v", err)
	}
	l2, err := audit.New(p)
	if err != nil {
		t.Fatalf("New l2: %v", err)
	}

	_ = l1.Record("scan", "first")
	_ = l2.Record("scan", "second")
	_ = l1.Record("scan", "third")

	entries, err := audit.ReadAll(p)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestRoundTripPreservesFields(t *testing.T) {
	p := filepath.Join(t.TempDir(), "audit.log")
	l, _ := audit.New(p)
	_ = l.Record("baseline_violation", "9000/udp unexpected")

	entries, _ := audit.ReadAll(p)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	e := entries[0]
	if e.Event != "baseline_violation" {
		t.Errorf("event = %q", e.Event)
	}
	if e.Detail != "9000/udp unexpected" {
		t.Errorf("detail = %q", e.Detail)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp zero")
	}
}
