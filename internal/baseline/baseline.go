package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline represents a known-good snapshot of open ports.
type Baseline struct {
	mu      sync.RWMutex
	path    string
	Ports   []scanner.Port `json:"ports"`
	SavedAt time.Time      `json:"saved_at"`
}

// New creates a Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{path: path}
}

// Set replaces the current baseline with the provided ports and persists it.
func (b *Baseline) Set(ports []scanner.Port) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Ports = ports
	b.SavedAt = time.Now().UTC()
	return b.save()
}

// Get returns the current baseline ports.
func (b *Baseline) Get() []scanner.Port {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Ports
}

// Load reads the baseline from disk. If the file does not exist the baseline
// is left empty and no error is returned.
func (b *Baseline) Load() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	data, err := os.ReadFile(b.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, b)
}

// Clear removes the baseline file and resets in-memory state.
func (b *Baseline) Clear() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Ports = nil
	b.SavedAt = time.Time{}
	if err := os.Remove(b.path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// IsEmpty reports whether no baseline has been recorded yet.
func (b *Baseline) IsEmpty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.Ports) == 0
}

func (b *Baseline) save() error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o600)
}
