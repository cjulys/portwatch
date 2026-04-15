package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"portwatch/internal/config"
)

func TestDefaultValues(t *testing.T) {
	cfg := config.Default()
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected 30s interval, got %v", cfg.Interval)
	}
	if cfg.AlertLevel != "info" {
		t.Errorf("expected alert_level 'info', got %q", cfg.AlertLevel)
	}
	if len(cfg.Ports) != 0 {
		t.Errorf("expected empty ports slice, got %v", cfg.Ports)
	}
}

func TestLoadNonExistentFileReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("/tmp/portwatch_no_such_file_xyz.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected default interval, got %v", cfg.Interval)
	}
}

func TestLoadEmptyPathReturnsDefaults(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AlertLevel != "info" {
		t.Errorf("expected default alert_level, got %q", cfg.AlertLevel)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	orig := &config.Config{
		Interval:   10 * time.Second,
		Ports:      []uint16{80, 443, 8080},
		AlertLevel: "warn",
		LogFile:    "/var/log/portwatch.log",
	}

	if err := config.Save(orig, tmp.Name()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := config.Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Interval != orig.Interval {
		t.Errorf("interval mismatch: got %v want %v", loaded.Interval, orig.Interval)
	}
	if loaded.AlertLevel != orig.AlertLevel {
		t.Errorf("alert_level mismatch: got %q want %q", loaded.AlertLevel, orig.AlertLevel)
	}
	if loaded.LogFile != orig.LogFile {
		t.Errorf("log_file mismatch: got %q want %q", loaded.LogFile, orig.LogFile)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("ports length mismatch: got %d want %d", len(loaded.Ports), len(orig.Ports))
	}
}

func TestSaveProducesValidJSON(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := config.Save(config.Default(), tmp.Name()); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	data, _ := os.ReadFile(tmp.Name())
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Errorf("saved file is not valid JSON: %v", err)
	}
}
