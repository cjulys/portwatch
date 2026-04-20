package pipeline_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

// startTCPListener binds a TCP listener on a random port and returns it.
func startTCPListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startTCPListener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

// TestPipelineDetectsRealOpenPort verifies that the pipeline's Run method
// returns the port that is actually listening on the loopback interface.
func TestPipelineDetectsRealOpenPort(t *testing.T) {
	ln, port := startTCPListener(t)
	defer ln.Close()

	sc := scanner.New([]string{"127.0.0.1"}, []int{port})
	p := pipeline.New(sc, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if result.OpenCount == 0 {
		t.Errorf("expected at least one open port, got 0")
	}

	found := false
	for _, pt := range result.Ports {
		if pt.Port == port {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("port %d not found in scan results: %+v", port, result.Ports)
	}
}

// TestPipelineClosedPortNotReported verifies that after the listener is closed
// the port no longer appears in a subsequent scan.
func TestPipelineClosedPortNotReported(t *testing.T) {
	ln, port := startTCPListener(t)
	// Close before scanning.
	ln.Close()

	// Give the OS a moment to release the port.
	time.Sleep(20 * time.Millisecond)

	sc := scanner.New([]string{"127.0.0.1"}, []int{port})
	p := pipeline.New(sc, nil, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	for _, pt := range result.Ports {
		if pt.Port == port {
			t.Errorf("closed port %d should not appear in results", port)
		}
	}
}

// TestPipelineScanTimeIsRecent checks that the ScanTime recorded by Run is
// within a reasonable window of the current time.
func TestPipelineScanTimeIsRecent(t *testing.T) {
	sc := scanner.New([]string{"127.0.0.1"}, []int{})
	p := pipeline.New(sc, nil, nil)

	before := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := p.Run(ctx)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	after := time.Now()

	if result.ScanTime.Before(before) || result.ScanTime.After(after) {
		t.Errorf("ScanTime %v not between %v and %v", result.ScanTime, before, after)
	}
}
