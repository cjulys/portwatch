package rollup_test

import (
	"sync"
	"testing"
	"time"

	"portwatch/internal/rollup"
	"portwatch/internal/scanner"
)

func TestMultipleAddsMergeIntoOneEvent(t *testing.T) {
	var mu sync.Mutex
	var events []rollup.Event

	r := rollup.New(30*time.Millisecond, func(e rollup.Event) {
		mu.Lock()
		events = append(events, e)
		mu.Unlock()
	})

	for i := uint16(1); i <= 5; i++ {
		r.Add([]scanner.Port{{Port: i, Protocol: "tcp"}}, nil)
	}

	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) != 1 {
		t.Fatalf("expected 1 rolled-up event, got %d", len(events))
	}
	if len(events[0].Opened) != 5 {
		t.Errorf("expected 5 opened ports, got %d", len(events[0].Opened))
	}
}

func TestFlushResetsState(t *testing.T) {
	count := 0
	r := rollup.New(10*time.Second, func(e rollup.Event) { count++ })

	r.Add([]scanner.Port{{Port: 80, Protocol: "tcp"}}, nil)
	r.Flush()
	r.Flush() // second flush with nothing queued should not fire

	if count != 1 {
		t.Errorf("expected exactly 1 event, got %d", count)
	}
}
