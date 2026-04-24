package remap_test

import (
	"testing"

	"github.com/user/portwatch/internal/remap"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Proto: proto, Port: port, State: "open"}
}

func TestApplyKnownPortGetsAlias(t *testing.T) {
	m := remap.New([]remap.Rule{
		{Proto: "tcp", Port: 80, Alias: "http"},
	})
	ports := m.Apply([]scanner.Port{makePort("tcp", 80)})
	if ports[0].Alias != "http" {
		t.Fatalf("expected alias 'http', got %q", ports[0].Alias)
	}
}

func TestApplyUnknownPortLeavesAliasEmpty(t *testing.T) {
	m := remap.New([]remap.Rule{
		{Proto: "tcp", Port: 443, Alias: "https"},
	})
	ports := m.Apply([]scanner.Port{makePort("tcp", 22)})
	if ports[0].Alias != "" {
		t.Fatalf("expected empty alias, got %q", ports[0].Alias)
	}
}

func TestApplyProtocolDistinction(t *testing.T) {
	m := remap.New([]remap.Rule{
		{Proto: "tcp", Port: 53, Alias: "dns-tcp"},
		{Proto: "udp", Port: 53, Alias: "dns-udp"},
	})
	ports := m.Apply([]scanner.Port{
		makePort("tcp", 53),
		makePort("udp", 53),
	})
	if ports[0].Alias != "dns-tcp" {
		t.Errorf("tcp alias: want 'dns-tcp', got %q", ports[0].Alias)
	}
	if ports[1].Alias != "dns-udp" {
		t.Errorf("udp alias: want 'dns-udp', got %q", ports[1].Alias)
	}
}

func TestApplyEmptyInputReturnsEmpty(t *testing.T) {
	m := remap.New([]remap.Rule{{Proto: "tcp", Port: 80, Alias: "http"}})
	ports := m.Apply(nil)
	if len(ports) != 0 {
		t.Fatalf("expected empty slice, got %d elements", len(ports))
	}
}

func TestAliasDirectLookup(t *testing.T) {
	m := remap.New([]remap.Rule{{Proto: "tcp", Port: 22, Alias: "ssh"}})
	if got := m.Alias("tcp", 22); got != "ssh" {
		t.Fatalf("want 'ssh', got %q", got)
	}
	if got := m.Alias("tcp", 80); got != "" {
		t.Fatalf("want empty, got %q", got)
	}
}

func TestDuplicateRuleLastWins(t *testing.T) {
	m := remap.New([]remap.Rule{
		{Proto: "tcp", Port: 8080, Alias: "alt-http"},
		{Proto: "tcp", Port: 8080, Alias: "proxy"},
	})
	if got := m.Alias("tcp", 8080); got != "proxy" {
		t.Fatalf("want 'proxy', got %q", got)
	}
}

func TestLenReflectsUniqueRules(t *testing.T) {
	m := remap.New([]remap.Rule{
		{Proto: "tcp", Port: 80, Alias: "http"},
		{Proto: "tcp", Port: 443, Alias: "https"},
	})
	if m.Len() != 2 {
		t.Fatalf("expected 2 rules, got %d", m.Len())
	}
}
