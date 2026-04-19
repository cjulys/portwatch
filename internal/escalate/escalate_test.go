package escalate

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestFirstCallIsInfo(t *testing.T) {
	e := New(time.Minute, 2, 4)
	if got := e.Evaluate("k"); got != Info {
		t.Fatalf("want Info got %s", got)
	}
}

func TestEscalatesToWarning(t *testing.T) {
	e := New(time.Minute, 2, 4)
	e.Evaluate("k")
	if got := e.Evaluate("k"); got != Warning {
		t.Fatalf("want Warning got %s", got)
	}
}

func TestEscalatesToCritical(t *testing.T) {
	e := New(time.Minute, 2, 4)
	for i := 0; i < 3; i++ {
		e.Evaluate("k")
	}
	if got := e.Evaluate("k"); got != Critical {
		t.Fatalf("want Critical got %s", got)
	}
}

func TestWindowExpiryResetsCount(t *testing.T) {
	base := time.Now()
	e := New(time.Minute, 2, 4)
	e.now = fixedClock(base)
	e.Evaluate("k")
	e.Evaluate("k") // Warning
	// advance past window
	e.now = fixedClock(base.Add(2 * time.Minute))
	if got := e.Evaluate("k"); got != Info {
		t.Fatalf("want Info after window reset, got %s", got)
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	e := New(time.Minute, 2, 4)
	e.Evaluate("a")
	e.Evaluate("a") // Warning for a
	if got := e.Evaluate("b"); got != Info {
		t.Fatalf("want Info for new key b, got %s", got)
	}
}

func TestResetClearsKey(t *testing.T) {
	e := New(time.Minute, 2, 4)
	e.Evaluate("k")
	e.Evaluate("k")
	e.Reset("k")
	if got := e.Evaluate("k"); got != Info {
		t.Fatalf("want Info after reset, got %s", got)
	}
}

func TestFlushClearsAll(t *testing.T) {
	e := New(time.Minute, 2, 4)
	e.Evaluate("a")
	e.Evaluate("b")
	e.Flush()
	if got := e.Evaluate("a"); got != Info {
		t.Fatalf("want Info after flush, got %s", got)
	}
}

func TestLevelStrings(t *testing.T) {
	if Info.String() != "info" || Warning.String() != "warning" || Critical.String() != "critical" {
		t.Fatal("unexpected level string")
	}
}
