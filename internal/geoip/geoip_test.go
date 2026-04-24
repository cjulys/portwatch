package geoip

import (
	"strings"
	"testing"
)

const sampleCSV = `
# start_ip,end_ip,country,asn,org
1.0.0.0,1.0.0.255,AU,AS13335,Cloudflare
8.8.8.0,8.8.8.255,US,AS15169,Google LLC
10.0.0.0,10.255.255.255,ZZ,AS64512,Private
`

func loadSample(t *testing.T) *DB {
	t.Helper()
	db := New()
	if err := db.load(strings.NewReader(sampleCSV)); err != nil {
		t.Fatalf("load: %v", err)
	}
	return db
}

func TestLookupKnownIP(t *testing.T) {
	db := loadSample(t)
	r := db.Lookup("8.8.8.8")
	if r.Country != "US" {
		t.Errorf("country: got %q, want US", r.Country)
	}
	if r.ASN != "AS15169" {
		t.Errorf("asn: got %q, want AS15169", r.ASN)
	}
	if r.Org != "Google LLC" {
		t.Errorf("org: got %q, want Google LLC", r.Org)
	}
}

func TestLookupUnknownIPReturnsXX(t *testing.T) {
	db := loadSample(t)
	r := db.Lookup("5.5.5.5")
	if r.Country != "XX" {
		t.Errorf("expected XX for unknown IP, got %q", r.Country)
	}
}

func TestLookupInvalidIPReturnsXX(t *testing.T) {
	db := loadSample(t)
	r := db.Lookup("not-an-ip")
	if r.Country != "XX" {
		t.Errorf("expected XX for invalid IP, got %q", r.Country)
	}
}

func TestLookupPrivateRange(t *testing.T) {
	db := loadSample(t)
	r := db.Lookup("10.42.0.1")
	if r.Country != "ZZ" {
		t.Errorf("expected ZZ for private range, got %q", r.Country)
	}
}

func TestLoadCSVMissingFileReturnsError(t *testing.T) {
	db := New()
	if err := db.LoadCSV("/nonexistent/path.csv"); err == nil {
		t.Error("expected error for missing file")
	}
}

func TestEmptyDBReturnsXX(t *testing.T) {
	db := New()
	r := db.Lookup("1.2.3.4")
	if r.Country != "XX" {
		t.Errorf("empty db: expected XX, got %q", r.Country)
	}
}

func TestLookupBoundaryStart(t *testing.T) {
	db := loadSample(t)
	r := db.Lookup("1.0.0.0")
	if r.Country != "AU" {
		t.Errorf("boundary start: got %q, want AU", r.Country)
	}
}

func TestLookupBoundaryEnd(t *testing.T) {
	db := loadSample(t)
	r := db.Lookup("1.0.0.255")
	if r.Country != "AU" {
		t.Errorf("boundary end: got %q, want AU", r.Country)
	}
}
