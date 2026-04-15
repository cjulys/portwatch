package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// Interval between port scans.
	Interval time.Duration `json:"interval"`
	// Ports to watch; if empty, all detected ports are monitored.
	Ports []uint16 `json:"ports"`
	// AlertLevel controls minimum severity to emit ("info", "warn", "error").
	AlertLevel string `json:"alert_level"`
	// LogFile is an optional path to write alerts to (stdout if empty).
	LogFile string `json:"log_file"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Interval:   30 * time.Second,
		Ports:      []uint16{},
		AlertLevel: "info",
		LogFile:    "",
	}
}

// Load reads a JSON config file from path and merges it over the defaults.
// If path is empty or the file does not exist, the defaults are returned.
func Load(path string) (*Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config as indented JSON to path.
func Save(cfg *Config, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
