// Package snapshot captures and compares point-in-time port scan results.
package snapshot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a timestamped set of scanned ports.
type Snapshot struct {
	CapturedAt time.Time      `json:"captured_at"`
	Ports      []scanner.Port `json:"ports"`
}

// Take creates a new Snapshot from the given ports.
func Take(ports []scanner.Port) Snapshot {
	return Snapshot{
		CapturedAt: time.Now().UTC(),
		Ports:      ports,
	}
}

// SaveToFile writes the snapshot as JSON to the given path.
func SaveToFile(s Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(s)
}

// LoadFromFile reads a Snapshot from a JSON file.
func LoadFromFile(path string) (Snapshot, error) {
	var s Snapshot
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, nil
		}
		return s, err
	}
	defer f.Close()
	return s, json.NewDecoder(f).Decode(&s)
}
