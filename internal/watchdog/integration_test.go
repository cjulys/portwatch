package watchdog_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestMultipleBeatsKeepAlive(t *testing.T) {
	var buf bytes.Buffer
	wd := watchdog.New(60*time.Millisecond, &buf)
	defer wd.Stop()

	tick := time.NewTicker(20 * time.Millisecond)
	defer tick.Stop()
	done := time.After(150 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			wd.Beat()
		case <-done:
			goto check
		}
	}
check:
	if buf.Len() != 0 {
		t.Fatalf("unexpected stall alert: %s", buf.String())
	}
}

func TestAlertMessageContainsKeyword(t *testing.T) {
	var buf bytes.Buffer
	_ = watchdog.New(30*time.Millisecond, &buf)
	time.Sleep(80 * time.Millisecond)

	if !strings.Contains(buf.String(), "stalled") {
		t.Fatalf("alert message missing 'stalled': %q", buf.String())
	}
}
