package circuitbreaker

import (
	"bytes"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestClosedByDefault(t *testing.T) {
	cb := New(3, time.Second)
	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected closed, got %s", cb.CurrentState())
	}
}

func TestAllowUnderThreshold(t *testing.T) {
	cb := New(3, time.Second)
	cb.RecordFailure()
	cb.RecordFailure()
	if !cb.Allow() {
		t.Fatal("expected Allow=true below threshold")
	}
}

func TestOpensAfterThreshold(t *testing.T) {
	cb := New(3, time.Second)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.CurrentState() != StateOpen {
		t.Fatalf("expected open, got %s", cb.CurrentState())
	}
	if cb.Allow() {
		t.Fatal("expected Allow=false when open")
	}
}

func TestHalfOpenAfterReset(t *testing.T) {
	now := time.Now()
	cb := New(2, 10*time.Second)
	cb.now = fixedClock(now)
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.CurrentState() != StateOpen {
		t.Fatal("expected open")
	}
	cb.now = fixedClock(now.Add(11 * time.Second))
	if !cb.Allow() {
		t.Fatal("expected Allow=true after reset window")
	}
	if cb.CurrentState() != StateHalfOpen {
		t.Fatalf("expected half-open, got %s", cb.CurrentState())
	}
}

func TestSuccessClosesFromHalfOpen(t *testing.T) {
	now := time.Now()
	cb := New(1, 5*time.Second)
	cb.now = fixedClock(now)
	cb.RecordFailure()
	cb.now = fixedClock(now.Add(6 * time.Second))
	cb.Allow() // transitions to half-open
	cb.RecordSuccess()
	if cb.CurrentState() != StateClosed {
		t.Fatalf("expected closed after success, got %s", cb.CurrentState())
	}
}

func TestFallbackWriterOnOpen(t *testing.T) {
	var buf bytes.Buffer
	cb := New(1, time.Second)
	cb.fallback = &buf
	cb.RecordFailure()
	if buf.Len() == 0 {
		t.Fatal("expected fallback message written on open")
	}
}

func TestStateString(t *testing.T) {
	for _, tc := range []struct {
		s    State
		want string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	} {
		if tc.s.String() != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", tc.s, tc.s.String(), tc.want)
		}
	}
}
