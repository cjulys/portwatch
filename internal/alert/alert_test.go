package alert

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, port int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: port}
}

func TestNotifyNewPort(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	diffs := []scanner.Diff{
		{Port: makePort("tcp", 8080), State: scanner.StateNew},
	}
	n.Notify(diffs)

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level, got: %s", out)
	}
	if !strings.Contains(out, "new port opened") {
		t.Errorf("expected 'new port opened' message, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port number in output, got: %s", out)
	}
}

func TestNotifyClosedPort(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	diffs := []scanner.Diff{
		{Port: makePort("tcp", 443), State: scanner.StateClosed},
	}
	n.Notify(diffs)

	out := buf.String()
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level, got: %s", out)
	}
	if !strings.Contains(out, "port closed") {
		t.Errorf("expected 'port closed' message, got: %s", out)
	}
}

func TestNotifyNoDiffs(t *testing.T) {
	var buf bytes.Buffer
	n := New(&buf)

	n.Notify([]scanner.Diff{})

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diffs, got: %s", buf.String())
	}
}

func TestNotifyDefaultWriter(t *testing.T) {
	// Ensure New(nil) does not panic
	n := New(nil)
	if n.out == nil {
		t.Error("expected non-nil writer when nil is passed to New")
	}
}
