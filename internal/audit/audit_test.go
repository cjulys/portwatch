package audit

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "audit.log")
}

func TestRecordWritesJSON(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)
	if err := l.Record("scan", "ports checked"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	var e Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Event != "scan" {
		t.Errorf("event = %q, want scan", e.Event)
	}
	if e.Detail != "ports checked" {
		t.Errorf("detail = %q, want 'ports checked'", e.Detail)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestRecordTimestampIsUTC(t *testing.T) {
	var buf bytes.Buffer
	l := NewWithWriter(&buf)
	_ = l.Record("x", "y")
	var e Entry
	_ = json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &e)
	if e.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC, got %v", e.Timestamp.Location())
	}
}

func TestReadAllPersisted(t *testing.T) {
	p := tempPath(t)
	l, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = l.Record("open", "80/tcp")
	_ = l.Record("close", "443/tcp")

	entries, err := ReadAll(p)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Event != "open" {
		t.Errorf("entries[0].Event = %q", entries[0].Event)
	}
	if entries[1].Detail != "443/tcp" {
		t.Errorf("entries[1].Detail = %q", entries[1].Detail)
	}
}

func TestReadAllMissingFileReturnsNil(t *testing.T) {
	entries, err := ReadAll("/nonexistent/path/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil, got %v", entries)
	}
}

func TestNewCreatesFile(t *testing.T) {
	p := tempPath(t)
	l, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	_ = l.Record("init", "started")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		t.Error("expected file to exist")
	}
}
