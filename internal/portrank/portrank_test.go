package portrank

import (
	"testing"
)

func TestGetBuiltinSSHIsCritical(t *testing.T) {
	r := New(nil)
	if got := r.Get("tcp", 22); got != RankCritical {
		t.Fatalf("expected critical, got %s", got)
	}
}

func TestGetBuiltinHTTPIsHigh(t *testing.T) {
	r := New(nil)
	if got := r.Get("tcp", 80); got != RankHigh {
		t.Fatalf("expected high, got %s", got)
	}
}

func TestGetUnknownPortIsUnknown(t *testing.T) {
	r := New(nil)
	if got := r.Get("tcp", 9999); got != RankUnknown {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	ovr := map[string]Rank{"tcp:9999": RankCritical}
	r := New(ovr)
	if got := r.Get("tcp", 9999); got != RankCritical {
		t.Fatalf("expected critical from override, got %s", got)
	}
}

func TestOverrideCanDowngradeBuiltin(t *testing.T) {
	ovr := map[string]Rank{"tcp:22": RankLow}
	r := New(ovr)
	if got := r.Get("tcp", 22); got != RankLow {
		t.Fatalf("expected low from override, got %s", got)
	}
}

func TestIsCriticalTrueForSSH(t *testing.T) {
	r := New(nil)
	if !r.IsCritical("tcp", 22) {
		t.Fatal("expected IsCritical to be true for SSH")
	}
}

func TestIsCriticalFalseForHTTP(t *testing.T) {
	r := New(nil)
	if r.IsCritical("tcp", 80) {
		t.Fatal("expected IsCritical to be false for HTTP")
	}
}

func TestRankStringValues(t *testing.T) {
	cases := []struct {
		rank Rank
		want string
	}{
		{RankUnknown, "unknown"},
		{RankLow, "low"},
		{RankMedium, "medium"},
		{RankHigh, "high"},
		{RankCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.rank.String(); got != tc.want {
			t.Errorf("Rank(%d).String() = %q, want %q", tc.rank, got, tc.want)
		}
	}
}

func TestProtocolDistinction(t *testing.T) {
	r := New(nil)
	// UDP port 53 is high, TCP port 53 is not in builtins
	if got := r.Get("udp", 53); got != RankHigh {
		t.Fatalf("expected high for udp:53, got %s", got)
	}
	if got := r.Get("tcp", 53); got != RankUnknown {
		t.Fatalf("expected unknown for tcp:53, got %s", got)
	}
}
