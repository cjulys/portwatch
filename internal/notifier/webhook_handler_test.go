package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func TestWebhookHandlerSendsPayload(t *testing.T) {
	var received WebhookPayload

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("unexpected content-type: %s", ct)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	h := NewWebhookHandler(ts.URL, &bytes.Buffer{})
	e := Event{
		Timestamp: time.Now(),
		Level:     "warn",
		Message:   "new port opened",
		Port:      scanner.Port{Port: 8080, Protocol: "tcp"},
	}
	h.Handle(e)

	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", received.Protocol)
	}
	if received.Level != "warn" {
		t.Errorf("expected level warn, got %s", received.Level)
	}
	if received.Message != "new port opened" {
		t.Errorf("unexpected message: %s", received.Message)
	}
}

func TestWebhookHandlerLogsOnBadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	var errBuf bytes.Buffer
	h := NewWebhookHandler(ts.URL, &errBuf)
	h.Handle(Event{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test",
		Port:      scanner.Port{Port: 22, Protocol: "tcp"},
	})

	if errBuf.Len() == 0 {
		t.Error("expected error output for non-2xx status, got none")
	}
}

func TestWebhookHandlerLogsOnUnreachable(t *testing.T) {
	var errBuf bytes.Buffer
	h := NewWebhookHandler("http://127.0.0.1:1", &errBuf)
	h.Handle(Event{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "unreachable",
		Port:      scanner.Port{Port: 443, Protocol: "tcp"},
	})

	if errBuf.Len() == 0 {
		t.Error("expected error output for unreachable host, got none")
	}
}
