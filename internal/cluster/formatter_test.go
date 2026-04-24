package cluster_test

import (
	"testing"

	"github.com/user/portwatch/internal/cluster"
	"github.com/user/portwatch/internal/scanner"
)

func TestSummaryEmpty(t *testing.T) {
	if got := cluster.Summary(nil); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestSummaryMultipleRanges(t *testing.T) {
	ports := []scanner.Port{
		{Protocol: "tcp", Port: 22, State: "open"},
		{Protocol: "tcp", Port: 80, State: "open"},
		{Protocol: "tcp", Port: 81, State: "open"},
		{Protocol: "udp", Port: 53, State: "open"},
	}
	got := cluster.Summary(ports)
	// Expect three ranges: tcp/22, tcp/80-81, udp/53
	for _, want := range []string{"tcp/22", "tcp/80-81", "udp/53"} {
		if !containsSubstr(got, want) {
			t.Errorf("summary %q missing expected segment %q", got, want)
		}
	}
}

func TestRangeCountReduction(t *testing.T) {
	ports := make([]scanner.Port, 10)
	for i := range ports {
		ports[i] = scanner.Port{Protocol: "tcp", Port: uint16(8000 + i), State: "open"}
	}
	if got := cluster.RangeCount(ports); got != 1 {
		t.Errorf("expected 1 range for 10 contiguous ports, got %d", got)
	}
}

func TestRangeCountNoReduction(t *testing.T) {
	ports := []scanner.Port{
		{Protocol: "tcp", Port: 80, State: "open"},
		{Protocol: "tcp", Port: 443, State: "open"},
		{Protocol: "tcp", Port: 8080, State: "open"},
	}
	if got := cluster.RangeCount(ports); got != 3 {
		t.Errorf("expected 3 ranges, got %d", got)
	}
}

func containsSubstr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 &&
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
