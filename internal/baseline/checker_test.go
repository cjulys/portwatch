package baseline_test

import (
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func port(n int) scanner.Port {
	return scanner.Port{Number: n, Protocol: "tcp", State: "open"}
}

func TestCheckNoViolations(t *testing.T) {
	base := []scanner.Port{port(80), port(443)}
	current := []scanner.Port{port(80), port(443)}

	v := baseline.Check(base, current)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheckUnexpectedOpen(t *testing.T) {
	base := []scanner.Port{port(80)}
	current := []scanner.Port{port(80), port(9999)}

	v := baseline.Check(base, current)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Reason != "unexpected_open" {
		t.Errorf("expected unexpected_open, got %s", v[0].Reason)
	}
	if v[0].Port.Number != 9999 {
		t.Errorf("expected port 9999, got %d", v[0].Port.Number)
	}
}

func TestCheckExpectedClosed(t *testing.T) {
	base := []scanner.Port{port(80), port(22)}
	current := []scanner.Port{port(80)}

	v := baseline.Check(base, current)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Reason != "expected_closed" {
		t.Errorf("expected expected_closed, got %s", v[0].Reason)
	}
}

func TestCheckEmptyBaseAllViolations(t *testing.T) {
	var base []scanner.Port
	current := []scanner.Port{port(80), port(443)}

	v := baseline.Check(base, current)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestViolationCount(t *testing.T) {
	base := []scanner.Port{port(22)}
	current := []scanner.Port{port(80), port(443)}

	v := baseline.Check(base, current)
	unexpected, closed := baseline.ViolationCount(v)
	if unexpected != 2 {
		t.Errorf("expected 2 unexpected, got %d", unexpected)
	}
	if closed != 1 {
		t.Errorf("expected 1 closed, got %d", closed)
	}
}
