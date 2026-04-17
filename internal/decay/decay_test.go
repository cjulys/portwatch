package decay

import (
	"testing"
	"time"
)

var testSteps = []time.Duration{
	1 * time.Minute,
	10 * time.Minute,
	1 * time.Hour,
}

func TestFirstEvaluationIsCritical(t *testing.T) {
	tr := New(testSteps)
	if got := tr.Evaluate("tcp:8080"); got != LevelCritical {
		t.Fatalf("expected critical, got %s", got)
	}
}

func TestLevelDecaysOverTime(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr := New(testSteps)
	tr.clock = func() time.Time { return now }

	tr.Evaluate("tcp:9000") // register

	tr.clock = func() time.Time { return now.Add(2 * time.Minute) }
	if got := tr.Evaluate("tcp:9000"); got != LevelWarning {
		t.Fatalf("expected warning after 2m, got %s", got)
	}

	tr.clock = func() time.Time { return now.Add(15 * time.Minute) }
	if got := tr.Evaluate("tcp:9000"); got != LevelInfo {
		t.Fatalf("expected info after 15m, got %s", got)
	}

	tr.clock = func() time.Time { return now.Add(2 * time.Hour) }
	if got := tr.Evaluate("tcp:9000"); got != LevelSilenced {
		t.Fatalf("expected silenced after 2h, got %s", got)
	}
}

func TestResetRestartsCritical(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr := New(testSteps)
	tr.clock = func() time.Time { return now }
	tr.Evaluate("tcp:443")

	tr.clock = func() time.Time { return now.Add(30 * time.Minute) }
	tr.Reset("tcp:443")

	if got := tr.Evaluate("tcp:443"); got != LevelCritical {
		t.Fatalf("expected critical after reset, got %s", got)
	}
}

func TestFlushClearsAll(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr := New(testSteps)
	tr.clock = func() time.Time { return now }
	tr.Evaluate("tcp:80")
	tr.Evaluate("tcp:443")
	tr.Flush()

	if got := tr.Evaluate("tcp:80"); got != LevelCritical {
		t.Fatalf("expected critical after flush, got %s", got)
	}
}

func TestLevelStringValues(t *testing.T) {
	cases := map[Level]string{
		LevelCritical: "critical",
		LevelWarning:  "warning",
		LevelInfo:     "info",
		LevelSilenced: "silenced",
	}
	for l, want := range cases {
		if got := l.String(); got != want {
			t.Errorf("Level(%d).String() = %q, want %q", l, got, want)
		}
	}
}

func TestIndependentKeysDoNotInterfere(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr := New(testSteps)
	tr.clock = func() time.Time { return now }
	tr.Evaluate("tcp:22")

	tr.clock = func() time.Time { return now.Add(5 * time.Minute) }
	tr.Evaluate("tcp:22") // warning

	// new key should still be critical
	if got := tr.Evaluate("tcp:3306"); got != LevelCritical {
		t.Fatalf("new key should be critical, got %s", got)
	}
}
