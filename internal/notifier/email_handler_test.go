package notifier

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestEmailHandlerNoHostWritesFallback(t *testing.T) {
	var buf bytes.Buffer
	h := NewEmailHandler(EmailConfig{}, &buf)

	e := Event{
		Message:   "port 8080/tcp opened",
		Level:     "warn",
		Timestamp: time.Now(),
	}
	h.Handle(e)

	out := buf.String()
	if !strings.Contains(out, "no SMTP host configured") {
		t.Errorf("expected fallback message, got: %q", out)
	}
	if !strings.Contains(out, "port 8080/tcp opened") {
		t.Errorf("expected event message in fallback output, got: %q", out)
	}
}

func TestEmailHandlerNoRecipientsWritesFallback(t *testing.T) {
	var buf bytes.Buffer
	cfg := EmailConfig{
		Host: "smtp.example.com",
		Port: 587,
		To:   []string{},
	}
	h := NewEmailHandler(cfg, &buf)

	e := Event{
		Message:   "port 22/tcp closed",
		Level:     "warn",
		Timestamp: time.Now(),
	}
	h.Handle(e)

	out := buf.String()
	ifno SMTP host configured") {
		t.Errorf("expected fallback message, got: %q", out)
	}
}

func TestEmailHandlerDefaultFallbackStderr(t *testing.T) {
	// Passing nil writer should not panic; it defaults to os.Stderr.
	h := NewEmailHandler(EmailConfig{}, nil)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestEmailHandlerUnreachableHostWritesFallback(t *testing.T) {
	var buf bytes.Buffer
	cfg := EmailConfig{
		Host:     "127.0.0.1",
		Port:     59999, // nothing listening here
		Username: "user",
		Password: "pass",
		From:     "portwatch@localhost",
		To:       []string{"admin@localhost"},
	}
	h := NewEmailHandler(cfg, &buf)

	e := Event{
		Message:   "port 443/tcp opened",
		Level:     "warn",
		Timestamp: time.Now(),
	}
	h.Handle(e)

	out := buf.String()
	if !strings.Contains(out, "failed to send email") {
		t.Errorf("expected send-failure message, got: %q", out)
	}
}
