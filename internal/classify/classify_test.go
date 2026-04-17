package classify_test

import (
	"testing"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port uint16, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, State: "open"}
}

func TestClassifyOpenedWarningByDefault(t *testing.T) {
	c := classify.New(nil)
	r := c.ClassifyOpened(makePort(8080, "tcp"))
	if r.Severity != classify.SeverityWarning {
		t.Fatalf("expected warning, got %s", r.Severity)
	}
	if r.Category != classify.CategoryNewPort {
		t.Fatalf("expected new_port, got %s", r.Category)
	}
}

func TestClassifyOpenedCriticalForKnownPort(t *testing.T) {
	c := classify.New([]uint16{22, 3306})
	r := c.ClassifyOpened(makePort(22, "tcp"))
	if r.Severity != classify.SeverityCritical {
		t.Fatalf("expected critical, got %s", r.Severity)
	}
}

func TestClassifyOpenedNonCriticalPort(t *testing.T) {
	c := classify.New([]uint16{22})
	r := c.ClassifyOpened(makePort(9090, "tcp"))
	if r.Severity != classify.SeverityWarning {
		t.Fatalf("expected warning, got %s", r.Severity)
	}
}

func TestClassifyClosedIsInfo(t *testing.T) {
	c := classify.New([]uint16{22})
	r := c.ClassifyClosed(makePort(22, "tcp"))
	if r.Severity != classify.SeverityInfo {
		t.Fatalf("expected info, got %s", r.Severity)
	}
	if r.Category != classify.CategoryClosedPort {
		t.Fatalf("expected closed_port, got %s", r.Category)
	}
}

func TestClassifyLabelFields(t *testing.T) {
	c := classify.New(nil)
	op := c.ClassifyOpened(makePort(80, "tcp"))
	cl := c.ClassifyClosed(makePort(80, "tcp"))
	if op.Label != "opened" {
		t.Fatalf("expected 'opened', got %s", op.Label)
	}
	if cl.Label != "closed" {
		t.Fatalf("expected 'closed', got %s", cl.Label)
	}
}

func TestClassifyPortIsPreserved(t *testing.T) {
	c := classify.New(nil)
	p := makePort(443, "tcp")
	r := c.ClassifyOpened(p)
	if r.Port != p {
		t.Fatalf("port not preserved in result")
	}
}
