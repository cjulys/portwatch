package score

import (
	"testing"
)

func TestComputeAllZeroIsLow(t *testing.T) {
	r := Compute(Input{})
	if r.Score != 0 {
		t.Fatalf("expected 0, got %v", r.Score)
	}
	if r.Level != LevelLow {
		t.Fatalf("expected low, got %v", r.Level)
	}
}

func TestComputeAllOneIsHigh(t *testing.T) {
	r := Compute(Input{DriftScore: 1, ReachableScore: 1, EscalationScore: 1})
	if r.Score != 100 {
		t.Fatalf("expected 100, got %v", r.Score)
	}
	if r.Level != LevelHigh {
		t.Fatalf("expected high, got %v", r.Level)
	}
}

func TestComputeWeightsApplied(t *testing.T) {
	// Only drift signal active.
	r := Compute(Input{DriftScore: 1, ReachableScore: 0, EscalationScore: 0})
	want := weightDrift * 100
	if r.Score != want {
		t.Fatalf("expected %v, got %v", want, r.Score)
	}
}

func TestComputeClampsAboveOne(t *testing.T) {
	r := Compute(Input{DriftScore: 5, ReachableScore: 5, EscalationScore: 5})
	if r.Score != 100 {
		t.Fatalf("expected 100 after clamping, got %v", r.Score)
	}
}

func TestComputeClampsBelowZero(t *testing.T) {
	r := Compute(Input{DriftScore: -1, ReachableScore: -2, EscalationScore: -3})
	if r.Score != 0 {
		t.Fatalf("expected 0 after clamping, got %v", r.Score)
	}
}

func TestLevelBoundaries(t *testing.T) {
	cases := []struct {
		score float64
		want  Level
	}{
		{0, LevelLow},
		{33, LevelLow},
		{34, LevelMedium},
		{66, LevelMedium},
		{67, LevelHigh},
		{100, LevelHigh},
	}
	for _, tc := range cases {
		got := levelFor(tc.score)
		if got != tc.want {
			t.Errorf("levelFor(%v) = %v, want %v", tc.score, got, tc.want)
		}
	}
}

func TestLevelString(t *testing.T) {
	if LevelLow.String() != "low" {
		t.Errorf("unexpected: %v", LevelLow.String())
	}
	if LevelMedium.String() != "medium" {
		t.Errorf("unexpected: %v", LevelMedium.String())
	}
	if LevelHigh.String() != "high" {
		t.Errorf("unexpected: %v", LevelHigh.String())
	}
}
