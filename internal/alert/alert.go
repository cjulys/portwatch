package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single notification event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier sends alerts to a configured output.
type Notifier struct {
	out io.Writer
}

// New creates a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify formats and writes alerts for each diff entry.
func (n *Notifier) Notify(diffs []scanner.Diff) {
	for _, d := range diffs {
		level := levelForDiff(d)
		msg := messageForDiff(d)
		a := Alert{
			Timestamp: time.Now(),
			Level:     level,
			Message:   msg,
		}
		fmt.Fprintf(n.out, "[%s] %s %s\n", a.Level, a.Timestamp.Format(time.RFC3339), a.Message)
	}
}

func levelForDiff(d scanner.Diff) Level {
	switch d.State {
	case scanner.StateNew:
		return LevelAlert
	case scanner.StateClosed:
		return LevelWarn
	default:
		return LevelInfo
	}
}

func messageForDiff(d scanner.Diff) string {
	switch d.State {
	case scanner.StateNew:
		return fmt.Sprintf("new port opened: %s", d.Port)
	case scanner.StateClosed:
		return fmt.Sprintf("port closed: %s", d.Port)
	default:
		return fmt.Sprintf("port state changed: %s", d.Port)
	}
}
