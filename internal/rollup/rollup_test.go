package rollup

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func port(p uint16, proto string) scanner.Port {
	return scanner.Port{Port: p, Protocol: proto}
}

func TestFlushEmitsAccumulatedDiffs(t *testing.T) {
	var mu sync.Mutex
	var got []Event
	r := New(10*time.Second, func(e Event) {
		mu.Lock()
		got = append(got, e)
		mu.Unlock()
	})
	r.Add([]scanner.Port{port(80, "tcp")}, nil)
	r.Add(nil, []scanner.Port{port(443, "tcp")})
	r.Flush()
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 event, got %d", len(got))
	}
	if len(got[0].Opened) != 1 || got[0].Opened[0].Port != 80 {
		t.Errorf("unexpected opened: %v", got[0].Opened)
	}
	if len(got[0].Closed) != 1 || got[0].Closed[0].Port != 443 {
		t.Errorf("unexpected closed: %v", got[0].Closed)
	}
}

func TestNoFlushWhenEmpty(t *testing.T) {
	called := false
	r := New(10*time.Second, func(e Event) { called = true })
	r.Flush()
	if called {
		t.Error("flush fn should not be called when nothing was added")
	}
}

func TestWindowFiresAutomatically(t *testing.T) {
	ch := make(chan Event, 1)
	r := New(50*time.Millisecond, func(e Event) { ch <- e })
	r.Add([]scanner.Port{port(22, "tcp")}, nil)
	select {
	case e := <-ch:
		if len(e.Opened) != 1 {
			t.Errorf("expected 1 opened port, got %d", len(e.Opened))
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for automatic flush")
	}
}

func TestEventTimestampIsUTC(t *testing.T) {
	ch := make(chan Event, 1)
	r := New(10*time.Second, func(e Event) { ch <- e })
	r.Add([]scanner.Port{port(8080, "tcp")}, nil)
	r.Flush()
	e := <-ch
	if e.At.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %v", e.At.Location())
	}
}
