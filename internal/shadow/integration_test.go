package shadow_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/shadow"
)

// TestThreeCycleFlap verifies that a port appearing and disappearing across
// three consecutive commits is correctly reported as flapping in each
// transition cycle.
func TestThreeCycleFlap(t *testing.T) {
	s := shadow.New()

	p := makePort("10.0.0.1", 8443, "tcp")

	// Cycle 1: port present
	s.Commit([]scanner.Port{p})
	// Cycle 2: port gone — should show as flapping
	s.Commit([]scanner.Port{})
	if got := s.Flapping(); len(got) != 1 {
		t.Fatalf("cycle 2: expected 1 flapping port, got %d", len(got))
	}

	// Cycle 3: port back — should again show as flapping
	s.Commit([]scanner.Port{p})
	if got := s.Flapping(); len(got) != 1 {
		t.Fatalf("cycle 3: expected 1 flapping port, got %d", len(got))
	}
}

// TestProtocolDistinctionInFlapping ensures tcp:80 and udp:80 are tracked
// independently so only the truly changed protocol is reported.
func TestProtocolDistinctionInFlapping(t *testing.T) {
	s := shadow.New()

	tcp80 := makePort("127.0.0.1", 80, "tcp")
	udp80 := makePort("127.0.0.1", 80, "udp")

	s.Commit([]scanner.Port{tcp80, udp80})
	// Remove only the UDP port.
	s.Commit([]scanner.Port{tcp80})

	flapping := s.Flapping()
	if len(flapping) != 1 {
		t.Fatalf("expected exactly 1 flapping port (udp:80), got %d: %+v", len(flapping), flapping)
	}
	if flapping[0].Protocol != "udp" {
		t.Fatalf("expected flapping protocol=udp, got %s", flapping[0].Protocol)
	}
}
