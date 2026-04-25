package reachable

import (
	"testing"
	"time"
)

func TestGetEmptyReturnsZeroScore(t *testing.T) {
	tr := New(time.Minute)
	s := tr.Get("tcp:80")
	if s.Total != 0 || s.Seen != 0 || s.Ratio != 0 {
		t.Fatalf("expected zero score, got %+v", s)
	}
}

func TestRecordOpenIncrementsSeen(t *testing.T) {
	tr := New(time.Minute)
	tr.Record("tcp:80", true)
	tr.Record("tcp:80", true)
	s := tr.Get("tcp:80")
	if s.Seen != 2 || s.Total != 2 {
		t.Fatalf("expected seen=2 total=2, got %+v", s)
	}
	if s.Ratio != 1.0 {
		t.Fatalf("expected ratio=1.0, got %v", s.Ratio)
	}
}

func TestRecordMixedComputesRatio(t *testing.T) {
	tr := New(time.Minute)
	tr.Record("tcp:443", true)
	tr.Record("tcp:443", false)
	tr.Record("tcp:443", true)
	tr.Record("tcp:443", false)
	s := tr.Get("tcp:443")
	if s.Total != 4 || s.Seen != 2 {
		t.Fatalf("unexpected counts: %+v", s)
	}
	if s.Ratio != 0.5 {
		t.Fatalf("expected ratio=0.5, got %v", s.Ratio)
	}
}

func TestWindowEvictsOldEntries(t *testing.T) {
	tr := New(50 * time.Millisecond)
	base := time.Now()
	clock := base
	tr.now = func() time.Time { return clock }

	tr.Record("tcp:22", true) // t=0 — will be evicted

	clock = base.Add(60 * time.Millisecond)
	tr.Record("tcp:22", false) // t=60ms — within window from new now

	s := tr.Get("tcp:22")
	if s.Total != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", s.Total)
	}
	if s.Seen != 0 {
		t.Fatalf("expected seen=0, got %d", s.Seen)
	}
}

func TestFlushClearsAllKeys(t *testing.T) {
	tr := New(time.Minute)
	tr.Record("tcp:80", true)
	tr.Record("udp:53", true)
	tr.Flush()
	if s := tr.Get("tcp:80"); s.Total != 0 {
		t.Fatalf("expected empty after flush, got %+v", s)
	}
	if s := tr.Get("udp:53"); s.Total != 0 {
		t.Fatalf("expected empty after flush, got %+v", s)
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	tr := New(time.Minute)
	tr.Record("tcp:80", true)
	tr.Record("tcp:80", true)
	tr.Record("udp:53", false)

	http := tr.Get("tcp:80")
	dns := tr.Get("udp:53")

	if http.Seen != 2 {
		t.Fatalf("tcp:80 seen=%d, want 2", http.Seen)
	}
	if dns.Seen != 0 {
		t.Fatalf("udp:53 seen=%d, want 0", dns.Seen)
	}
}
