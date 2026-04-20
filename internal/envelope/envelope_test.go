package envelope_test

import (
	"testing"
	"time"

	"portwatch/internal/envelope"
	"portwatch/internal/scanner"
)

func makePort(num int) scanner.Port {
	return scanner.Port{Number: num, Protocol: "tcp", State: "open"}
}

func TestNewSetsCreatedAt(t *testing.T) {
	before := time.Now().UTC()
	e := envelope.New("e1", []scanner.Port{makePort(80)}, envelope.PriorityNormal)
	after := time.Now().UTC()
	if e.CreatedAt.Before(before) || e.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v outside expected range", e.CreatedAt)
	}
}

func TestNewStoresPorts(t *testing.T) {
	ports := []scanner.Port{makePort(22), makePort(443)}
	e := envelope.New("e2", ports, envelope.PriorityHigh)
	if len(e.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(e.Ports))
	}
}

func TestPriorityString(t *testing.T) {
	cases := []struct {
		p    envelope.Priority
		want string
	}{
		{envelope.PriorityLow, "low"},
		{envelope.PriorityNormal, "normal"},
		{envelope.PriorityHigh, "high"},
	}
	for _, c := range cases {
		if got := c.p.String(); got != c.want {
			t.Errorf("Priority.String() = %q, want %q", got, c.want)
		}
	}
}

func TestWithTagAndHasTag(t *testing.T) {
	e := envelope.New("e3", nil, envelope.PriorityLow)
	e.WithTag("critical").WithTag("external")
	if !e.HasTag("critical") {
		t.Error("expected tag 'critical'")
	}
	if e.HasTag("missing") {
		t.Error("unexpected tag 'missing'")
	}
}

func TestWithMeta(t *testing.T) {
	e := envelope.New("e4", nil, envelope.PriorityNormal)
	e.WithMeta("host", "localhost").WithMeta("zone", "us-east")
	if e.Meta["host"] != "localhost" {
		t.Errorf("meta host = %q, want 'localhost'", e.Meta["host"])
	}
	if e.Meta["zone"] != "us-east" {
		t.Errorf("meta zone = %q, want 'us-east'", e.Meta["zone"])
	}
}

func TestMetaInitialisedOnNew(t *testing.T) {
	e := envelope.New("e5", nil, envelope.PriorityLow)
	if e.Meta == nil {
		t.Error("Meta map should be initialised")
	}
}
