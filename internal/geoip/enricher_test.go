package geoip

import (
	"strings"
	"testing"

	"portwatch/internal/scanner"
)

func makePort(addr string) scanner.Port {
	return scanner.Port{Address: addr, Port: 80, Protocol: "tcp"}
}

func loadEnricher(t *testing.T) *Enricher {
	t.Helper()
	db := New()
	if err := db.load(strings.NewReader(sampleCSV)); err != nil {
		t.Fatalf("load: %v", err)
	}
	return NewEnricher(db)
}

func hasTag(tags []string, prefix string) bool {
	for _, t := range tags {
		if strings.HasPrefix(t, prefix) {
			return true
		}
	}
	return false
}

func tagValue(tags []string, prefix string) string {
	for _, t := range tags {
		if strings.HasPrefix(t, prefix) {
			return strings.TrimPrefix(t, prefix)
		}
	}
	return ""
}

func TestEnrichKnownIP(t *testing.T) {
	e := loadEnricher(t)
	p := e.Enrich(makePort("8.8.8.8"))
	if v := tagValue(p.Tags, "country="); v != "US" {
		t.Errorf("country tag: got %q, want US", v)
	}
	if !hasTag(p.Tags, "asn=") {
		t.Error("expected asn tag")
	}
}

func TestEnrichLoopback(t *testing.T) {
	e := loadEnricher(t)
	p := e.Enrich(makePort("127.0.0.1"))
	if v := tagValue(p.Tags, "country="); v != "LO" {
		t.Errorf("loopback: got %q, want LO", v)
	}
}

func TestEnrichEmptyAddressUnchanged(t *testing.T) {
	e := loadEnricher(t)
	p := makePort("")
	out := e.Enrich(p)
	if len(out.Tags) != 0 {
		t.Errorf("expected no tags for empty address, got %v", out.Tags)
	}
}

func TestEnrichUnknownIPGetsXX(t *testing.T) {
	e := loadEnricher(t)
	p := e.Enrich(makePort("5.5.5.5"))
	if v := tagValue(p.Tags, "country="); v != "XX" {
		t.Errorf("unknown: got %q, want XX", v)
	}
}

func TestEnrichAllPreservesLength(t *testing.T) {
	e := loadEnricher(t)
	ports := []scanner.Port{makePort("8.8.8.8"), makePort("1.0.0.1"), makePort("")}
	out := e.EnrichAll(ports)
	if len(out) != len(ports) {
		t.Errorf("length: got %d, want %d", len(out), len(ports))
	}
}

func TestEnrichPreservesExistingTags(t *testing.T) {
	e := loadEnricher(t)
	p := makePort("8.8.8.8")
	p.Tags = []string{"existing=yes"}
	out := e.Enrich(p)
	if !hasTag(out.Tags, "existing=") {
		t.Error("existing tag was lost")
	}
	if !hasTag(out.Tags, "country=") {
		t.Error("country tag missing")
	}
}
