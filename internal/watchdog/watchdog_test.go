package watchdog_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestBeatPreventsAlert(t *testing.T) {
	var buf bytes.Buffer
	wd := watchdog.New(80*time.Millisecond, &buf)
	defer wd.Stop()

	for i := 0; i < 3; i++ {
		time.Sleep(40 * time.Millisecond)
		wd.Beat()
	}
	time.Sleep(20 * time.Millisecond)

	if buf.Len() != 0 {
		t.Fatalf("expected no alert, got: %s", buf.String())
	}
}

func TestFiresWhenStalled(t *testing.T) {
	var buf bytes.Buffer
	_ = watchdog.New(50*time.Millisecond, &buf)
	// no Beat called
	time.Sleep(100 * time.Millisecond)

	if buf.Len() == 0 {
		t.Fatal("expected stall alert, got none")
	}
}

func TestStopSuppressesFire(t *testing.T) {
	var buf bytes.Buffer
	wd := watchdog.New(50*time.Millisecond, &buf)
	wd.Stop()
	time.Sleep(100 * time.Millisecond)

	if buf.Len() != 0 {
		t.Fatalf("expected no alert after Stop, got: %s", buf.String())
	}
}

func TestDefaultWriterIsStderr(t *testing.T) {
	// Just ensure New doesn't panic with nil writer.
	wd := watchdog.New(1*time.Hour, nil)
	wd.Stop()
}
