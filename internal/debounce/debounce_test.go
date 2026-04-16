package debounce

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeDiff(port uint16) scanner.Diff {
	return scanner.Diff{Port: scanner.Port{Port: port, Proto: "tcp"}, Kind: "opened"}
}

func TestFlushDeliversBuffered(t *testing.T) {
	var mu sync.Mutex
	var got []scanner.Diff
	d := New(500*time.Millisecond, func(diffs []scanner.Diff) {
		mu.Lock()
		got = append(got, diffs...)
		mu.Unlock()
	})
	d.Add([]scanner.Diff{makeDiff(80)})
	d.Add([]scanner.Diff{makeDiff(443)})
	d.Flush()
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(got))
	}
}

func TestAddEmptyIsNoop(t *testing.T) {
	called := false
	d := New(10*time.Millisecond, func(_ []scanner.Diff) { called = true })
	d.Add(nil)
	d.Add([]scanner.Diff{})
	time.Sleep(30 * time.Millisecond)
	if called {
		t.Fatal("handler should not be called for empty adds")
	}
}

func TestTimerFiresAfterWindow(t *testing.T) {
	var mu sync.Mutex
	var got []scanner.Diff
	d := New(30*time.Millisecond, func(diffs []scanner.Diff) {
		mu.Lock()
		got = append(got, diffs...)
		mu.Unlock()
	})
	d.Add([]scanner.Diff{makeDiff(22)})
	time.Sleep(80 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(got) == 0 {
		t.Fatal("expected handler to fire after window")
	}
}

func TestFlushOnEmptyBufferIsNoop(t *testing.T) {
	called := false
	d := New(100*time.Millisecond, func(_ []scanner.Diff) { called = true })
	d.Flush()
	if called {
		t.Fatal("flush on empty buffer should not call handler")
	}
}

func TestResetTimerOnConsecutiveAdds(t *testing.T) {
	var mu sync.Mutex
	var calls int
	d := New(50*time.Millisecond, func(_ []scanner.Diff) {
		mu.Lock()
		calls++
		mu.Unlock()
	})
	for i := 0; i < 5; i++ {
		d.Add([]scanner.Diff{makeDiff(uint16(i + 1))})
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected exactly 1 handler call, got %d", calls)
	}
}
