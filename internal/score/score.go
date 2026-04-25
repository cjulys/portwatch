// Package score computes a composite risk score for a port scan result
// by combining drift, reachability, and escalation signals into a single
// normalised value in the range [0, 100].
package score

import "math"

// Weights for each signal component. Must sum to 1.0.
const (
	weightDrift       = 0.40
	weightReachable   = 0.35
	weightEscalation  = 0.25
)

// Level describes the severity bucket derived from a Score.
type Level int

const (
	LevelLow      Level = iota // 0–33
	LevelMedium               // 34–66
	LevelHigh                 // 67–100
)

func (l Level) String() string {
	switch l {
	case LevelLow:
		return "low"
	case LevelMedium:
		return "medium"
	case LevelHigh:
		return "high"
	default:
		return "unknown"
	}
}

// Input holds the normalised [0,1] signals fed into the scorer.
type Input struct {
	// DriftScore is the fraction of ports that deviated from baseline.
	DriftScore float64
	// ReachableScore is the observed reachability ratio (0 = never reached, 1 = always).
	ReachableScore float64
	// EscalationScore is the escalation pressure normalised to [0,1]
	// (e.g. 0 = first occurrence, 1 = max escalation tier reached).
	EscalationScore float64
}

// Result is the computed composite risk assessment.
type Result struct {
	Score float64
	Level Level
}

// Compute returns a composite risk Result for the given Input.
// Each signal is clamped to [0,1] before weighting.
func Compute(in Input) Result {
	ds := clamp(in.DriftScore)
	rs := clamp(in.ReachableScore)
	es := clamp(in.EscalationScore)

	raw := (ds*weightDrift + rs*weightReachable + es*weightEscalation) * 100
	score := math.Round(raw*10) / 10 // one decimal place

	return Result{
		Score: score,
		Level: levelFor(score),
	}
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func levelFor(score float64) Level {
	switch {
	case score <= 33:
		return LevelLow
	case score <= 66:
		return LevelMedium
	default:
		return LevelHigh
	}
}
