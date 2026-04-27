package label_test

import (
	"testing"

	"github.com/user/portwatch/internal/label"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, Address: "127.0.0.1"}
}

func TestBuiltinTCPLabel(t *testing.T) {
	l := label.New(nil)
	got := l.Apply(makePort("tcp", 22))
	if got != "SSH" {
		t.Fatalf("expected SSH, got %q", got)
	}
}

func TestBuiltinUDPLabel(t *testing.T) {
	l := label.New(nil)
	got := l.Apply(makePort("udp", 53))
	if got != "DNS" {
		t.Fatalf("expected DNS, got %q", got)
	}
}

func TestUnknownPortReturnsEmpty(t *testing.T) {
	l := label.New(nil)
	got := l.Apply(makePort("tcp", 9999))
	if got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestOverrideRuleTakesPrecedence(t *testing.T) {
	rules := []label.Rule{
		{Port: 22, Protocol: "tcp", Label: "Bastion"},
	}
	l := label.New(rules)
	got := l.Apply(makePort("tcp", 22))
	if got != "Bastion" {
		t.Fatalf("expected Bastion, got %q", got)
	}
}

func TestOverrideAnyProtocol(t *testing.T) {
	rules := []label.Rule{
		{Port: 1234, Protocol: "", Label: "custom"},
	}
	l := label.New(rules)
	if got := l.Apply(makePort("tcp", 1234)); got != "custom" {
		t.Fatalf("tcp: expected custom, got %q", got)
	}
	if got := l.Apply(makePort("udp", 1234)); got != "custom" {
		t.Fatalf("udp: expected custom, got %q", got)
	}
}

func TestProtocolCaseInsensitive(t *testing.T) {
	l := label.New(nil)
	got := l.Apply(makePort("TCP", 80))
	if got != "HTTP" {
		t.Fatalf("expected HTTP, got %q", got)
	}
}

func TestApplyAllReturnsMap(t *testing.T) {
	l := label.New(nil)
	ports := []scanner.Port{
		makePort("tcp", 80),
		makePort("tcp", 443),
		makePort("tcp", 9999),
	}
	m := l.ApplyAll(ports)
	if m["tcp:80"] != "HTTP" {
		t.Fatalf("tcp:80 expected HTTP, got %q", m["tcp:80"])
	}
	if m["tcp:443"] != "HTTPS" {
		t.Fatalf("tcp:443 expected HTTPS, got %q", m["tcp:443"])
	}
	if _, ok := m["tcp:9999"]; ok {
		t.Fatal("tcp:9999 should not appear in map")
	}
}
