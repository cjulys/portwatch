package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "state.json")
}

func makePorts() []scanner.Port {
	return []scanner.Port{
		{Number: 80, Protocol: "tcp", State: "open"},
		{Number: 443, Protocol: "tcp", State: "open"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	store := state.New(tempPath(t))
	ports := makePorts()

	if err := store.Save(ports); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != len(ports) {
		t.Errorf("expected %d ports, got %d", len(ports), len(snap.Ports))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	store := state.New(tempPath(t))
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %d", len(snap.Ports))
	}
}

func TestClearRemovesFile(t *testing.T) {
	path := tempPath(t)
	store := state.New(path)

	if err := store.Save(makePorts()); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := store.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected file to be removed after Clear")
	}
}

func TestClearNonExistentFileIsNoop(t *testing.T) {
	store := state.New(tempPath(t))
	if err := store.Clear(); err != nil {
		t.Errorf("Clear on missing file should not error: %v", err)
	}
}

func TestSaveOverwritesPreviousState(t *testing.T) {
	store := state.New(tempPath(t))

	_ = store.Save(makePorts())
	newPorts := []scanner.Port{{Number: 22, Protocol: "tcp", State: "open"}}
	_ = store.Save(newPorts)

	snap, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(snap.Ports) != 1 {
		t.Errorf("expected 1 port after overwrite, got %d", len(snap.Ports))
	}
	if snap.Ports[0].Number != 22 {
		t.Errorf("expected port 22, got %d", snap.Ports[0].Number)
	}
}
