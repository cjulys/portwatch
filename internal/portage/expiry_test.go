package portage

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeEntry(port int, proto string, lastSeen time.Time) Entry {
	return Entry{
		Port:      scanner.Port{Port: port, Protocol: proto},
		FirstSeen: lastSeen.Add(-1 * time.Hour),
		LastSeen:  lastSeen,
	}
}

func TestDefaultExpiryPolicyMaxAge(t *testing.T) {
	p := DefaultExpiryPolicy()
	if p.MaxAge != 72*time.Hour {
		t.Fatalf("expected 72h, got %v", p.MaxAge)
	}
}

func TestExpiredReturnsTrueWhenOld(t *testing.T) {
	now := time.Now()
	policy := ExpiryPolicy{MaxAge: 1 * time.Hour}
	entry := makeEntry(80, "tcp", now.Add(-2*time.Hour))
	if !policy.Expired(entry, now) {
		t.Fatal("expected entry to be expired")
	}
}

func TestExpiredReturnsFalseWhenFresh(t *testing.T) {
	now := time.Now()
	policy := ExpiryPolicy{MaxAge: 1 * time.Hour}
	entry := makeEntry(80, "tcp", now.Add(-30*time.Minute))
	if policy.Expired(entry, now) {
		t.Fatal("expected entry to be fresh")
	}
}

func TestExpiredZeroMaxAgeNeverExpires(t *testing.T) {
	now := time.Now()
	policy := ExpiryPolicy{MaxAge: 0}
	entry := makeEntry(80, "tcp", now.Add(-1000*time.Hour))
	if policy.Expired(entry, now) {
		t.Fatal("zero MaxAge should never expire")
	}
}

func TestPruneExpiredSeparatesEntries(t *testing.T) {
	now := time.Now()
	policy := ExpiryPolicy{MaxAge: 1 * time.Hour}

	entries := []Entry{
		makeEntry(22, "tcp", now.Add(-2*time.Hour)),  // expired
		makeEntry(80, "tcp", now.Add(-30*time.Minute)), // fresh
		makeEntry(443, "tcp", now.Add(-90*time.Minute)), // expired
	}

	kept, evicted := PruneExpired(entries, policy, now)

	if len(kept) != 1 {
		t.Fatalf("expected 1 kept, got %d", len(kept))
	}
	if kept[0].Port.Port != 80 {
		t.Fatalf("expected port 80 to be kept, got %d", kept[0].Port.Port)
	}
	if len(evicted) != 2 {
		t.Fatalf("expected 2 evicted, got %d", len(evicted))
	}
}

func TestPruneExpiredEmptyInput(t *testing.T) {
	now := time.Now()
	policy := DefaultExpiryPolicy()
	kept, evicted := PruneExpired(nil, policy, now)
	if kept != nil || evicted != nil {
		t.Fatal("expected nil slices for empty input")
	}
}
