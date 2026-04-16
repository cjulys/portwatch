package suppress_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"portwatch/internal/suppress"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "suppress.json")
}

func TestAddAndIsSuppressed(t *testing.T) {
	s := suppress.New(tempPath(t))
	s.Add(suppress.Rule{Port: 8080, Protocol: "tcp"})
	if !s.IsSuppressed(8080, "tcp") {
		t.Fatal("expected port to be suppressed")
	}
}

func TestNotSuppressedWithoutRule(t *testing.T) {
	s := suppress.New(tempPath(t))
	if s.IsSuppressed(9090, "tcp") {
		t.Fatal("expected port not to be suppressed")
	}
}

func TestExpiredRuleNotSuppressed(t *testing.T) {
	s := suppress.New(tempPath(t))
	s.Add(suppress.Rule{Port: 443, Protocol: "tcp", Until: time.Now().Add(-time.Second)})
	if s.IsSuppressed(443, "tcp") {
		t.Fatal("expired rule should not suppress")
	}
}

func TestFutureRuleIsSuppressed(t *testing.T) {
	s := suppress.New(tempPath(t))
	s.Add(suppress.Rule{Port: 443, Protocol: "tcp", Until: time.Now().Add(time.Hour)})
	if !s.IsSuppressed(443, "tcp") {
		t.Fatal("future rule should suppress")
	}
}

func TestRemoveRule(t *testing.T) {
	s := suppress.New(tempPath(t))
	s.Add(suppress.Rule{Port: 22, Protocol: "tcp"})
	s.Remove(22, "tcp")
	if s.IsSuppressed(22, "tcp") {
		t.Fatal("removed rule should not suppress")
	}
}

func TestPersistenceAcrossReload(t *testing.T) {
	p := tempPath(t)
	s1 := suppress.New(p)
	s1.Add(suppress.Rule{Port: 80, Protocol: "tcp"})

	s2 := suppress.New(p)
	if !s2.IsSuppressed(80, "tcp") {
		t.Fatal("rule should persist across reload")
	}
}

func TestMissingFileReturnsEmptyStore(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	s := suppress.New(p)
	if s.IsSuppressed(1234, "udp") {
		t.Fatal("empty store should suppress nothing")
	}
	_, err := os.Stat(p)
	if err == nil {
		t.Fatal("file should not be created on load failure")
	}
}
