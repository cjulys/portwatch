package history

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "history.json")
}

func makePorts(numbers ...int) []scanner.Port {
	ports := make([]scanner.Port, len(numbers))
	for i, n := range numbers {
		ports[i] = scanner.Port{Number: n, Protocol: "tcp", State: "open"}
	}
	return ports
}

func TestAddAndRetrieve(t *testing.T) {
	h, err := New(tempPath(t), 10)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := h.Add(makePorts(80, 443)); err != nil {
		t.Fatalf("Add: %v", err)
	}
	entries := h.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if len(entries[0].Ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(entries[0].Ports))
	}
}

func TestMaxSizeEviction(t *testing.T) {
	h, err := New(tempPath(t), 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 5; i++ {
		if err := h.Add(makePorts(i + 1)); err != nil {
			t.Fatalf("Add: %v", err)
		}
	}
	if got := len(h.Entries()); got != 3 {
		t.Errorf("expected 3 entries after eviction, got %d", got)
	}
}

func TestPersistenceAcrossReload(t *testing.T) {
	path := tempPath(t)
	h, _ := New(path, 10)
	_ = h.Add(makePorts(22, 8080))

	h2, err := New(path, 10)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := len(h2.Entries()); got != 1 {
		t.Errorf("expected 1 entry after reload, got %d", got)
	}
}

func TestClearRemovesFile(t *testing.T) {
	path := tempPath(t)
	h, _ := New(path, 10)
	_ = h.Add(makePorts(80))
	if err := h.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Clear")
	}
	if got := len(h.Entries()); got != 0 {
		t.Errorf("expected 0 entries after Clear, got %d", got)
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	h, err := New(tempPath(t), 10)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := len(h.Entries()); got != 0 {
		t.Errorf("expected 0 entries for missing file, got %d", got)
	}
}
