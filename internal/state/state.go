package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a recorded set of port states at a point in time.
type Snapshot struct {
	Timestamp time.Time      `json:"timestamp"`
	Ports     []scanner.Port `json:"ports"`
}

// Store persists and retrieves port snapshots from disk.
type Store struct {
	path string
}

// New creates a Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes the current snapshot to disk, overwriting any previous state.
func (s *Store) Save(ports []scanner.Port) error {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

// Load reads the last saved snapshot from disk.
// If the file does not exist, an empty snapshot is returned without error.
func (s *Store) Load() (Snapshot, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return Snapshot{}, nil
	}
	if err != nil {
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}

// Clear removes the persisted state file from disk.
func (s *Store) Clear() error {
	err := os.Remove(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
