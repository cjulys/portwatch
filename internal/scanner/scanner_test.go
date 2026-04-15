package scanner_test

import (
	"net"
	"testing"

	"github.com/yourorg/portwatch/internal/scanner"
)

func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScanDetectsOpenPort(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	s := scanner.New([]string{"tcp"})
	states, err := s.Scan([2]int{port, port})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected 1 open port, got %d", len(states))
	}
	if states[0].Port != port {
		t.Errorf("expected port %d, got %d", port, states[0].Port)
	}
	if states[0].Protocol != "TCP" {
		t.Errorf("expected protocol TCP, got %s", states[0].Protocol)
	}
}

func TestScanNoOpenPorts(t *testing.T) {
	s := scanner.New([]string{"tcp"})
	// Use a port range unlikely to have anything listening in CI
	states, err := s.Scan([2]int{19999, 19999})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 0 {
		t.Errorf("expected 0 open ports, got %d", len(states))
	}
}

func TestPortStateString(t *testing.T) {
	ps := scanner.PortState{Protocol: "TCP", Port: 8080, Address: "127.0.0.1"}
	expected := "127.0.0.1:8080 (TCP)"
	if ps.String() != expected {
		t.Errorf("expected %q, got %q", expected, ps.String())
	}
}
