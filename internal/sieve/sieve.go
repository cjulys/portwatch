// Package sieve provides a probabilistic duplicate-suppression filter
// using a fixed-size bitset (Bloom-filter style) keyed on port fingerprints.
// It is intentionally lossy: false positives are possible, false negatives
// are not. Use it to avoid re-alerting on ports that were already reported
// within the current scan cycle.
package sieve

import (
	"hash/fnv"
	"sync"
)

const defaultBuckets = 1024

// Sieve is a thread-safe probabilistic seen-set.
type Sieve struct {
	mu      sync.Mutex
	bits    []bool
	buckets uint32
}

// New returns a Sieve with the given number of buckets.
// If buckets is zero, defaultBuckets is used.
func New(buckets uint32) *Sieve {
	if buckets == 0 {
		buckets = defaultBuckets
	}
	return &Sieve{
		bits:    make([]bool, buckets),
		buckets: buckets,
	}
}

// TestAndSet returns true if key was already present, then marks it as seen.
// Subsequent calls with the same key always return true until Reset is called.
func (s *Sieve) TestAndSet(key string) bool {
	idx := s.index(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	alreadySeen := s.bits[idx]
	s.bits[idx] = true
	return alreadySeen
}

// Seen reports whether key has been marked without modifying state.
func (s *Sieve) Seen(key string) bool {
	idx := s.index(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.bits[idx]
}

// Reset clears all bits, making every key appear unseen again.
func (s *Sieve) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.bits {
		s.bits[i] = false
	}
}

// Len returns the number of buckets in the sieve.
func (s *Sieve) Len() int { return int(s.buckets) }

func (s *Sieve) index(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32() % s.buckets
}
