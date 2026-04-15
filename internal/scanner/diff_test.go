package scanner_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/scanner"
)

func makePort(proto, addr string, port int) scanner.PortState {
	return scanner.PortState{Protocol: proto, Address: addr, Port: port}
}

func TestCompareNoDiff(t *testing.T) {
	prev := []scanner.PortState{makePort("TCP", "127.0.0.1", 80)}
	curr := []scanner.PortState{makePort("TCP", "127.0.0.1", 80)}

	d := scanner.Compare(prev, curr)
	if d.HasChanges() {
		t.Error("expected no changes")
	}
}

func TestCompareNewPortOpened(t *testing.T) {
	prev := []scanner.PortState{makePort("TCP", "127.0.0.1", 80)}
	curr := []scanner.PortState{
		makePort("TCP", "127.0.0.1", 80),
		makePort("TCP", "127.0.0.1", 443),
	}

	d := scanner.Compare(prev, curr)
	if !d.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(d.Opened) != 1 {
		t.Fatalf("expected 1 opened port, got %d", len(d.Opened))
	}
	if d.Opened[0].Port != 443 {
		t.Errorf("expected opened port 443, got %d", d.Opened[0].Port)
	}
	if len(d.Closed) != 0 {
		t.Errorf("expected 0 closed ports, got %d", len(d.Closed))
	}
}

func TestComparePortClosed(t *testing.T) {
	prev := []scanner.PortState{
		makePort("TCP", "127.0.0.1", 80),
		makePort("TCP", "127.0.0.1", 8080),
	}
	curr := []scanner.PortState{makePort("TCP", "127.0.0.1", 80)}

	d := scanner.Compare(prev, curr)
	if !d.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(d.Closed) != 1 {
		t.Fatalf("expected 1 closed port, got %d", len(d.Closed))
	}
	if d.Closed[0].Port != 8080 {
		t.Errorf("expected closed port 8080, got %d", d.Closed[0].Port)
	}
}

func TestCompareEmptySlices(t *testing.T) {
	d := scanner.Compare(nil, nil)
	if d.HasChanges() {
		t.Error("expected no changes for empty slices")
	}
}
