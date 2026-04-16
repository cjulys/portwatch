// Package audit records scan events to an append-only log for later review.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Entry is a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Detail    string    `json:"detail"`
}

// Logger appends audit entries to a file.
type Logger struct {
	mu   sync.Mutex
	out  io.Writer
	path string
}

// New opens (or creates) the audit log at path.
func New(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return &Logger{out: f, path: path}, nil
}

// NewWithWriter creates a Logger that writes to w (useful for tests).
func NewWithWriter(w io.Writer) *Logger {
	return &Logger{out: w}
}

// Record writes an audit entry.
func (l *Logger) Record(event, detail string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Detail:    detail,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err = fmt.Fprintf(l.out, "%s\n", b)
	return err
}

// ReadAll parses all entries from path.
func ReadAll(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read %s: %w", path, err)
	}
	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err == nil {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

func splitLines(b []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, c := range b {
		if c == '\n' {
			lines = append(lines, b[start:i])
			start = i + 1
		}
	}
	if start < len(b) {
		lines = append(lines, b[start:])
	}
	return lines
}
