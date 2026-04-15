package notifier

import (
	"fmt"
	"io"
	"net/smtp"
	"os"
	"strings"
)

// EmailConfig holds SMTP configuration for the email handler.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// emailHandler sends alert notifications via SMTP email.
type emailHandler struct {
	cfg    EmailConfig
	writer io.Writer
}

// NewEmailHandler creates a Handler that sends emails on port change events.
// If cfg.Host is empty, messages are written to fallbackWriter (or os.Stderr).
func NewEmailHandler(cfg EmailConfig, fallbackWriter io.Writer) Handler {
	if fallbackWriter == nil {
		fallbackWriter = os.Stderr
	}
	return &emailHandler{cfg: cfg, writer: fallbackWriter}
}

func (h *emailHandler) Handle(e Event) {
	subject := fmt.Sprintf("[portwatch] %s", e.Message)
	body := fmt.Sprintf("Time: %s\nLevel: %s\nMessage: %s\n",
		e.Timestamp.Format("2006-01-02 15:04:05"),
		e.Level,
		e.Message,
	)

	if h.cfg.Host == "" || len(h.cfg.To) == 0 {
		fmt.Fprintf(h.writer, "[email_handler] no SMTP host configured, dropping: %s\n", subject)
		return
	}

	addr := fmt.Sprintf("%s:%d", h.cfg.Host, h.cfg.Port)
	auth := smtp.PlainAuth("", h.cfg.Username, h.cfg.Password, h.cfg.Host)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s",
		strings.Join(h.cfg.To, ", "),
		h.cfg.From,
		subject,
		body,
	))

	if err := smtp.SendMail(addr, auth, h.cfg.From, h.cfg.To, msg); err != nil {
		fmt.Fprintf(h.writer, "[email_handler] failed to send email: %v\n", err)
	}
}
