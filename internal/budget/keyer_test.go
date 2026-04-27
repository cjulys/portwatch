package budget

import (
	"strings"
	"testing"
)

func TestScanKeyFormat(t *testing.T) {
	k := ScanKey("localhost")
	if !strings.HasPrefix(k, "scan:") {
		t.Fatalf("ScanKey should start with 'scan:', got %q", k)
	}
	if !strings.Contains(k, "localhost") {
		t.Fatalf("ScanKey should contain target, got %q", k)
	}
}

func TestHostKeyFormat(t *testing.T) {
	k := HostKey("192.168.1.1")
	if !strings.HasPrefix(k, "host:") {
		t.Fatalf("HostKey should start with 'host:', got %q", k)
	}
	if !strings.Contains(k, "192.168.1.1") {
		t.Fatalf("HostKey should contain host, got %q", k)
	}
}

func TestScanKeyAndHostKeyAreDifferent(t *testing.T) {
	if ScanKey("x") == HostKey("x") {
		t.Fatal("ScanKey and HostKey with same argument should produce different keys")
	}
}

func TestScanKeysDifferByTarget(t *testing.T) {
	if ScanKey("a") == ScanKey("b") {
		t.Fatal("ScanKeys for different targets should differ")
	}
}
