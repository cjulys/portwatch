package geoip

import (
	"fmt"
	"net"

	"portwatch/internal/scanner"
)

// Enricher attaches GeoIP metadata to scanner.Port values as string tags
// using the format "country=XX", "asn=AS1234", "org=Some Org".
type Enricher struct {
	db *DB
}

// NewEnricher returns an Enricher backed by db.
func NewEnricher(db *DB) *Enricher {
	return &Enricher{db: db}
}

// Enrich returns a copy of p with GeoIP tags appended to p.Tags.
// If the port has no associated remote address the original port is returned
// unchanged. Localhost addresses are skipped (country "LO" is injected).
func (e *Enricher) Enrich(p scanner.Port) scanner.Port {
	addr := p.Address
	if addr == "" {
		return p
	}
	ip := net.ParseIP(addr)
	if ip == nil {
		return p
	}
	if ip.IsLoopback() {
		return withTags(p, "country=LO", "", "")
	}
	r := e.db.Lookup(addr)
	return withTags(p, fmt.Sprintf("country=%s", r.Country), fmt.Sprintf("asn=%s", r.ASN), fmt.Sprintf("org=%s", r.Org))
}

// EnrichAll applies Enrich to every port in the slice and returns the result.
func (e *Enricher) EnrichAll(ports []scanner.Port) []scanner.Port {
	out := make([]scanner.Port, len(ports))
	for i, p := range ports {
		out[i] = e.Enrich(p)
	}
	return out
}

func withTags(p scanner.Port, country, asn, org string) scanner.Port {
	tags := make([]string, len(p.Tags))
	copy(tags, p.Tags)
	if country != "" {
		tags = append(tags, country)
	}
	if asn != "" {
		tags = append(tags, asn)
	}
	if org != "" {
		tags = append(tags, org)
	}
	p.Tags = tags
	return p
}
