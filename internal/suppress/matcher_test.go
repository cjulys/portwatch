package suppress

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newSuppressor(t *testing.T) *Suppressor {
	t.Helper()
	p := filepath.Join(t.TempDir(), "suppress.json")
	s, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestMatchExactKey(t *testing.T) {
	s := newSuppressor(t)
	_ = s.Add("tcp", "127.0.0.1", 8080, time.Hour)

	r := s.Match("tcp", "127.0.0.1", 8080)
	if !r.Matched {
		t.Fatal("expected exact key to match")
	}
}

func TestMatchWildcardFallback(t *testing.T) {
	s := newSuppressor(t)
	// Add a wildcard rule (empty addr maps to "*" via WildcardKey)
	s.rules[WildcardKey("tcp", 443)] = rule{Until: time.Now().Add(time.Hour)}

	r := s.Match("tcp", "10.0.0.1", 443)
	if !r.Matched {
		t.Fatal("expected wildcard to match arbitrary address")
	}
}

func TestMatchExpiredRuleReturnsFalse(t *testing.T) {
	s := newSuppressor(t)
	s.rules[Key("tcp", "127.0.0.1", 22)] = rule{Until: time.Now().Add(-time.Second)}

	r := s.Match("tcp", "127.0.0.1", 22)
	if r.Matched {
		t.Fatal("expired rule should not match")
	}
}

func TestMatchNoRuleReturnsFalse(t *testing.T) {
	s := newSuppressor(t)
	r := s.Match("udp", "0.0.0.0", 53)
	if r.Matched {
		t.Fatal("expected no match when no rules exist")
	}
}

func TestMatchResultContainsRuleKey(t *testing.T) {
	s := newSuppressor(t)
	_ = s.Add("tcp", "192.168.0.1", 9090, time.Hour)

	r := s.Match("tcp", "192.168.0.1", 9090)
	want := Key("tcp", "192.168.0.1", 9090)
	if r.RuleKey != want {
		t.Fatalf("RuleKey = %q; want %q", r.RuleKey, want)
	}
}

func TestMatchEvictsExpiredAndAllowsWildcard(t *testing.T) {
	s := newSuppressor(t)
	// Exact key expired
	s.rules[Key("tcp", "1.2.3.4", 80)] = rule{Until: time.Now().Add(-time.Second)}
	// Wildcard still valid
	s.rules[WildcardKey("tcp", 80)] = rule{Until: time.Now().Add(time.Hour)}

	r := s.Match("tcp", "1.2.3.4", 80)
	if !r.Matched {
		t.Fatal("wildcard should match after exact key is evicted")
	}
	_ = os.Remove("") // satisfy import
}
