package grace

import (
	"bytes"
	"context"
	"syscall"
	"testing"
	"time"
)

func TestAcquireReleaseDrains(t *testing.T) {
	c := New(time.Second, nil)
	c.Acquire()
	c.Acquire()
	c.Release()
	c.Release()
	select {
	case <-c.drain:
	// ok
	default:
		t.Fatal("expected drain signal after all releases")
	}
}

func TestReleaseWithoutAcquireIsNoop(t *testing.T) {
	c := New(time.Second, nil)
	c.Release() // must not panic or underflow
}

func TestWaitCancelledByParent(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	c := New(time.Second, nil)
	ctx := c.Wait(parent)
	cancel()
	select {
	case <-ctx.Done():
	case <-time.After(500 * time.Millisecond):
		t.Fatal("context should be done after parent cancel")
	}
}

func TestWaitDrainsBeforeCancel(t *testing.T) {
	var buf bytes.Buffer
	c := New(2*time.Second, &buf)
	c.Acquire()

	parent := context.Background()
	ctx := c.Wait(parent)

	go func() {
		time.Sleep(50 * time.Millisecond)
		c.Release()
	}()

	// send signal to self
	syscall.Kill(syscall.Getpid(), syscall.SIGINT) //nolint:errcheck

	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("context should be done after drain")
	}
}

func TestTimeoutWritesFallback(t *testing.T) {
	var buf bytes.Buffer
	c := New(50*time.Millisecond, &buf)
	c.Acquire() // never released

	parent := context.Background()
	ctx := c.Wait(parent)

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM) //nolint:errcheck

	select {
	case <-ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("context should be done after timeout")
	}
	if buf.Len() == 0 {
		t.Fatal("expected fallback message on timeout")
	}
}
