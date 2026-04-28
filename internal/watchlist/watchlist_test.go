package watchlist

import (
	"testing"
)

func TestAddAndContains(t *testing.T) {
	w := New()
	w.Add(22, "tcp")
	if !w.Contains(22, "tcp") {
		t.Fatal("expected port 22/tcp to be on the watchlist")
	}
}

func TestContainsReturnsFalseForUnknownPort(t *testing.T) {
	w := New()
	if w.Contains(80, "tcp") {
		t.Fatal("expected empty watchlist to return false")
	}
}

func TestRemoveDeletesEntry(t *testing.T) {
	w := New()
	w.Add(443, "tcp")
	w.Remove(443, "tcp")
	if w.Contains(443, "tcp") {
		t.Fatal("expected port to be removed")
	}
}

func TestRemoveNonExistentIsNoop(t *testing.T) {
	w := New()
	w.Remove(9999, "tcp") // must not panic
}

func TestProtocolDistinction(t *testing.T) {
	w := New()
	w.Add(53, "tcp")
	if w.Contains(53, "udp") {
		t.Fatal("tcp and udp entries must be independent")
	}
}

func TestDuplicateAddDoesNotGrow(t *testing.T) {
	w := New()
	w.Add(22, "tcp")
	w.Add(22, "tcp")
	if w.Len() != 1 {
		t.Fatalf("expected len 1, got %d", w.Len())
	}
}

func TestAllReturnsSnapshot(t *testing.T) {
	w := New()
	w.Add(22, "tcp")
	w.Add(80, "tcp")
	w.Add(53, "udp")
	all := w.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
}

func TestLenReflectsCurrentState(t *testing.T) {
	w := New()
	if w.Len() != 0 {
		t.Fatal("expected empty watchlist")
	}
	w.Add(8080, "tcp")
	if w.Len() != 1 {
		t.Fatal("expected len 1 after add")
	}
	w.Remove(8080, "tcp")
	if w.Len() != 0 {
		t.Fatal("expected len 0 after remove")
	}
}
