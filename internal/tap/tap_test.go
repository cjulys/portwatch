package tap_test

import (
	"bytes"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tap"
)

func makeDiff(port uint16, dir string) scanner.Diff {
	return scanner.Diff{
		Port:      scanner.Port{Port: port, Protocol: "tcp"},
		Direction: dir,
	}
}

func TestSendDeliveredToSink(t *testing.T) {
	tr := tap.New(nil)
	var got []scanner.Diff
	tr.Register(func(diffs []scanner.Diff) { got = append(got, diffs...) })

	tr.Send([]scanner.Diff{makeDiff(80, "opened")})

	if len(got) != 1 || got[0].Port.Port != 80 {
		t.Fatalf("expected diff for port 80, got %v", got)
	}
}

func TestSendDeliveredToMultipleSinks(t *testing.T) {
	tr := tap.New(nil)
	var a, b []scanner.Diff
	tr.Register(func(d []scanner.Diff) { a = append(a, d...) })
	tr.Register(func(d []scanner.Diff) { b = append(b, d...) })

	tr.Send([]scanner.Diff{makeDiff(443, "opened")})

	if len(a) != 1 || len(b) != 1 {
		t.Fatalf("expected both sinks to receive diff, got a=%d b=%d", len(a), len(b))
	}
}

func TestSendEmptyIsNoop(t *testing.T) {
	tr := tap.New(nil)
	called := false
	tr.Register(func(_ []scanner.Diff) { called = true })

	tr.Send(nil)

	if called {
		t.Fatal("sink should not be called for empty diff slice")
	}
}

func TestSinkPanicIsRecovered(t *testing.T) {
	buf := &bytes.Buffer{}
	tr := tap.New(buf)
	tr.Register(func(_ []scanner.Diff) { panic("boom") })

	// Must not propagate the panic.
	tr.Send([]scanner.Diff{makeDiff(22, "opened")})

	if buf.Len() == 0 {
		t.Fatal("expected fallback writer to receive panic notice")
	}
}

func TestRegisterNilSinkIsIgnored(t *testing.T) {
	tr := tap.New(nil)
	tr.Register(nil)
	if tr.Len() != 0 {
		t.Fatalf("expected 0 sinks, got %d", tr.Len())
	}
}

func TestDefaultFallbackIsStderr(t *testing.T) {
	// Constructing with nil fallback should not panic.
	tr := tap.New(nil)
	if tr == nil {
		t.Fatal("expected non-nil tap")
	}
}
