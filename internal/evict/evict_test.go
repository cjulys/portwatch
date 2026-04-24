package evict

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestTouchAndLenIncrements(t *testing.T) {
	tr := New(time.Minute)
	tr.Touch(makePort("tcp", 80))
	tr.Touch(makePort("tcp", 443))
	if tr.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", tr.Len())
	}
}

func TestEvictRemovesExpiredEntry(t *testing.T) {
	tr := New(time.Second)
	now := time.Now()
	tr.now = func() time.Time { return now }

	tr.Touch(makePort("tcp", 8080))

	// Advance clock past TTL.
	tr.now = func() time.Time { return now.Add(2 * time.Second) }

	evicted := tr.Evict()
	if len(evicted) != 1 {
		t.Fatalf("expected 1 evicted port, got %d", len(evicted))
	}
	if evicted[0].Number != 8080 {
		t.Errorf("expected port 8080, got %d", evicted[0].Number)
	}
	if tr.Len() != 0 {
		t.Errorf("expected tracker to be empty after eviction")
	}
}

func TestEvictKeepsFreshEntry(t *testing.T) {
	tr := New(time.Minute)
	now := time.Now()
	tr.now = func() time.Time { return now }

	tr.Touch(makePort("udp", 53))

	// Advance clock but stay within TTL.
	tr.now = func() time.Time { return now.Add(30 * time.Second) }

	evicted := tr.Evict()
	if len(evicted) != 0 {
		t.Errorf("expected no evictions, got %d", len(evicted))
	}
	if tr.Len() != 1 {
		t.Errorf("expected 1 tracked entry, got %d", tr.Len())
	}
}

func TestEvictMixedExpiry(t *testing.T) {
	tr := New(10 * time.Second)
	start := time.Now()
	tr.now = func() time.Time { return start }

	tr.Touch(makePort("tcp", 22))

	// Advance and touch a second port (fresh).
	tr.now = func() time.Time { return start.Add(5 * time.Second) }
	tr.Touch(makePort("tcp", 443))

	// Advance past TTL of first port only.
	tr.now = func() time.Time { return start.Add(15 * time.Second) }

	evicted := tr.Evict()
	if len(evicted) != 1 {
		t.Fatalf("expected 1 eviction, got %d", len(evicted))
	}
	if evicted[0].Number != 22 {
		t.Errorf("expected port 22 evicted, got %d", evicted[0].Number)
	}
}

func TestEvictEmptyTrackerReturnsNil(t *testing.T) {
	tr := New(time.Second)
	evicted := tr.Evict()
	if evicted != nil {
		t.Errorf("expected nil on empty tracker, got %v", evicted)
	}
}

func TestTouchRefreshesTTL(t *testing.T) {
	tr := New(10 * time.Second)
	start := time.Now()
	tr.now = func() time.Time { return start }

	tr.Touch(makePort("tcp", 80))

	// Re-touch just before expiry.
	tr.now = func() time.Time { return start.Add(9 * time.Second) }
	tr.Touch(makePort("tcp", 80))

	// Advance past original TTL but within refreshed TTL.
	tr.now = func() time.Time { return start.Add(15 * time.Second) }

	evicted := tr.Evict()
	if len(evicted) != 0 {
		t.Errorf("expected no evictions after refresh, got %d", len(evicted))
	}
}
