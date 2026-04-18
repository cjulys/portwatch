package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Format controls the output format of the reporter.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Entry represents a single snapshot report entry.
type Entry struct {
	Timestamp time.Time      `json:"timestamp"`
	Ports     []scanner.Port `json:"ports"`
	Count     int            `json:"count"`
}

// Reporter writes periodic port snapshot reports to a writer.
type Reporter struct {
	w      io.Writer
	format Format
}

// New creates a Reporter that writes to w in the given format.
// If w is nil, os.Stdout is used.
func New(w io.Writer, format Format) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	if format == "" {
		format = FormatText
	}
	return &Reporter{w: w, format: format}
}

// Report writes a snapshot of the current ports to the underlying writer.
func (r *Reporter) Report(ports []scanner.Port) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Ports:     ports,
		Count:     len(ports),
	}
	switch r.format {
	case FormatJSON:
		return r.writeJSON(entry)
	default:
		return r.writeText(entry)
	}
}

func (r *Reporter) writeText(e Entry) error {
	_, err := fmt.Fprintf(r.w, "[%s] open ports (%d):\n",
		e.Timestamp.Format(time.RFC3339), e.Count)
	if err != nil {
		return err
	}
	for _, p := range e.Ports {
		_, err = fmt.Fprintf(r.w, "  %s\n", p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Reporter) writeJSON(e Entry) error {
	enc := json.NewEncoder(r.w)
	return enc.Encode(e)
}

// SetFormat updates the output format used by the reporter.
func (r *Reporter) SetFormat(format Format) {
	if format == "" {
		return
	}
	r.format = format
}

// Format returns the current output format of the reporter.
func (r *Reporter) Format() Format {
	return r.format
}
