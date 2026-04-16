package fingerprint_test

import (
	"testing"

	"portwatch/internal/fingerprint"
	"portwatch/internal/scanner"
)

func makePort(proto string, num int, state string) scanner.Port {
	return scanner.Port{Protocol: proto, Number: num, State: state}
}

func TestPortFingerprintIsStable(t *testing.T) {
	p := makePort("tcp", 80, "open")
	if fingerprint.Port(p) != fingerprint.Port(p) {
		t.Fatal("fingerprint not stable across calls")
	}
}

func TestPortFingerprintDiffersOnChange(t *testing.T) {
	a := fingerprint.Port(makePort("tcp", 80, "open"))
	b := fingerprint.Port(makePort("tcp", 443, "open"))
	if a == b {
		t.Fatal("expected different fingerprints for different ports")
	}
}

func TestPortsOrderIndependent(t *testing.T) {
	ports1 := []scanner.Port{makePort("tcp", 80, "open"), makePort("tcp", 443, "open")}
	ports2 := []scanner.Port{makePort("tcp", 443, "open"), makePort("tcp", 80, "open")}
	if fingerprint.Ports(ports1) != fingerprint.Ports(ports2) {
		t.Fatal("fingerprint should be order-independent")
	}
}

func TestPortsEmptyReturnsEmpty(t *testing.T) {
	if fingerprint.Ports(nil) != "" {
		t.Fatal("expected empty string for nil slice")
	}
	if fingerprint.Ports([]scanner.Port{}) != "" {
		t.Fatal("expected empty string for empty slice")
	}
}

func TestChangedDetectsDifference(t *testing.T) {
	prev := []scanner.Port{makePort("tcp", 22, "open")}
	curr := []scanner.Port{makePort("tcp", 22, "open"), makePort("tcp", 8080, "open")}
	if !fingerprint.Changed(prev, curr) {
		t.Fatal("expected Changed to return true")
	}
}

func TestChangedReturnsFalseWhenSame(t *testing.T) {
	ports := []scanner.Port{makePort("tcp", 22, "open")}
	if fingerprint.Changed(ports, ports) {
		t.Fatal("expected Changed to return false for identical slices")
	}
}
