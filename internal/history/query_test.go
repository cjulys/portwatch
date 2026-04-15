package history

import (
	"testing"
	"time"
)

func TestLatestNilOnEmpty(t *testing.T) {
	h, _ := New(tempPath(t), 10)
	if h.Latest() != nil {
		t.Error("expected nil Latest on empty history")
	}
}

func TestLatestReturnsNewest(t *testing.T) {
	h, _ := New(tempPath(t), 10)
	_ = h.Add(makePorts(80))
	_ = h.Add(makePorts(443))
	latest := h.Latest()
	if latest == nil {
		t.Fatal("expected non-nil Latest")
	}
	if len(latest.Ports) != 1 || latest.Ports[0].Number != 443 {
		t.Errorf("unexpected latest port: %+v", latest.Ports)
	}
}

func TestSinceFiltersEntries(t *testing.T) {
	h, _ := New(tempPath(t), 10)
	_ = h.Add(makePorts(22))
	cutoff := time.Now().UTC()
	time.Sleep(2 * time.Millisecond)
	_ = h.Add(makePorts(80))
	_ = h.Add(makePorts(443))

	results := h.Since(cutoff)
	if len(results) != 2 {
		t.Errorf("expected 2 entries since cutoff, got %d", len(results))
	}
}

func TestPortSeenTrue(t *testing.T) {
	h, _ := New(tempPath(t), 10)
	_ = h.Add(makePorts(8080))
	if !h.PortSeen(8080, "tcp") {
		t.Error("expected PortSeen to return true for port 8080/tcp")
	}
}

func TestPortSeenFalse(t *testing.T) {
	h, _ := New(tempPath(t), 10)
	_ = h.Add(makePorts(80))
	if h.PortSeen(9999, "tcp") {
		t.Error("expected PortSeen to return false for port 9999/tcp")
	}
}

func TestUniquePortsEver(t *testing.T) {
	h, _ := New(tempPath(t), 10)
	_ = h.Add(makePorts(80, 443))
	_ = h.Add(makePorts(80, 8080))

	unique := h.UniquePortsEver()
	if len(unique) != 3 {
		t.Errorf("expected 3 unique ports, got %d", len(unique))
	}
}
