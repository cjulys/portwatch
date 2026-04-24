// Package geoip provides lightweight country/ASN tagging for IP addresses
// observed on open ports. It uses a simple in-memory map seeded from a
// CSV file so the daemon has zero external runtime dependencies.
package geoip

import (
	"encoding/csv"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Record holds the enrichment data returned for a single IP lookup.
type Record struct {
	IP      string
	Country string // ISO 3166-1 alpha-2, e.g. "US"
	ASN     string // e.g. "AS15169"
	Org     string // human-readable organisation name
}

// DB is a compiled lookup table loaded from a CSV source.
type DB struct {
	entries []entry
}

type entry struct {
	start   net.IP
	end     net.IP
	country string
	asn     string
	org     string
}

// New returns an empty DB. Call LoadCSV to populate it.
func New() *DB { return &DB{} }

// LoadCSV reads a CSV file with columns: start_ip,end_ip,country,asn,org.
// Existing entries are replaced.
func (db *DB) LoadCSV(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("geoip: open %s: %w", path, err)
	}
	defer f.Close()
	return db.load(f)
}

func (db *DB) load(r io.Reader) error {
	cr := csv.NewReader(r)
	cr.Comment = '#'
	cr.TrimLeadingSpace = true
	records, err := cr.ReadAll()
	if err != nil {
		return fmt.Errorf("geoip: parse csv: %w", err)
	}
	entries := make([]entry, 0, len(records))
	for _, row := range records {
		if len(row) < 5 {
			continue
		}
		e := entry{
			start:   net.ParseIP(strings.TrimSpace(row[0])),
			end:     net.ParseIP(strings.TrimSpace(row[1])),
			country: strings.TrimSpace(row[2]),
			asn:     strings.TrimSpace(row[3]),
			org:     strings.TrimSpace(row[4]),
		}
		if e.start == nil || e.end == nil {
			continue
		}
		entries = append(entries, e)
	}
	db.entries = entries
	return nil
}

// Lookup returns the Record for ip. If no range matches, Country is "XX".
func (db *DB) Lookup(ip string) Record {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return Record{IP: ip, Country: "XX"}
	}
	for _, e := range db.entries {
		if bytesGTE(parsed, e.start) && bytesLTE(parsed, e.end) {
			return Record{IP: ip, Country: e.country, ASN: e.asn, Org: e.org}
		}
	}
	return Record{IP: ip, Country: "XX"}
}

func bytesGTE(a, b net.IP) bool {
	a4, b4 := a.To4(), b.To4()
	if a4 != nil && b4 != nil {
		a, b = a4, b4
	}
	for i := range a {
		if i >= len(b) {
			break
		}
		if a[i] > b[i] {
			return true
		}
		if a[i] < b[i] {
			return false
		}
	}
	return true
}

func bytesLTE(a, b net.IP) bool {
	a4, b4 := a.To4(), b.To4()
	if a4 != nil && b4 != nil {
		a, b = a4, b4
	}
	for i := range a {
		if i >= len(b) {
			break
		}
		if a[i] < b[i] {
			return true
		}
		if a[i] > b[i] {
			return false
		}
	}
	return true
}
