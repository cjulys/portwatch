package config

import (
	"errors"
	"time"
)

var validAlertLevels = map[string]struct{}{
	"info":  {},
	"warn":  {},
	"error": {},
}

// ValidationError collects one or more validation failures.
type ValidationError struct {
	Errs []error
}

func (v *ValidationError) Error() string {
	msg := "config validation failed:"
	for _, e := range v.Errs {
		msg += " " + e.Error() + ";"
	}
	return msg
}

// Validate checks cfg for invalid or contradictory values.
// It returns a *ValidationError listing all problems, or nil if cfg is valid.
func Validate(cfg *Config) error {
	var errs []error

	if cfg.Interval < time.Second {
		errs = append(errs, errors.New("interval must be at least 1s"))
	}

	if _, ok := validAlertLevels[cfg.AlertLevel]; !ok {
		errs = append(errs, errors.New("alert_level must be one of: info, warn, error"))
	}

	seen := make(map[uint16]bool)
	for _, p := range cfg.Ports {
		if p == 0 {
			errs = append(errs, errors.New("port 0 is not a valid watch target"))
			continue
		}
		if seen[p] {
			errs = append(errs, errors.New("duplicate port in watch list"))
		}
		seen[p] = true
	}

	if len(errs) > 0 {
		return &ValidationError{Errs: errs}
	}
	return nil
}
