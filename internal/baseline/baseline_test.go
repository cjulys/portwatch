package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func makePorts(nums ...int) []scanner.Port {
	ports := make([]scanner.Port, len(nums))
	for i, n := range nums {
		ports[i] = scanner.Port{Number: n, Protocol: "tcp", State: "open"}
	}
	return ports
}

func TestSetAndGet(t *testing.T) {
	b := baseline.New(tempPath(t))
	ports := makePorts(80, 443)
	if err := b.Set(ports); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got := b.Get()
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(got))
	}
}

func TestLoadPersisted(t *testing.T) {
	path := tempPath(t)
	b1 := baseline.New(path)
	_ = b1.Set(makePorts(22, 8080))

	b2 := baseline.New(path)
	if err := b2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(b2.Get()) != 2 {
		t.Fatalf("expected 2 ports after reload, got %d", len(b2.Get()))
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	b := baseline.New(tempPath(t))
	if err := b.Load(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !b.IsEmpty() {
		t.Fatal("expected empty baseline")
	}
}

func TestClearRemovesFile(t *testing.T) {
	path := tempPath(t)
	b := baseline.New(path)
	_ = b.Set(makePorts(80))

	if err := b.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if !b.IsEmpty() {
		t.Fatal("expected empty after clear")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected file to be removed")
	}
}

func TestClearIdempotent(t *testing.T) {
	b := baseline.New(tempPath(t))
	if err := b.Clear(); err != nil {
		t.Fatalf("Clear on missing file: %v", err)
	}
}
