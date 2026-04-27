package shadow_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/shadow"
)

func makePort(addr string, port int, proto string) scanner.Port {
	return scanner.Port{Address: addr, Port: port, Protocol: proto, State: "open"}
}

func TestCurrentNilBeforeCommit(t *testing.T) {
	s := shadow.New()
	if s.Current() != nil {
		t.Fatal("expected nil before any commit")
	}
}

func TestPreviousNilAfterOneCommit(t *testing.T) {
	s := shadow.New()
	s.Commit([]scanner.Port{makePort("127.0.0.1", 80, "tcp")})
	if s.Previous() != nil {
		t.Fatal("expected nil previous after only one commit")
	}
}

func TestCurrentReflectsLastCommit(t *testing.T) {
	s := shadow.New()
	ports := []scanner.Port{makePort("127.0.0.1", 443, "tcp")}
	s.Commit(ports)
	entry := s.Current()
	if entry == nil {
		t.Fatal("expected non-nil current entry")
	}
	if len(entry.Ports) != 1 || entry.Ports[0].Port != 443 {
		t.Fatalf("unexpected ports in current entry: %+v", entry.Ports)
	}
}

func TestPreviousAfterTwoCommits(t *testing.T) {
	s := shadow.New()
	first := []scanner.Port{makePort("127.0.0.1", 22, "tcp")}
	second := []scanner.Port{makePort("127.0.0.1", 80, "tcp")}
	s.Commit(first)
	s.Commit(second)
	prev := s.Previous()
	if prev == nil {
		t.Fatal("expected non-nil previous entry after two commits")
	}
	if len(prev.Ports) != 1 || prev.Ports[0].Port != 22 {
		t.Fatalf("previous entry should hold first scan: %+v", prev.Ports)
	}
}

func TestFlappingNilBeforeTwoCommits(t *testing.T) {
	s := shadow.New()
	if got := s.Flapping(); got != nil {
		t.Fatalf("expected nil flapping before two commits, got %v", got)
	}
	s.Commit([]scanner.Port{makePort("127.0.0.1", 8080, "tcp")})
	if got := s.Flapping(); got != nil {
		t.Fatalf("expected nil flapping after one commit, got %v", got)
	}
}

func TestFlappingDetectsDisappearedPort(t *testing.T) {
	s := shadow.New()
	s.Commit([]scanner.Port{makePort("127.0.0.1", 9090, "tcp")})
	s.Commit([]scanner.Port{}) // port gone
	flapping := s.Flapping()
	if len(flapping) != 1 {
		t.Fatalf("expected 1 flapping port, got %d", len(flapping))
	}
	if flapping[0].Port != 9090 {
		t.Fatalf("expected port 9090, got %d", flapping[0].Port)
	}
}

func TestFlappingDetectsNewPort(t *testing.T) {
	s := shadow.New()
	s.Commit([]scanner.Port{})
	s.Commit([]scanner.Port{makePort("127.0.0.1", 3306, "tcp")})
	flapping := s.Flapping()
	if len(flapping) != 1 {
		t.Fatalf("expected 1 flapping port, got %d", len(flapping))
	}
}

func TestFlappingEmptyWhenStable(t *testing.T) {
	s := shadow.New()
	ports := []scanner.Port{makePort("127.0.0.1", 443, "tcp")}
	s.Commit(ports)
	s.Commit(ports)
	if got := s.Flapping(); len(got) != 0 {
		t.Fatalf("expected no flapping for stable ports, got %v", got)
	}
}

func TestCommitRecordedAtIsSet(t *testing.T) {
	s := shadow.New()
	s.Commit(nil)
	entry := s.Current()
	if entry.RecordedAt.IsZero() {
		t.Fatal("expected RecordedAt to be set")
	}
}
