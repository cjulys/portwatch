package normalize_test

import (
	"testing"

	"github.com/user/portwatch/internal/normalize"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port int, proto, addr string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Address: addr}
}

func TestApplyNoOptions(t *testing.T) {
	n := normalize.New()
	input := []scanner.Port{makePort(80, "TCP", " 0.0.0.0 ")}
	out := n.Apply(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 port, got %d", len(out))
	}
	if out[0].Protocol != "TCP" {
		t.Errorf("expected protocol unchanged, got %q", out[0].Protocol)
	}
	if out[0].Address != " 0.0.0.0 " {
		t.Errorf("expected address unchanged, got %q", out[0].Address)
	}
}

func TestWithLowerProtocol(t *testing.T) {
	n := normalize.New(normalize.WithLowerProtocol())
	out := n.Apply([]scanner.Port{makePort(443, "TCP", "0.0.0.0")})
	if out[0].Protocol != "tcp" {
		t.Errorf("expected 'tcp', got %q", out[0].Protocol)
	}
}

func TestWithTrimAddress(t *testing.T) {
	n := normalize.New(normalize.WithTrimAddress())
	out := n.Apply([]scanner.Port{makePort(22, "tcp", "  127.0.0.1  ")})
	if out[0].Address != "127.0.0.1" {
		t.Errorf("expected trimmed address, got %q", out[0].Address)
	}
}

func TestWithDeduplication(t *testing.T) {
	n := normalize.New(normalize.WithLowerProtocol(), normalize.WithDeduplication())
	input := []scanner.Port{
		makePort(80, "tcp", "0.0.0.0"),
		makePort(80, "tcp", "0.0.0.0"),
		makePort(443, "tcp", "0.0.0.0"),
	}
	out := n.Apply(input)
	if len(out) != 2 {
		t.Errorf("expected 2 unique ports, got %d", len(out))
	}
}

func TestDeduplicationProtocolDistinct(t *testing.T) {
	n := normalize.New(normalize.WithDeduplication())
	input := []scanner.Port{
		makePort(53, "tcp", "0.0.0.0"),
		makePort(53, "udp", "0.0.0.0"),
	}
	out := n.Apply(input)
	if len(out) != 2 {
		t.Errorf("tcp/53 and udp/53 should be distinct, got %d ports", len(out))
	}
}

func TestApplyEmptySlice(t *testing.T) {
	n := normalize.New(normalize.WithLowerProtocol(), normalize.WithDeduplication())
	out := n.Apply(nil)
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d", len(out))
	}
}

func TestOriginalSliceUnmodified(t *testing.T) {
	n := normalize.New(normalize.WithLowerProtocol())
	orig := []scanner.Port{makePort(8080, "TCP", "0.0.0.0")}
	_ = n.Apply(orig)
	if orig[0].Protocol != "TCP" {
		t.Error("Apply must not modify the original slice")
	}
}
