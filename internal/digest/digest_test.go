package digest_test

import (
	"testing"

	"portwatch/internal/digest"
	"portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestEmptySliceReturnsEmpty(t *testing.T) {
	if got := digest.Compute(nil); got != digest.Empty {
		t.Fatalf("expected Empty, got %q", got)
	}
	if got := digest.Compute([]scanner.Port{}); got != digest.Empty {
		t.Fatalf("expected Empty for zero-length slice, got %q", got)
	}
}

func TestDeterministic(t *testing.T) {
	ports := []scanner.Port{makePort("tcp", 80), makePort("tcp", 443)}
	a := digest.Compute(ports)
	b := digest.Compute(ports)
	if a != b {
		t.Fatalf("same input produced different digests: %q vs %q", a, b)
	}
}

func TestOrderIndependent(t *testing.T) {
	a := digest.Compute([]scanner.Port{makePort("tcp", 80), makePort("tcp", 443)})
	b := digest.Compute([]scanner.Port{makePort("tcp", 443), makePort("tcp", 80)})
	if a != b {
		t.Fatalf("order should not affect digest: %q vs %q", a, b)
	}
}

func TestDifferentPortsProduceDifferentDigests(t *testing.T) {
	a := digest.Compute([]scanner.Port{makePort("tcp", 80)})
	b := digest.Compute([]scanner.Port{makePort("tcp", 8080)})
	if digest.Equal(a, b) {
		t.Fatal("different ports should not produce equal digests")
	}
}

func TestProtoDistinction(t *testing.T) {
	a := digest.Compute([]scanner.Port{makePort("tcp", 53)})
	b := digest.Compute([]scanner.Port{makePort("udp", 53)})
	if digest.Equal(a, b) {
		t.Fatal("tcp:53 and udp:53 should produce different digests")
	}
}

func TestEqualHelper(t *testing.T) {
	p := []scanner.Port{makePort("tcp", 22)}
	if !digest.Equal(digest.Compute(p), digest.Compute(p)) {
		t.Fatal("Equal should return true for identical digests")
	}
}
