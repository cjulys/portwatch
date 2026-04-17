package cooldown

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestFirstCallAlwaysAllowed(t *testing.T) {
	c := New(5 * time.Second)
	if !c.Record("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestSecondCallWithinWindowSuppressed(t *testing.T) {
	base := time.Now()
	c := New(10 * time.Second)
	c.now = fixedClock(base)
	c.Record("k")
	c.now = fixedClock(base.Add(3 * time.Second))
	if c.Record("k") {
		t.Fatal("expected second call within window to be suppressed")
	}
}

func TestCallAfterWindowAllowed(t *testing.T) {
	base := time.Now()
	c := New(5 * time.Second)
	c.now = fixedClock(base)
	c.Record("k")
	c.now = fixedClock(base.Add(6 * time.Second))
	if !c.Record("k") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestWindowResetsOnEachEvent(t *testing.T) {
	base := time.Now()
	c := New(10 * time.Second)
	c.now = fixedClock(base)
	c.Record("k")
	// second event at t+4 — still within window, resets timer
	c.now = fixedClock(base.Add(4 * time.Second))
	c.Record("k")
	// third event at t+9 — only 5s since last, still suppressed
	c.now = fixedClock(base.Add(9 * time.Second))
	if c.Record("k") {
		t.Fatal("expected suppression because window was reset")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	base := time.Now()
	c := New(10 * time.Second)
	c.now = fixedClock(base)
	c.Record("a")
	c.Record("b")
	c.now = fixedClock(base.Add(2 * time.Second))
	if c.Record("a") {
		t.Fatal("a should be suppressed")
	}
	if c.Record("b") {
		t.Fatal("b should be suppressed")
	}
}

func TestFlushResetsState(t *testing.T) {
	base := time.Now()
	c := New(10 * time.Second)
	c.now = fixedClock(base)
	c.Record("k")
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected 0 keys after flush, got %d", c.Len())
	}
	if !c.Record("k") {
		t.Fatal("expected first call after flush to be allowed")
	}
}

func TestLenTracksKeys(t *testing.T) {
	c := New(time.Minute)
	c.Record("x")
	c.Record("y")
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
