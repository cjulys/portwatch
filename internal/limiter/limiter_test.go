package limiter

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowUnderLimit(t *testing.T) {
	l := New(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !l.Allow("k") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllowBlocksOverLimit(t *testing.T) {
	l := New(time.Minute, 2)
	l.Allow("k")
	l.Allow("k")
	if l.Allow("k") {
		t.Fatal("expected block on third call")
	}
}

func TestAllowAfterWindowExpires(t *testing.T) {
	base := time.Now()
	l := New(time.Minute, 1)
	l.now = fixedClock(base)
	l.Allow("k")

	l.now = fixedClock(base.Add(61 * time.Second))
	if !l.Allow("k") {
		t.Fatal("expected allow after window expired")
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	l := New(time.Minute, 1)
	l.Allow("a")
	if !l.Allow("b") {
		t.Fatal("key b should be independent of key a")
	}
}

func TestCountReflectsWindow(t *testing.T) {
	base := time.Now()
	l := New(time.Minute, 10)
	l.now = fixedClock(base)
	l.Allow("k")
	l.Allow("k")

	l.now = fixedClock(base.Add(61 * time.Second))
	if c := l.Count("k"); c != 0 {
		t.Fatalf("expected 0 after window, got %d", c)
	}
}

func TestResetClearsAll(t *testing.T) {
	l := New(time.Minute, 1)
	l.Allow("k")
	l.Reset()
	if !l.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}
