package notifier

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Level represents the severity of a notification.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event holds a single notification event.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.Port
}

// Notifier dispatches port change events to one or more handlers.
type Notifier struct {
	handlers []Handler
	out      io.Writer
}

// Handler is a function that receives a notification event.
type Handler func(e Event)

// New creates a Notifier that writes plain-text fallback output to out.
// If out is nil, os.Stderr is used.
func New(out io.Writer, handlers ...Handler) *Notifier {
	if out == nil {
		out = os.Stderr
	}
	return &Notifier{handlers: handlers, out: out}
}

// Dispatch sends an event to all registered handlers.
// If no handlers are registered the event is printed to the fallback writer.
func (n *Notifier) Dispatch(level Level, msg string, port scanner.Port) {
	e := Event{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   msg,
		Port:      port,
	}
	if len(n.handlers) == 0 {
		fmt.Fprintf(n.out, "[%s] %s %s\n", e.Level, e.Timestamp.Format(time.RFC3339), e.Message)
		return
	}
	for _, h := range n.handlers {
		h(e)
	}
}

// AddHandler appends a handler to the notifier at runtime.
func (n *Notifier) AddHandler(h Handler) {
	n.handlers = append(n.handlers, h)
}
