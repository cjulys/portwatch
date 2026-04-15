package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to the webhook endpoint.
type WebhookPayload struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
}

// WebhookHandler sends notifier events to an HTTP endpoint as JSON.
type WebhookHandler struct {
	URL    string
	client *http.Client
	out    io.Writer
}

// NewWebhookHandler creates a WebhookHandler that posts to url.
// errOut receives diagnostic messages on delivery failure.
func NewWebhookHandler(url string, errOut io.Writer) *WebhookHandler {
	return &WebhookHandler{
		URL:    url,
		client: &http.Client{Timeout: 5 * time.Second},
		out:    errOut,
	}
}

// Handle implements the Handler interface expected by the Notifier.
func (w *WebhookHandler) Handle(e Event) {
	payload := WebhookPayload{
		Timestamp: e.Timestamp,
		Level:     e.Level,
		Message:   e.Message,
		Port:      e.Port.Port,
		Protocol:  e.Port.Protocol,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(w.out, "webhook: marshal error: %v\n", err)
		return
	}

	resp, err := w.client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(w.out, "webhook: post error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(w.out, "webhook: unexpected status %d for %s\n", resp.StatusCode, w.URL)
	}
}
