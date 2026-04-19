package probe_test

import (
	"net"
	"testing"
	"time"

	"portwatch/internal/probe"
)

func startTCP(t *testing.T) (port int, stop func()) {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return l.Addr().(*net.TCPAddr).Port, func() { l.Close() }
}

func TestProbeReachablePort(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	p := probe.New(time.Second)
	r := p.Probe("127.0.0.1", port, "tcp")
	if !r.Reachable {
		t.Fatalf("expected reachable, got err: %v", r.Err)
	}
	if r.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbeUnreachablePort(t *testing.T) {
	p := probe.New(200 * time.Millisecond)
	r := p.Probe("127.0.0.1", 1, "tcp")
	if r.Reachable {
		t.Fatal("expected unreachable")
	}
	if r.Err == nil {
		t.Fatal("expected error")
	}
}

func TestProbeUnsupportedProtocol(t *testing.T) {
	p := probe.New(time.Second)
	r := p.Probe("127.0.0.1", 80, "udp")
	if r.Err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

func TestProbeAllReturnsAllResults(t *testing.T) {
	port, stop := startTCP(t)
	defer stop()

	p := probe.New(time.Second)
	targets := []probe.Target{
		{Port: port, Protocol: "tcp"},
		{Port: 1, Protocol: "tcp"},
	}
	results := p.ProbeAll("127.0.0.1", targets)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Reachable {
		t.Error("first target should be reachable")
	}
	if results[1].Reachable {
		t.Error("second target should not be reachable")
	}
}

func TestNewDefaultTimeout(t *testing.T) {
	p := probe.New(0)
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}
