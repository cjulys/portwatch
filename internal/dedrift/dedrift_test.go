package dedrift_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/dedrift"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, num int) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestEvaluateNoDrift(t *testing.T) {
	d := dedrift.New(nil)
	base := []scanner.Port{makePort("tcp", 80), makePort("tcp", 443)}
	curr := []scanner.Port{makePort("tcp", 80), makePort("tcp", 443)}

	r := d.Evaluate(base, curr)

	if r.HasDrift() {
		t.Fatalf("expected no drift, got %d entries", len(r.Entries))
	}
}

func TestEvaluateAddedPort(t *testing.T) {
	d := dedrift.New(nil)
	base := []scanner.Port{makePort("tcp", 80)}
	curr := []scanner.Port{makePort("tcp", 80), makePort("tcp", 8080)}

	r := d.Evaluate(base, curr)

	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Kind != dedrift.KindAdded {
		t.Errorf("expected KindAdded, got %s", r.Entries[0].Kind)
	}
	if r.Entries[0].Port.Number != 8080 {
		t.Errorf("expected port 8080, got %d", r.Entries[0].Port.Number)
	}
}

func TestEvaluateRemovedPort(t *testing.T) {
	d := dedrift.New(nil)
	base := []scanner.Port{makePort("tcp", 80), makePort("tcp", 22)}
	curr := []scanner.Port{makePort("tcp", 80)}

	r := d.Evaluate(base, curr)

	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Kind != dedrift.KindRemoved {
		t.Errorf("expected KindRemoved, got %s", r.Entries[0].Kind)
	}
}

func TestEvaluateEmptyBaseline(t *testing.T) {
	d := dedrift.New(nil)
	curr := []scanner.Port{makePort("tcp", 443)}

	r := d.Evaluate(nil, curr)

	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Kind != dedrift.KindAdded {
		t.Errorf("expected KindAdded, got %s", r.Entries[0].Kind)
	}
}

func TestSummariseNoDrift(t *testing.T) {
	d := dedrift.New(nil)
	r := dedrift.Report{}
	var buf bytes.Buffer
	d.Summarise(r, &buf)
	if !strings.Contains(buf.String(), "no drift") {
		t.Errorf("expected 'no drift' message, got %q", buf.String())
	}
}

func TestSummariseListsDriftedPorts(t *testing.T) {
	d := dedrift.New(nil)
	base := []scanner.Port{makePort("udp", 53)}
	curr := []scanner.Port{makePort("tcp", 80)}
	r := d.Evaluate(base, curr)

	var buf bytes.Buffer
	d.Summarise(r, &buf)
	out := buf.String()

	if !strings.Contains(out, "added") {
		t.Errorf("expected 'added' in output, got %q", out)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' in output, got %q", out)
	}
}
