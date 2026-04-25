package drift

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func port(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Port: num}
}

func TestNoDriftWhenIdentical(t *testing.T) {
	base := []scanner.Port{port("tcp", 80), port("tcp", 443)}
	tr := New(base)
	s := tr.Evaluate(base)
	if s.Total != 0 {
		t.Fatalf("expected 0 drift, got %d", s.Total)
	}
}

func TestAddedPortIncreasesScore(t *testing.T) {
	base := []scanner.Port{port("tcp", 80)}
	current := []scanner.Port{port("tcp", 80), port("tcp", 8080)}
	tr := New(base)
	s := tr.Evaluate(current)
	if s.Added != 1 {
		t.Fatalf("expected 1 added, got %d", s.Added)
	}
	if s.Removed != 0 {
		t.Fatalf("expected 0 removed, got %d", s.Removed)
	}
	if s.Total != 1 {
		t.Fatalf("expected total 1, got %d", s.Total)
	}
}

func TestRemovedPortIncreasesScore(t *testing.T) {
	base := []scanner.Port{port("tcp", 80), port("tcp", 443)}
	current := []scanner.Port{port("tcp", 80)}
	tr := New(base)
	s := tr.Evaluate(current)
	if s.Removed != 1 {
		t.Fatalf("expected 1 removed, got %d", s.Removed)
	}
	if s.Total != 1 {
		t.Fatalf("expected total 1, got %d", s.Total)
	}
}

func TestMixedChanges(t *testing.T) {
	base := []scanner.Port{port("tcp", 80), port("tcp", 443)}
	current := []scanner.Port{port("tcp", 80), port("tcp", 9090)}
	tr := New(base)
	s := tr.Evaluate(current)
	if s.Added != 1 || s.Removed != 1 || s.Total != 2 {
		t.Fatalf("unexpected score %+v", s)
	}
}

func TestSetBaselineReplacesOld(t *testing.T) {
	tr := New([]scanner.Port{port("tcp", 80)})
	tr.SetBaseline([]scanner.Port{port("tcp", 443)})
	s := tr.Evaluate([]scanner.Port{port("tcp", 443)})
	if s.Total != 0 {
		t.Fatalf("expected 0 drift after baseline update, got %d", s.Total)
	}
}

func TestLastReturnsMostRecentScore(t *testing.T) {
	tr := New([]scanner.Port{port("tcp", 80)})
	_ = tr.Evaluate([]scanner.Port{port("tcp", 8080)})
	l := tr.Last()
	if l.Total != 1 {
		t.Fatalf("Last() should return most recent score, got %+v", l)
	}
}

func TestEmptyBaselineAllAdded(t *testing.T) {
	tr := New(nil)
	s := tr.Evaluate([]scanner.Port{port("tcp", 80), port("udp", 53)})
	if s.Added != 2 {
		t.Fatalf("expected 2 added, got %d", s.Added)
	}
}

func TestScoreTimestampIsSet(t *testing.T) {
	tr := New(nil)
	s := tr.Evaluate(nil)
	if s.At.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}
