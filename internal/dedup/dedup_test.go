package dedup

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Proto: proto, State: "open"}
}

func TestFirstCallIsNotDuplicate(t *testing.T) {
	d := New(10 * time.Second)
	p := makePort(80, "tcp")
	if d.IsDuplicate(p, "opened") {
		t.Fatal("first call should never be a duplicate")
	}
}

func TestSecondCallWithinWindowIsDuplicate(t *testing.T) {
	d := New(10 * time.Second)
	p := makePort(443, "tcp")
	d.IsDuplicate(p, "opened")
	if !d.IsDuplicate(p, "opened") {
		t.Fatal("second call within window should be duplicate")
	}
}

func TestCallAfterWindowIsNotDuplicate(t *testing.T) {
	now := time.Now()
	d := New(5 * time.Second)
	d.now = func() time.Time { return now }

	p := makePort(22, "tcp")
	d.IsDuplicate(p, "opened")

	d.now = func() time.Time { return now.Add(6 * time.Second) }
	if d.IsDuplicate(p, "opened") {
		t.Fatal("call after window expiry should not be duplicate")
	}
}

func TestDifferentDirectionsAreIndependent(t *testing.T) {
	d := New(10 * time.Second)
	p := makePort(8080, "tcp")
	d.IsDuplicate(p, "opened")
	if d.IsDuplicate(p, "closed") {
		t.Fatal("different directions should be tracked independently")
	}
}

func TestDifferentPortsAreIndependent(t *testing.T) {
	d := New(10 * time.Second)
	p1 := makePort(80, "tcp")
	p2 := makePort(81, "tcp")
	d.IsDuplicate(p1, "opened")
	if d.IsDuplicate(p2, "opened") {
		t.Fatal("different ports should be tracked independently")
	}
}

func TestFlushResetsState(t *testing.T) {
	d := New(10 * time.Second)
	p := makePort(53, "udp")
	d.IsDuplicate(p, "opened")
	d.Flush()
	if d.IsDuplicate(p, "opened") {
		t.Fatal("after Flush, call should not be duplicate")
	}
}

func TestStatsCountsRepeats(t *testing.T) {
	d := New(30 * time.Second)
	p := makePort(3306, "tcp")
	d.IsDuplicate(p, "opened") // count = 1
	d.IsDuplicate(p, "opened") // count = 2
	d.IsDuplicate(p, "opened") // count = 3
	if got := d.Stats(p, "opened"); got != 3 {
		t.Fatalf("expected count 3, got %d", got)
	}
}

func TestStatsUnknownKeyIsZero(t *testing.T) {
	d := New(10 * time.Second)
	p := makePort(9999, "tcp")
	if got := d.Stats(p, "opened"); got != 0 {
		t.Fatalf("expected 0 for unseen key, got %d", got)
	}
}
