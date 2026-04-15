package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(port uint16, proto string) scanner.PortState {
	return scanner.PortState{Port: port, Protocol: proto, Open: true}
}

func TestEmptyFilterPassesAll(t *testing.T) {
	f := filter.New(nil)
	ports := []scanner.PortState{
		makePort(80, "tcp"),
		makePort(443, "tcp"),
		makePort(53, "udp"),
	}
	got := f.Apply(ports)
	if len(got) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(got))
	}
}

func TestFilterAllowsOnlyMatchingPorts(t *testing.T) {
	rules := []filter.Rule{
		{Port: 80, Protocol: "tcp"},
		{Port: 443, Protocol: "tcp"},
	}
	f := filter.New(rules)
	ports := []scanner.PortState{
		makePort(80, "tcp"),
		makePort(443, "tcp"),
		makePort(8080, "tcp"),
		makePort(53, "udp"),
	}
	got := f.Apply(ports)
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
	for _, p := range got {
		if p.Port != 80 && p.Port != 443 {
			t.Errorf("unexpected port %d in result", p.Port)
		}
	}
}

func TestFilterProtocolDistinction(t *testing.T) {
	rules := []filter.Rule{
		{Port: 53, Protocol: "tcp"},
	}
	f := filter.New(rules)
	ports := []scanner.PortState{
		makePort(53, "tcp"),
		makePort(53, "udp"), // same port, different protocol — should be excluded
	}
	got := f.Apply(ports)
	if len(got) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got))
	}
	if got[0].Protocol != "tcp" {
		t.Errorf("expected tcp, got %s", got[0].Protocol)
	}
}

func TestContainsEmptyFilter(t *testing.T) {
	f := filter.New(nil)
	if !f.Contains(9999, "tcp") {
		t.Error("empty filter should contain every port")
	}
}

func TestContainsWithRules(t *testing.T) {
	rules := []filter.Rule{{Port: 22, Protocol: "tcp"}}
	f := filter.New(rules)
	if !f.Contains(22, "tcp") {
		t.Error("expected port 22/tcp to be contained")
	}
	if f.Contains(22, "udp") {
		t.Error("expected port 22/udp to NOT be contained")
	}
	if f.Contains(80, "tcp") {
		t.Error("expected port 80/tcp to NOT be contained")
	}
}

func TestApplyEmptyPortList(t *testing.T) {
	rules := []filter.Rule{{Port: 80, Protocol: "tcp"}}
	f := filter.New(rules)
	got := f.Apply([]scanner.PortState{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d ports", len(got))
	}
}
