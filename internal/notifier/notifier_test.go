package notifier_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, num uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num}
}

func TestDispatchCallsHandler(t *testing.T) {
	var received []notifier.Event
	h := func(e notifier.Event) { received = append(received, e) }

	n := notifier.New(nil, h)
	n.Dispatch(notifier.LevelAlert, "port opened", makePort("tcp", 8080))

	if len(received) != 1 {
		t.Fatalf("expected 1 event, got %d", len(received))
	}
	if received[0].Level != notifier.LevelAlert {
		t.Errorf("expected ALERT, got %s", received[0].Level)
	}
	if received[0].Port.Number != 8080 {
		t.Errorf("expected port 8080, got %d", received[0].Port.Number)
	}
}

func TestDispatchMultipleHandlers(t *testing.T) {
	count := 0
	inc := func(e notifier.Event) { count++ }

	n := notifier.New(nil, inc, inc, inc)
	n.Dispatch(notifier.LevelInfo, "test", makePort("udp", 53))

	if count != 3 {
		t.Errorf("expected 3 handler calls, got %d", count)
	}
}

func TestDispatchFallbackWriter(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf)
	n.Dispatch(notifier.LevelWarn, "unexpected port", makePort("tcp", 22))

	output := buf.String()
	if !strings.Contains(output, "WARN") {
		t.Errorf("expected WARN in output, got: %s", output)
	}
	if !strings.Contains(output, "unexpected port") {
		t.Errorf("expected message in output, got: %s", output)
	}
}

func TestEventTimestampIsRecent(t *testing.T) {
	before := time.Now().UTC()
	var got notifier.Event
	n := notifier.New(nil, func(e notifier.Event) { got = e })
	n.Dispatch(notifier.LevelInfo, "check", makePort("tcp", 443))
	after := time.Now().UTC()

	if got.Timestamp.Before(before) || got.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", got.Timestamp, before, after)
	}
}

func TestAddHandlerAtRuntime(t *testing.T) {
	n := notifier.New(nil)
	called := false
	n.AddHandler(func(e notifier.Event) { called = true })
	n.Dispatch(notifier.LevelInfo, "runtime handler", makePort("tcp", 80))

	if !called {
		t.Error("expected runtime handler to be called")
	}
}
