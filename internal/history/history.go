package history

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry records a snapshot of port state at a point in time.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Ports     []scanner.Port `json:"ports"`
}

// History manages a rolling log of port scan entries.
type History struct {
	path    string
	maxSize int
	entries []Entry
}

// New creates a History backed by the given file path.
// maxSize controls how many entries are retained.
func New(path string, maxSize int) (*History, error) {
	h := &History{path: path, maxSize: maxSize}
	if err := h.load(); err != nil {
		return nil, err
	}
	return h, nil
}

// Add appends a new entry and persists the history file.
func (h *History) Add(ports []scanner.Port) error {
	h.entries = append(h.entries, Entry{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	})
	if len(h.entries) > h.maxSize {
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
	return h.save()
}

// Entries returns a copy of all stored entries.
func (h *History) Entries() []Entry {
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Clear removes all entries and deletes the backing file.
func (h *History) Clear() error {
	h.entries = nil
	return os.Remove(h.path)
}

func (h *History) load() error {
	data, err := os.ReadFile(h.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &h.entries)
}

func (h *History) save() error {
	data, err := json.MarshalIndent(h.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(h.path, data, 0o644)
}
