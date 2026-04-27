package census

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePorts(specs [][2]string) []scanner.Port {
	var out []scanner.Port
	for _, s := range specs {
		out = append(out, scanner.Port{Protocol: s[0], Address: s[1]})
	}
	return out
}

func TestRecordFirstScanDeltaEqualsSnapshot(t *testing.T) {
	c := New()
	ports := makePorts([][2]string{{"tcp", "0.0.0.0:80"}, {"tcp", "0.0.0.0:443"}, {"udp", "0.0.0.0:53"}})
	snap, delta := c.Record(ports)

	if snap.Total != 3 {
		t.Fatalf("want Total=3 got %d", snap.Total)
	}
	if snap.TCP != 2 {
		t.Fatalf("want TCP=2 got %d", snap.TCP)
	}
	if snap.UDP != 1 {
		t.Fatalf("want UDP=1 got %d", snap.UDP)
	}
	if delta.Total != snap.Total || delta.TCP != snap.TCP || delta.UDP != snap.UDP {
		t.Fatalf("first-scan delta should equal snapshot counts, got %+v", delta)
	}
}

func TestRecordDeltaReflectsChange(t *testing.T) {
	c := New()
	first := makePorts([][2]string{{"tcp", "0.0.0.0:80"}})
	c.Record(first)

	second := makePorts([][2]string{{"tcp", "0.0.0.0:80"}, {"tcp", "0.0.0.0:8080"}, {"udp", "0.0.0.0:53"}})
	snap, delta := c.Record(second)

	if snap.Total != 3 {
		t.Fatalf("want Total=3 got %d", snap.Total)
	}
	if delta.Total != 2 {
		t.Fatalf("want delta.Total=2 got %d", delta.Total)
	}
	if delta.TCP != 1 {
		t.Fatalf("want delta.TCP=1 got %d", delta.TCP)
	}
	if delta.UDP != 1 {
		t.Fatalf("want delta.UDP=1 got %d", delta.UDP)
	}
}

func TestRecordNegativeDelta(t *testing.T) {
	c := New()
	first := makePorts([][2]string{{"tcp", "0.0.0.0:80"}, {"tcp", "0.0.0.0:443"}})
	c.Record(first)

	_, delta := c.Record([]scanner.Port{})
	if delta.Total != -2 {
		t.Fatalf("want delta.Total=-2 got %d", delta.Total)
	}
}

func TestLastReturnsNilBeforeFirstScan(t *testing.T) {
	c := New()
	if c.Last() != nil {
		t.Fatal("expected nil before any scan")
	}
}

func TestLastReturnsSnapshot(t *testing.T) {
	c := New()
	ports := makePorts([][2]string{{"tcp", "0.0.0.0:22"}})
	before := time.Now().UTC()
	c.Record(ports)
	after := time.Now().UTC()

	snap := c.Last()
	if snap == nil {
		t.Fatal("expected non-nil snapshot")
	}
	if snap.Total != 1 {
		t.Fatalf("want Total=1 got %d", snap.Total)
	}
	if snap.At.Before(before) || snap.At.After(after) {
		t.Fatalf("timestamp %v not in expected range [%v, %v]", snap.At, before, after)
	}
}

func TestResetClearsState(t *testing.T) {
	c := New()
	c.Record(makePorts([][2]string{{"tcp", "0.0.0.0:80"}}))
	c.Reset()

	if c.Last() != nil {
		t.Fatal("expected nil after reset")
	}

	// After reset the next Record should behave like the first scan.
	_, delta := c.Record(makePorts([][2]string{{"tcp", "0.0.0.0:80"}}))
	if delta.Total != 1 {
		t.Fatalf("want delta.Total=1 after reset, got %d", delta.Total)
	}
}
