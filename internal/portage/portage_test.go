package portage

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: number, Address: "127.0.0.1", State: "open"}
}

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestUpdateAddsNewPort(t *testing.T) {
	tr := New()
	now := time.Now()
	tr.clock = fixedClock(now)

	ports := []scanner.Port{makePort("tcp", 80)}
	tr.Update(ports)

	e, ok := tr.Get(makePort("tcp", 80))
	if !ok {
		t.Fatal("expected port 80 to be tracked")
	}
	if !e.FirstSeen.Equal(now) {
		t.Fatalf("expected FirstSeen=%v, got %v", now, e.FirstSeen)
	}
}

func TestUpdateRemovesClosedPort(t *testing.T) {
	tr := New()
	tr.Update([]scanner.Port{makePort("tcp", 80)})
	tr.Update([]scanner.Port{})

	if tr.Len() != 0 {
		t.Fatalf("expected 0 entries after removal, got %d", tr.Len())
	}
}

func TestFirstSeenNotOverwrittenOnSubsequentUpdate(t *testing.T) {
	tr := New()
	t0 := time.Now()
	tr.clock = fixedClock(t0)
	tr.Update([]scanner.Port{makePort("tcp", 443)})

	// Advance clock and update again with the same port.
	tr.clock = fixedClock(t0.Add(10 * time.Second))
	tr.Update([]scanner.Port{makePort("tcp", 443)})

	e, ok := tr.Get(makePort("tcp", 443))
	if !ok {
		t.Fatal("expected port 443 to still be tracked")
	}
	if !e.FirstSeen.Equal(t0) {
		t.Fatalf("FirstSeen should not change on re-update: got %v want %v", e.FirstSeen, t0)
	}
}

func TestAgeCalculation(t *testing.T) {
	tr := New()
	t0 := time.Now()
	tr.clock = fixedClock(t0)
	tr.Update([]scanner.Port{makePort("tcp", 22)})

	e, _ := tr.Get(makePort("tcp", 22))
	age := e.Age(t0.Add(5 * time.Minute))
	if age != 5*time.Minute {
		t.Fatalf("expected 5m age, got %v", age)
	}
}

func TestProtocolDistinction(t *testing.T) {
	tr := New()
	tr.Update([]scanner.Port{makePort("tcp", 53), makePort("udp", 53)})

	if tr.Len() != 2 {
		t.Fatalf("expected 2 entries (tcp+udp), got %d", tr.Len())
	}
}

func TestSnapshotLength(t *testing.T) {
	tr := New()
	tr.Update([]scanner.Port{makePort("tcp", 80), makePort("tcp", 443), makePort("tcp", 8080)})

	snap := tr.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected snapshot length 3, got %d", len(snap))
	}
}

func TestGetUnknownPortReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get(makePort("tcp", 9999))
	if ok {
		t.Fatal("expected false for untracked port")
	}
}
