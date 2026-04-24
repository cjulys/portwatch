// Package geoip enriches scanner.Port values with geographic and network
// metadata (country, ASN, organisation) derived from an IP-range CSV database.
//
// Typical usage:
//
//	db := geoip.New()
//	if err := db.LoadCSV("/etc/portwatch/geoip.csv"); err != nil {
//		log.Println("geoip unavailable:", err)
//	}
//	enricher := geoip.NewEnricher(db)
//	enrichedPorts := enricher.EnrichAll(scannedPorts)
//
// The CSV format is:
//
//	start_ip,end_ip,country,asn,org
//
// Lines beginning with '#' are treated as comments and ignored.
// If the database is empty or an IP does not match any range the country
// field is set to "XX" (unknown).
package geoip
