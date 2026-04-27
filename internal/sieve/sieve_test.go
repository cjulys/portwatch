package sieve

import (
	"fmt"
	"testing"
)

func TestTestAndSetFirstCallReturnsFalse(t *testing.T) {
	s := New(256)
	if s.TestAndSet("tcp:22") {
		t.Fatal("expected false on first TestAndSet")
	}
}

func TestTestAndSetSecondCallReturnsTrue(t *testing.T) {
	s := New(256)
	s.TestAndSet("tcp:22")
	if !s.TestAndSet("tcp:22") {
		t.Fatal("expected true on second TestAndSet for same key")
	}
}

func TestSeenReflectsState(t *testing.T) {
	s := New(256)
	if s.Seen("udp:53") {
		t.Fatal("key should not be seen before TestAndSet")
	}
	s.TestAndSet("udp:53")
	if !s.Seen("udp:53") {
		t.Fatal("key should be seen after TestAndSet")
	}
}

func TestResetClearsAllBits(t *testing.T) {
	s := New(256)
	keys := []string{"tcp:80", "tcp:443", "udp:123"}
	for _, k := range keys {
		s.TestAndSet(k)
	}
	s.Reset()
	for _, k := range keys {
		if s.Seen(k) {
			t.Fatalf("key %q should not be seen after Reset", k)
		}
	}
}

func TestDifferentKeysAreIndependent(t *testing.T) {
	s := New(1024)
	s.TestAndSet("tcp:22")
	if s.Seen("tcp:23") {
		t.Fatal("unrelated key should not be seen")
	}
}

func TestDefaultBucketsUsedWhenZero(t *testing.T) {
	s := New(0)
	if s.Len() != defaultBuckets {
		t.Fatalf("expected %d buckets, got %d", defaultBuckets, s.Len())
	}
}

func TestConcurrentAccessDoesNotPanic(t *testing.T) {
	s := New(512)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func(n int) {
			key := fmt.Sprintf("tcp:%d", n)
			s.TestAndSet(key)
			s.Seen(key)
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 50; i++ {
		<-done
	}
}
