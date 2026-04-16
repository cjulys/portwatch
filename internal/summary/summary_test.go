package summary_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/scanner"
	"portwatch/internal/summary"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return t.TempDir() + "/history.json"
}

func makePorts(nums ...int) []scanner.Port {
	out := make([]scanner.Port, len(nums))
	for i, n := range nums {
		out[i] = scanner.Port{Port: n, Protocol: "tcp", State: "open"}
	}
	return out
}

func TestBuildEmptyHistory(t *testing.T) {
	h := history.New(tempPath(t), 100)
	b := summary.New(h, nil)
	r := b.Build(time.Hour)
	if len(r.NewPorts) != 0 || len(r.ClosedPorts) != 0 {
		t.Fatalf("expected empty report, got %+v", r)
	}
}

func TestBuildCountsNewAndClosed(t *testing.T) {
	h := history.New(tempPath(t), 100)
	h.Add(history.Entry{
		Timestamp: time.Now().UTC(),
		Opened:    makePorts(80, 443),
		Closed:    makePorts(8080),
	})

	b := summary.New(h, nil)
	r := b.Build(time.Hour)

	if len(r.NewPorts) != 2 {
		t.Errorf("expected 2 new ports, got %d", len(r.NewPorts))
	}
	if len(r.ClosedPorts) != 1 {
		t.Errorf("expected 1 closed port, got %d", len(r.ClosedPorts))
	}
}

func TestBuildWindowFiltersOldEntries(t *testing.T) {
	h := history.New(tempPath(t), 100)
	h.Add(history.Entry{
		Timestamp: time.Now().UTC().Add(-2 * time.Hour),
		Opened:    makePorts(22),
	})

	b := summary.New(h, nil)
	r := b.Build(time.Hour) // window is only 1 hour

	if len(r.NewPorts) != 0 {
		t.Errorf("expected old entry to be filtered, got %d new ports", len(r.NewPorts))
	}
}

func TestPrintContainsSummaryFields(t *testing.T) {
	h := history.New(tempPath(t), 100)
	var buf bytes.Buffer
	b := summary.New(h, &buf)
	r := b.Build(time.Hour)
	b.Print(r)

	out := buf.String()
	for _, want := range []string{"Summary", "New ports", "Ports closed", "Distinct open"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}
