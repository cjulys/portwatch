package digest_test

import (
	"testing"

	"portwatch/internal/digest"
	"portwatch/internal/scanner"
)

func ports(pairs ...interface{}) []scanner.Port {
	var out []scanner.Port
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, scanner.Port{
			Number:   pairs[i].(int),
			Protocol: pairs[i+1].(string),
		})
	}
	return out
}

func TestEqualSameContents(t *testing.T) {
	a := ports(80, "tcp", 443, "tcp")
	b := ports(443, "tcp", 80, "tcp")
	if !digest.Equal(a, b) {
		t.Error("expected Equal to return true for same ports in different order")
	}
}

func TestEqualDifferentContents(t *testing.T) {
	a := ports(80, "tcp")
	b := ports(8080, "tcp")
	if digest.Equal(a, b) {
		t.Error("expected Equal to return false for different ports")
	}
}

func TestEqualBothEmpty(t *testing.T) {
	if !digest.Equal(nil, nil) {
		t.Error("expected Equal to return true for two nil slices")
	}
}

func TestEqualOneEmpty(t *testing.T) {
	a := ports(22, "tcp")
	if digest.Equal(a, nil) {
		t.Error("expected Equal to return false when one slice is empty")
	}
}

func TestComputeStableAcrossCalls(t *testing.T) {
	p := ports(22, "tcp", 80, "tcp", 443, "tcp")
	d1 := digest.Compute(p)
	d2 := digest.Compute(p)
	if d1 != d2 {
		t.Errorf("expected stable digest, got %q and %q", d1, d2)
	}
}
