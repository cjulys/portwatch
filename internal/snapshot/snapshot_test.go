package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "snap.json")
}

func makePorts() []scanner.Port {
	return []scanner.Port{
		{Port: 80, Proto: "tcp", State: "open"},
		{Port: 443, Proto: "tcp", State: "open"},
	}
}

func TestTakeSetsTimestamp(t *testing.T) {
	before := time.Now().UTC()
	s := snapshot.Take(makePorts())
	after := time.Now().UTC()
	if s.CapturedAt.Before(before) || s.CapturedAt.After(after) {
		t.Errorf("unexpected timestamp: %v", s.CapturedAt)
	}
}

func TestTakeStoresPorts(t *testing.T) {
	ports := makePorts()
	s := snapshot.Take(ports)
	if len(s.Ports) != len(ports) {
		t.Fatalf("expected %d ports, got %d", len(ports), len(s.Ports))
	}
}

func TestSaveAndLoad(t *testing.T) {
	path := tempPath(t)
	orig := snapshot.Take(makePorts())
	if err := snapshot.SaveToFile(orig, path); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := snapshot.LoadFromFile(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("port count mismatch: got %d want %d", len(loaded.Ports), len(orig.Ports))
	}
	if !loaded.CapturedAt.Equal(orig.CapturedAt) {
		t.Errorf("timestamp mismatch")
	}
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	s, err := snapshot.LoadFromFile("/nonexistent/snap.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Ports) != 0 {
		t.Errorf("expected empty ports")
	}
}

func TestSaveCreatesFile(t *testing.T) {
	path := tempPath(t)
	if err := snapshot.SaveToFile(snapshot.Take(makePorts()), path); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not created: %v", err)
	}
}
