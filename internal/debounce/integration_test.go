package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/scanner"
)

func TestFlapCancelledWithinWindow(t *testing.T) {
	// Simulate port 8080 closing then reopening within the quiet window.
	// Net result: handler should see both events but only fire once.
	var mu sync.Mutex
	var invocations int
	var collected []scanner.Diff

	d := debounce.New(60*time.Millisecond, func(diffs []scanner.Diff) {
		mu.Lock()
		invocations++
		collected = append(collected, diffs...)
		mu.Unlock()
	})

	d.Add([]scanner.Diff{{Port: scanner.Port{Port: 8080, Proto: "tcp"}, Kind: "closed"}})
	time.Sleep(20 * time.Millisecond)
	d.Add([]scanner.Diff{{Port: scanner.Port{Port: 8080, Proto: "tcp"}, Kind: "opened"}})

	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if invocations != 1 {
		t.Fatalf("expected 1 invocation, got %d", invocations)
	}
	if len(collected) != 2 {
		t.Fatalf("expected 2 diffs collected, got %d", len(collected))
	}
}

func TestExplicitFlushBeforeWindow(t *testing.T) {
	var mu sync.Mutex
	var got []scanner.Diff

	d := debounce.New(500*time.Millisecond, func(diffs []scanner.Diff) {
		mu.Lock()
		got = append(got, diffs...)
		mu.Unlock()
	})

	d.Add([]scanner.Diff{{Port: scanner.Port{Port: 9000, Proto: "udp"}, Kind: "opened"}})
	d.Flush()

	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 diff after explicit flush, got %d", len(got))
	}
}
