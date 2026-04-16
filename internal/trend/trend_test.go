package trend_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/scanner"
	"portwatch/internal/trend"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func makePorts(n int) []scanner.Port {
	ports := make([]scanner.Port, n)
	for i := range ports {
		ports[i] = scanner.Port{Port: 8000 + i, Protocol: "tcp", State: "open"}
	}
	return ports
}

func buildHistory(t *testing.T, snapshots [][]scanner.Port, base time.Time) *history.History {
	t.Helper()
	h := history.New(tempPath(t), 100)
	for i, ports := range snapshots {
		h.Add(base.Add(time.Duration(i)*time.Minute), ports)
	}
	return h
}

func TestStableWhenCountUnchanged(t *testing.T) {
	now := time.Now().UTC()
	h := buildHistory(t, [][]scanner.Port{makePorts(3), makePorts(3)}, now.Add(-4*time.Minute))
	a := trend.New(h, 10*time.Minute)
	r := a.Analyze(now)
	if r.Direction != trend.Stable {
		t.Fatalf("expected Stable, got %s", r.Direction)
	}
	if r.Delta != 0 {
		t.Fatalf("expected delta 0, got %d", r.Delta)
	}
}

func TestRisingWhenPortsIncrease(t *testing.T) {
	now := time.Now().UTC()
	h := buildHistory(t, [][]scanner.Port{makePorts(2), makePorts(5)}, now.Add(-4*time.Minute))
	a := trend.New(h, 10*time.Minute)
	r := a.Analyze(now)
	if r.Direction != trend.Rising {
		t.Fatalf("expected Rising, got %s", r.Direction)
	}
	if r.Delta != 3 {
		t.Fatalf("expected delta 3, got %d", r.Delta)
	}
}

func TestFallingWhenPortsDecrease(t *testing.T) {
	now := time.Now().UTC()
	h := buildHistory(t, [][]scanner.Port{makePorts(5), makePorts(1)}, now.Add(-4*time.Minute))
	a := trend.New(h, 10*time.Minute)
	r := a.Analyze(now)
	if r.Direction != trend.Falling {
		t.Fatalf("expected Falling, got %s", r.Direction)
	}
	if r.Delta != -4 {
		t.Fatalf("expected delta -4, got %d", r.Delta)
	}
}

func TestEmptyHistoryIsStable(t *testing.T) {
	h := history.New(filepath.Join(t.TempDir(), "h.json"), 100)
	a := trend.New(h, 5*time.Minute)
	r := a.Analyze(time.Now().UTC())
	if r.Direction != trend.Stable {
		t.Fatalf("expected Stable on empty history, got %s", r.Direction)
	}
	if r.Samples != 0 {
		t.Fatalf("expected 0 samples, got %d", r.Samples)
	}
}

func TestDefaultWindowAppliedOnZero(t *testing.T) {
	_ = os.Getenv // suppress unused import lint
	h := history.New(filepath.Join(t.TempDir(), "h.json"), 100)
	a := trend.New(h, 0)
	r := a.Analyze(time.Now().UTC())
	if r.Window != 5*time.Minute {
		t.Fatalf("expected default window 5m, got %s", r.Window)
	}
}
