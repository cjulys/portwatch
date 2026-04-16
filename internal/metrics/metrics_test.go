package metrics

import (
	"testing"
	"time"
)

func TestInitialSnapshotIsZero(t *testing.T) {
	c := New()
	s := c.Snapshot()
	if s.ScansTotal != 0 || s.AlertsTotal != 0 || s.OpenPortCount != 0 {
		t.Fatalf("expected zero snapshot, got %+v", s)
	}
	if !s.LastScanAt.IsZero() || !s.LastAlertAt.IsZero() {
		t.Fatal("expected zero times")
	}
}

func TestRecordScanIncrements(t *testing.T) {
	c := New()
	before := time.Now()
	c.RecordScan(7)
	s := c.Snapshot()
	if s.ScansTotal != 1 {
		t.Fatalf("want 1 scan, got %d", s.ScansTotal)
	}
	if s.OpenPortCount != 7 {
		t.Fatalf("want 7 open ports, got %d", s.OpenPortCount)
	}
	if s.LastScanAt.Before(before) {
		t.Fatal("LastScanAt should be recent")
	}
}

func TestRecordAlertIncrements(t *testing.T) {
	c := New()
	before := time.Now()
	c.RecordAlert()
	c.RecordAlert()
	s := c.Snapshot()
	if s.AlertsTotal != 2 {
		t.Fatalf("want 2 alerts, got %d", s.AlertsTotal)
	}
	if s.LastAlertAt.Before(before) {
		t.Fatal("LastAlertAt should be recent")
	}
}

func TestOpenPortCountUpdated(t *testing.T) {
	c := New()
	c.RecordScan(3)
	c.RecordScan(10)
	if c.Snapshot().OpenPortCount != 10 {
		t.Fatal("expected latest open port count")
	}
}

func TestResetClearsAll(t *testing.T) {
	c := New()
	c.RecordScan(5)
	c.RecordAlert()
	c.Reset()
	s := c.Snapshot()
	if s.ScansTotal != 0 || s.AlertsTotal != 0 {
		t.Fatal("expected zeroed metrics after reset")
	}
}
