package tap_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tap"
)

// TestConcurrentSendIsSafe verifies that concurrent Send calls and sink
// registrations do not trigger the race detector.
func TestConcurrentSendIsSafe(t *testing.T) {
	tr := tap.New(nil)

	var mu sync.Mutex
	var total int

	const workers = 8
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(n int) {
			defer wg.Done()
			tr.Register(func(diffs []scanner.Diff) {
				mu.Lock()
				total += len(diffs)
				mu.Unlock()
			})
			tr.Send([]scanner.Diff{makeDiff(uint16(n+1), "opened")})
		}(i)
	}
	wg.Wait()

	if total == 0 {
		t.Fatal("expected at least one diff to be delivered")
	}
}

// TestPanicInOneSinkDoesNotBlockOthers ensures a panicking sink does not
// prevent subsequent sinks from receiving the diff.
func TestPanicInOneSinkDoesNotBlockOthers(t *testing.T) {
	tr := tap.New(nil)
	tr.Register(func(_ []scanner.Diff) { panic("intentional") })

	received := false
	tr.Register(func(_ []scanner.Diff) { received = true })

	tr.Send([]scanner.Diff{makeDiff(8080, "opened")})

	if !received {
		t.Fatal("second sink should have received diff despite first sink panic")
	}
}
