package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"portwatch/internal/retry"
)

func TestSuccessAfterTransientFailures(t *testing.T) {
	p := retry.Policy{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond, Factor: 1.5}
	r := retry.New(p)
	attempts := 0
	err := r.Do(context.Background(), func() error {
		attempts++
		if attempts < 4 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if attempts != 4 {
		t.Fatalf("expected 4 attempts, got %d", attempts)
	}
}

func TestOpKeyFormat(t *testing.T) {
	key := retry.OpKey("probe", "localhost:8080")
	if key != "probe:localhost:8080" {
		t.Errorf("unexpected key: %s", key)
	}
}

func TestScanKeyAndWebhookKey(t *testing.T) {
	sk := retry.ScanKey("192.168.1.1")
	wk := retry.WebhookKey("http://example.com/hook")
	if sk == wk {
		t.Error("scan and webhook keys should differ")
	}
}
