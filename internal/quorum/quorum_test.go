package quorum_test

import (
	"testing"

	"github.com/user/portwatch/internal/quorum"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Proto: proto}
}

func TestThresholdOneFiresImmediately(t *testing.T) {
	q := quorum.New(1)
	p := makePort(80, "tcp")
	if !q.Observe(p, "opened") {
		t.Fatal("expected true on first observation with threshold=1")
	}
}

func TestThresholdTwoRequiresTwoCalls(t *testing.T) {
	q := quorum.New(2)
	p := makePort(443, "tcp")

	if q.Observe(p, "opened") {
		t.Fatal("should not fire on first observation")
	}
	if !q.Observe(p, "opened") {
		t.Fatal("should fire on second observation")
	}
	// third call: threshold already passed, should not fire again
	if q.Observe(p, "opened") {
		t.Fatal("should not fire again after threshold reached")
	}
}

func TestDirectionsAreIndependent(t *testing.T) {
	q := quorum.New(2)
	p := makePort(22, "tcp")

	q.Observe(p, "opened") // count=1 for opened
	if q.Observe(p, "closed") {
		t.Fatal("closed direction should not fire after one observation")
	}
}

func TestResetAllowsRecount(t *testing.T) {
	q := quorum.New(2)
	p := makePort(8080, "tcp")

	q.Observe(p, "opened") // count=1
	q.Reset(p, "opened")
	q.Observe(p, "opened") // count=1 again after reset
	if q.Observe(p, "opened") {
		// count=2 after reset, should fire
		// (this branch means it fired on 2nd post-reset call — correct)
		return
	}
	t.Fatal("expected threshold to fire on second post-reset observation")
}

func TestFlushClearsAllKeys(t *testing.T) {
	q := quorum.New(2)
	p1 := makePort(80, "tcp")
	p2 := makePort(443, "tcp")

	q.Observe(p1, "opened")
	q.Observe(p2, "closed")
	q.Flush()

	// After flush both counts are zero; second observation should not fire.
	q.Observe(p1, "opened")
	if q.Observe(p1, "opened") {
		// fired on 2nd call post-flush — correct for threshold=2
		return
	}
	t.Fatal("expected threshold to fire on second post-flush observation")
}

func TestBelowOneThresholdClamped(t *testing.T) {
	q := quorum.New(0) // clamped to 1
	p := makePort(53, "udp")
	if !q.Observe(p, "opened") {
		t.Fatal("clamped threshold=1 should fire immediately")
	}
}

func TestProtocolDistinction(t *testing.T) {
	q := quorum.New(2)
	tcp := makePort(53, "tcp")
	udp := makePort(53, "udp")

	q.Observe(tcp, "opened") // tcp count=1
	if q.Observe(udp, "opened") {
		t.Fatal("udp and tcp should be tracked independently")
	}
}
