package correlate

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeDiff(port uint16, dir string) scanner.Diff {
	return scanner.Diff{
		Port:      scanner.Port{Port: port, Protocol: "tcp"},
		Direction: dir,
	}
}

func TestFlushEmitsEvent(t *testing.T) {
	var mu sync.Mutex
	var got []Event

	c := New(100*time.Millisecond, func(ev Event) {
		mu.Lock()
		got = append(got, ev)
		mu.Unlock()
	})

	c.Add([]scanner.Diff{makeDiff(80, "opened")})
	c.Flush()

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if len(got[0].Diffs) != 1 {
		t.Errorf("expected 1 diff in event, got %d", len(got[0].Diffs))
	}
}

func TestEventIDIsNonEmpty(t *testing.T) {
	var got Event
	c := New(50*time.Millisecond, func(ev Event) { got = ev })
	c.Add([]scanner.Diff{makeDiff(443, "opened")})
	c.Flush()
	if got.ID == "" {
		t.Error("expected non-empty correlation ID")
	}
}

func TestMultipleAddsGroupedIntoOneEvent(t *testing.T) {
	var mu sync.Mutex
	var events []Event

	c := New(200*time.Millisecond, func(ev Event) {
		mu.Lock()
		events = append(events, ev)
		mu.Unlock()
	})

	c.Add([]scanner.Diff{makeDiff(22, "opened")})
	c.Add([]scanner.Diff{makeDiff(80, "opened")})
	c.Add([]scanner.Diff{makeDiff(443, "opened")})
	c.Flush()

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 1 {
		t.Fatalf("expected 1 correlated event, got %d", len(events))
	}
	if len(events[0].Diffs) != 3 {
		t.Errorf("expected 3 diffs, got %d", len(events[0].Diffs))
	}
}

func TestFlushOnEmptyBufferIsNoop(t *testing.T) {
	called := false
	c := New(50*time.Millisecond, func(_ Event) { called = true })
	c.Flush()
	if called {
		t.Error("flush on empty buffer should not invoke callback")
	}
}

func TestAutoFlushFiresAfterWindow(t *testing.T) {
	var mu sync.Mutex
	var events []Event

	c := New(60*time.Millisecond, func(ev Event) {
		mu.Lock()
		events = append(events, ev)
		mu.Unlock()
	})

	c.Add([]scanner.Diff{makeDiff(8080, "opened")})
	time.Sleep(150 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) == 0 {
		t.Error("expected auto-flush to fire after window")
	}
}

func TestCreatedAtIsUTC(t *testing.T) {
	var got Event
	c := New(50*time.Millisecond, func(ev Event) { got = ev })
	c.Add([]scanner.Diff{makeDiff(9000, "closed")})
	c.Flush()
	if got.CreatedAt.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", got.CreatedAt.Location())
	}
}
