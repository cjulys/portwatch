package schedule

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRegisterAndDueImmediately(t *testing.T) {
	base := time.Now()
	s := New()
	s.now = fixedNow(base)

	if err := s.Register("scan", 5*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	due := s.Due()
	if len(due) != 1 || due[0] != "scan" {
		t.Fatalf("expected scan to be due immediately, got %v", due)
	}
}

func TestNotDueBeforeInterval(t *testing.T) {
	base := time.Now()
	s := New()
	s.now = fixedNow(base)
	s.Register("scan", 10*time.Second)
	s.Due() // consume first tick

	// advance less than interval
	s.now = fixedNow(base.Add(5 * time.Second))
	due := s.Due()
	if len(due) != 0 {
		t.Fatalf("expected nothing due, got %v", due)
	}
}

func TestDueAfterInterval(t *testing.T) {
	base := time.Now()
	s := New()
	s.now = fixedNow(base)
	s.Register("scan", 10*time.Second)
	s.Due() // consume first tick

	s.now = fixedNow(base.Add(10 * time.Second))
	due := s.Due()
	if len(due) != 1 || due[0] != "scan" {
		t.Fatalf("expected scan due, got %v", due)
	}
}

func TestRegisterInvalidInterval(t *testing.T) {
	s := New()
	if err := s.Register("bad", 0); err == nil {
		t.Fatal("expected error for zero interval")
	}
	if err := s.Register("bad", -1*time.Second); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestNextReturnsEarliest(t *testing.T) {
	base := time.Now()
	s := New()
	s.now = fixedNow(base)
	s.Register("a", 5*time.Second)
	s.Register("b", 20*time.Second)
	s.Due() // fire both, sets next = base+interval

	nxt := s.Next()
	expected := base.Add(5 * time.Second)
	if nxt != expected {
		t.Fatalf("expected %v, got %v", expected, nxt)
	}
}

func TestRemoveEntry(t *testing.T) {
	base := time.Now()
	s := New()
	s.now = fixedNow(base)
	s.Register("scan", 5*time.Second)
	s.Remove("scan")
	if len(s.Due()) != 0 {
		t.Fatal("expected no entries after remove")
	}
}

func TestNextZeroWhenEmpty(t *testing.T) {
	s := New()
	if !s.Next().IsZero() {
		t.Fatal("expected zero time for empty scheduler")
	}
}
