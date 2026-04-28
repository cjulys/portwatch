package watchlist

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func port(p uint16, proto string) scanner.Port {
	return scanner.Port{Port: p, Protocol: proto, State: "open"}
}

func TestCheckNoViolationsWhenAllPresent(t *testing.T) {
	wl := New()
	wl.Add(22, "tcp")
	wl.Add(80, "tcp")
	ports := []scanner.Port{port(22, "tcp"), port(80, "tcp")}
	v := Check(wl, ports)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheckMissingPortProducesViolation(t *testing.T) {
	wl := New()
	wl.Add(22, "tcp")
	v := Check(wl, nil)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Entry.Port != 22 {
		t.Fatalf("expected port 22, got %d", v[0].Entry.Port)
	}
}

func TestCheckEmptyWatchlistNeverViolates(t *testing.T) {
	wl := New()
	ports := []scanner.Port{port(22, "tcp"), port(80, "tcp")}
	v := Check(wl, ports)
	if len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}

func TestCheckProtocolDistinction(t *testing.T) {
	wl := New()
	wl.Add(53, "tcp")
	// Only UDP/53 is open, not TCP/53.
	ports := []scanner.Port{port(53, "udp")}
	v := Check(wl, ports)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
}

func TestViolationCountMatchesCheckLen(t *testing.T) {
	wl := New()
	wl.Add(22, "tcp")
	wl.Add(443, "tcp")
	ports := []scanner.Port{port(22, "tcp")}
	if ViolationCount(wl, ports) != 1 {
		t.Fatal("expected ViolationCount to return 1")
	}
}
