package jitter

import (
	"testing"
	"time"
)

const base = 10 * time.Second

func TestZeroFactorReturnsBase(t *testing.T) {
	j := New(0)
	for i := 0; i < 20; i++ {
		if got := j.Apply(base); got != base {
			t.Fatalf("expected %v, got %v", base, got)
		}
	}
}

func TestNegativeFactorClampsToZero(t *testing.T) {
	j := New(-0.5)
	if j.Factor() != 0 {
		t.Fatalf("expected factor 0, got %v", j.Factor())
	}
}

func TestFactorAboveOneClampsToOne(t *testing.T) {
	j := New(1.5)
	if j.Factor() != 1 {
		t.Fatalf("expected factor 1, got %v", j.Factor())
	}
}

func TestApplyStaysWithinBounds(t *testing.T) {
	j := New(0.2)
	min := time.Duration(float64(base) * 0.8)
	max := time.Duration(float64(base) * 1.2)
	for i := 0; i < 500; i++ {
		got := j.Apply(base)
		if got < min || got > max {
			t.Fatalf("value %v out of [%v, %v]", got, min, max)
		}
	}
}

func TestApplyZeroDurationReturnsZero(t *testing.T) {
	j := New(0.5)
	if got := j.Apply(0); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestApplyProducesVariation(t *testing.T) {
	j := New(0.3)
	seen := map[time.Duration]bool{}
	for i := 0; i < 100; i++ {
		seen[j.Apply(base)] = true
	}
	if len(seen) < 2 {
		t.Fatal("expected variation in jittered values")
	}
}
