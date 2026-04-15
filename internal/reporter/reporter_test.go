package reporter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/reporter"
	"github.com/user/portwatch/internal/scanner"
)

func makePorts(specs [][2]string) []scanner.Port {
	ports := make([]scanner.Port, 0, len(specs))
	for _, s := range specs {
		ports = append(ports, scanner.Port{Proto: s[0], Address: s[1]})
	}
	return ports
}

func TestReportTextFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)
	ports := makePorts([][2]string{{"tcp", "127.0.0.1:8080"}, {"tcp", "127.0.0.1:9090"}})

	if err := r.Report(ports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "open ports (2)") {
		t.Errorf("expected port count in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
}

func TestReportJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatJSON)
	ports := makePorts([][2]string{{"udp", "0.0.0.0:53"}})

	if err := r.Report(ports); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry reporter.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}
	if entry.Count != 1 {
		t.Errorf("expected count 1, got %d", entry.Count)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestReportEmptyPorts(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf, reporter.FormatText)

	if err := r.Report([]scanner.Port{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "open ports (0)") {
		t.Errorf("expected zero count, got: %s", buf.String())
	}
}

func TestDefaultWriterIsStdout(t *testing.T) {
	// Just verify New doesn't panic with nil writer and empty format.
	r := reporter.New(nil, "")
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
