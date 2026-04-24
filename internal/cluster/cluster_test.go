package cluster_test

import (
	"testing"

	"github.com/user/portwatch/internal/cluster"
	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, port uint16) scanner.Port {
	return scanner.Port{Protocol: proto, Port: port, State: "open"}
}

func TestGroupEmptyReturnsNil(t *testing.T) {
	if got := cluster.Group(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestGroupSinglePort(t *testing.T) {
	ports := []scanner.Port{makePort("tcp", 80)}
	ranges := cluster.Group(ports)
	if len(ranges) != 1 {
		t.Fatalf("expected 1 range, got %d", len(ranges))
	}
	if ranges[0].String() != "tcp/80" {
		t.Errorf("unexpected string: %s", ranges[0].String())
	}
	if ranges[0].Size() != 1 {
		t.Errorf("expected size 1, got %d", ranges[0].Size())
	}
}

func TestGroupContiguousRange(t *testing.T) {
	ports := []scanner.Port{
		makePort("tcp", 8080),
		makePort("tcp", 8081),
		makePort("tcp", 8082),
	}
	ranges := cluster.Group(ports)
	if len(ranges) != 1 {
		t.Fatalf("expected 1 range, got %d", len(ranges))
	}
	if ranges[0].String() != "tcp/8080-8082" {
		t.Errorf("unexpected string: %s", ranges[0].String())
	}
	if ranges[0].Size() != 3 {
		t.Errorf("expected size 3, got %d", ranges[0].Size())
	}
}

func TestGroupNonContiguousSplit(t *testing.T) {
	ports := []scanner.Port{
		makePort("tcp", 80),
		makePort("tcp", 443),
		makePort("tcp", 444),
	}
	ranges := cluster.Group(ports)
	if len(ranges) != 2 {
		t.Fatalf("expected 2 ranges, got %d", len(ranges))
	}
}

func TestGroupProtocolSeparation(t *testing.T) {
	ports := []scanner.Port{
		makePort("tcp", 53),
		makePort("udp", 53),
		makePort("udp", 54),
	}
	ranges := cluster.Group(ports)
	if len(ranges) != 2 {
		t.Fatalf("expected 2 ranges (one per protocol), got %d", len(ranges))
	}
}

func TestGroupUnsortedInput(t *testing.T) {
	ports := []scanner.Port{
		makePort("tcp", 8082),
		makePort("tcp", 8080),
		makePort("tcp", 8081),
	}
	ranges := cluster.Group(ports)
	if len(ranges) != 1 {
		t.Fatalf("expected 1 range for unsorted contiguous input, got %d", len(ranges))
	}
	if ranges[0].String() != "tcp/8080-8082" {
		t.Errorf("unexpected string: %s", ranges[0].String())
	}
}
