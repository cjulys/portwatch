package sampler

import (
	"testing"
	"time"
)

func TestInitialIntervalIsMax(t *testing.T) {
	s := New(time.Second, 10*time.Second, 0.5, 1.5)
	if s.Current() != 10*time.Second {
		t.Fatalf("expected 10s, got %v", s.Current())
	}
}

func TestRecordActivityDecreasesInterval(t *testing.T) {
	s := New(time.Second, 10*time.Second, 0.5, 1.5)
	s.RecordActivity()
	if s.Current() >= 10*time.Second {
		t.Fatalf("interval should have decreased, got %v", s.Current())
	}
}

func TestRecordActivityClampsToMin(t *testing.T) {
	s := New(time.Second, 2*time.Second, 0.1, 1.5)
	for i := 0; i < 20; i++ {
		s.RecordActivity()
	}
	if s.Current() < time.Second {
		t.Fatalf("interval dropped below min: %v", s.Current())
	}
	if s.Current() != time.Second {
		t.Fatalf("expected min 1s, got %v", s.Current())
	}
}

func TestRecordQuietIncreasesInterval(t *testing.T) {
	s := New(time.Second, 10*time.Second, 0.5, 1.5)
	s.RecordActivity() // bring it down first
	before := s.Current()
	s.RecordQuiet()
	if s.Current() <= before {
		t.Fatalf("interval should have grown, got %v", s.Current())
	}
}

func TestRecordQuietClampsToMax(t *testing.T) {
	s := New(time.Second, 10*time.Second, 0.5, 2.0)
	for i := 0; i < 20; i++ {
		s.RecordQuiet()
	}
	if s.Current() > 10*time.Second {
		t.Fatalf("interval exceeded max: %v", s.Current())
	}
}

func TestResetRestoresMax(t *testing.T) {
	s := New(time.Second, 10*time.Second, 0.5, 1.5)
	s.RecordActivity()
	s.RecordActivity()
	s.Reset()
	if s.Current() != 10*time.Second {
		t.Fatalf("expected 10s after reset, got %v", s.Current())
	}
}

func TestInvalidStepDownDefaulted(t *testing.T) {
	s := New(time.Second, 10*time.Second, -1, 1.5)
	s.RecordActivity()
	// with default 0.5 step-down, result should be 5s
	if s.Current() != 5*time.Second {
		t.Fatalf("expected 5s, got %v", s.Current())
	}
}

func TestMinGreaterThanMaxClamped(t *testing.T) {
	s := New(20*time.Second, 5*time.Second, 0.5, 1.5)
	if s.min > s.max {
		t.Fatal("min should not exceed max after clamping")
	}
}
