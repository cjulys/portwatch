package notifier

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// LogFormat selects the output format for LogHandler.
type LogFormat string

const (
	LogFormatText LogFormat = "text"
	LogFormatJSON LogFormat = "json"
)

type jsonEvent struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Protocol  string `json:"protocol"`
	Port      uint16 `json:"port"`
}

// LogHandler returns a Handler that writes events to w in the chosen format.
// If w is nil, os.Stdout is used.
func LogHandler(w io.Writer, format LogFormat) Handler {
	if w == nil {
		w = os.Stdout
	}
	return func(e Event) {
		switch format {
		case LogFormatJSON:
			je := jsonEvent{
				Timestamp: e.Timestamp.Format(time.RFC3339),
				Level:     string(e.Level),
				Message:   e.Message,
				Protocol:  e.Port.Protocol,
				Port:      e.Port.Number,
			}
			data, err := json.Marshal(je)
			if err != nil {
				fmt.Fprintf(w, "[error] failed to marshal event: %v\n", err)
				return
			}
			fmt.Fprintln(w, string(data))
		default:
			fmt.Fprintf(w, "[%s] %s | %s | %s/%d\n",
				e.Level,
				e.Timestamp.Format(time.RFC3339),
				e.Message,
				e.Port.Protocol,
				e.Port.Number,
			)
		}
	}
}
