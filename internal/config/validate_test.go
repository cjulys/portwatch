package config_test

import (
	"strings"
	"testing"
	"time"

	"portwatch/internal/config"
)

func TestValidateDefault(t *testing.T) {
	if err := config.Validate(config.Default()); err != nil {
		t.Errorf("default config should be valid, got: %v", err)
	}
}

func TestValidateIntervalTooShort(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 500 * time.Millisecond
	err := config.Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for short interval")
	}
	if !strings.Contains(err.Error(), "interval") {
		t.Errorf("error should mention 'interval', got: %v", err)
	}
}

func TestValidateBadAlertLevel(t *testing.T) {
	cfg := config.Default()
	cfg.AlertLevel = "critical"
	err := config.Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for unknown alert_level")
	}
	if !strings.Contains(err.Error(), "alert_level") {
		t.Errorf("error should mention 'alert_level', got: %v", err)
	}
}

func TestValidateDuplicatePort(t *testing.T) {
	cfg := config.Default()
	cfg.Ports = []uint16{80, 443, 80}
	err := config.Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for duplicate port")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Errorf("error should mention 'duplicate', got: %v", err)
	}
}

func TestValidatePortZero(t *testing.T) {
	cfg := config.Default()
	cfg.Ports = []uint16{0, 8080}
	err := config.Validate(cfg)
	if err == nil {
		t.Fatal("expected validation error for port 0")
	}
	if !strings.Contains(err.Error(), "port 0") {
		t.Errorf("error should mention 'port 0', got: %v", err)
	}
}

func TestValidateMultipleErrors(t *testing.T) {
	cfg := &config.Config{
		Interval:   0,
		AlertLevel: "bad",
		Ports:      []uint16{0},
	}
	err := config.Validate(cfg)
	if err == nil {
		t.Fatal("expected multiple validation errors")
	}
	ve, ok := err.(*config.ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errs) < 3 {
		t.Errorf("expected at least 3 errors, got %d", len(ve.Errs))
	}
}
