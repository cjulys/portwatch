package config

import (
	"fmt"
	"time"
)

// ValidationError holds a list of problems found in a Config.
type ValidationError struct {
	Problems []string
}

func (e *ValidationError) Error() string {
	if len(e.Problems) == 1 {
		return "config: " + e.Problems[0]
	}
	msg := fmt.Sprintf("config: %d validation errors", len(e.Problems))
	for _, p := range e.Problems {
		msg += "\n  - " + p
	}
	return msg
}

const minInterval = 5 * time.Second

var validAlertLevels = map[string]bool{
	"info":  true,
	"warn":  true,
	"error": true,
}

// Validate checks cfg for logical errors and returns a *ValidationError
// (wrapped as error) if any problems are found, or nil when cfg is valid.
func Validate(cfg *Config) error {
	var problems []string

	if cfg.Interval < minInterval {
		problems = append(problems,
			fmt.Sprintf("interval %s is below minimum %s", cfg.Interval, minInterval))
	}

	if !validAlertLevels[cfg.AlertLevel] {
		problems = append(problems,
			fmt.Sprintf("alert_level %q is not valid (want info|warn|error)", cfg.AlertLevel))
	}

	seen := make(map[int]bool, len(cfg.WatchPorts))
	for _, p := range cfg.WatchPorts {
		if p == 0 {
			problems = append(problems, "watch_ports contains port 0 which is not allowed")
			continue
		}
		if p < 0 || p > 65535 {
			problems = append(problems, fmt.Sprintf("watch_ports contains invalid port %d (must be 1-65535)", p))
			continue
		}
		if seen[p] {
			problems = append(problems, fmt.Sprintf("watch_ports contains duplicate port %d", p))
		} else {
			seen[p] = true
		}
	}

	if len(problems) == 0 {
		return nil
	}
	return &ValidationError{Problems: problems}
}
