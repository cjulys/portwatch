package anomaly

import (
	"math"
	"testing"
)

// feed populates d with n identical samples to establish a stable baseline.
func feed(d *Detector, value float64, n int) {
	for i := 0; i < n; i++ {
		d.Evaluate(int(value))
	}
}

func TestNoEventBeforeMinSamples(t *testing.T) {
	d := New(60, nil)
	for i := 0; i < 9; i++ {
		if ev := d.Evaluate(10); ev != nil {
			t.Fatalf("expected nil before 10 samples, got %v", ev)
		}
	}
}

func TestNoEventWhenStable(t *testing.T) {
	d := New(60, nil)
	feed(d, 20, 20)
	if ev := d.Evaluate(20); ev != nil {
		t.Fatalf("expected nil for stable series, got %v", ev)
	}
}

func TestSevereAnomalyDetected(t *testing.T) {
	d := New(60, nil)
	// baseline: 20 open ports, low variance
	for i := 0; i < 30; i++ {
		d.Evaluate(20)
	}
	d.Evaluate(19)
	d.Evaluate(21)

	// spike far from mean
	ev := d.Evaluate(200)
	if ev == nil {
		t.Fatal("expected anomaly event, got nil")
	}
	if ev.Level != LevelSevere {
		t.Errorf("expected severe, got %s", ev.Level)
	}
	if ev.OpenCount != 200 {
		t.Errorf("expected OpenCount=200, got %d", ev.OpenCount)
	}
	if ev.ZScore <= 3 {
		t.Errorf("expected z > 3, got %.2f", ev.ZScore)
	}
	if ev.At.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestMildAnomalyLevel(t *testing.T) {
	// Construct a detector with known mean/stddev via direct sample injection.
	d := New(60, nil)
	// mean=10, stddev≈1 after many samples of 10 with small jitter
	for i := 0; i < 28; i++ {
		d.Evaluate(10)
	}
	d.Evaluate(9)
	d.Evaluate(11)

	// 1.5σ above mean should be mild
	// stddev ≈ sqrt((28*0 + 1 + 1)/30) ≈ 0.365 — too small; use a wider spread.
	// Instead verify that levelFor mapping is correct independently.
	if levelFor(1.5) != LevelMild {
		t.Errorf("expected mild for z=1.5")
	}
	if levelFor(2.5) != LevelModerate {
		t.Errorf("expected moderate for z=2.5")
	}
	if levelFor(3.5) != LevelSevere {
		t.Errorf("expected severe for z=3.5")
	}
	if levelFor(0.5) != LevelNone {
		t.Errorf("expected none for z=0.5")
	}
}

func TestResetClearsHistory(t *testing.T) {
	d := New(60, nil)
	feed(d, 10, 30)
	d.Reset()
	// after reset fewer than 10 samples → no event
	for i := 0; i < 9; i++ {
		if ev := d.Evaluate(100); ev != nil {
			t.Fatalf("expected nil after reset, got %v", ev)
		}
	}
}

func TestZeroStdDevReturnsNil(t *testing.T) {
	d := New(60, nil)
	// all identical → stddev == 0 → no event
	for i := 0; i < 20; i++ {
		if ev := d.Evaluate(5); ev != nil {
			t.Fatalf("expected nil for zero stddev, got %v", ev)
		}
	}
}

func TestLevelString(t *testing.T) {
	cases := map[Level]string{
		LevelNone:     "none",
		LevelMild:     "mild",
		LevelModerate: "moderate",
		LevelSevere:   "severe",
	}
	for lvl, want := range cases {
		if got := lvl.String(); got != want {
			t.Errorf("Level(%d).String() = %q, want %q", lvl, got, want)
		}
	}
}

func TestMaxSamplesWindowEnforced(t *testing.T) {
	d := New(20, nil)
	// fill with 100 → window keeps last 20
	for i := 0; i < 50; i++ {
		d.Evaluate(100)
	}
	d.mu.Lock()
	if len(d.samples) > 20 {
		t.Errorf("expected at most 20 samples, got %d", len(d.samples))
	}
	d.mu.Unlock()
}

func TestStatsFunction(t *testing.T) {
	s := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	mean, stddev := stats(s)
	if math.Abs(mean-5.0) > 1e-9 {
		t.Errorf("mean: got %.4f, want 5.0", mean)
	}
	if math.Abs(stddev-2.0) > 1e-9 {
		t.Errorf("stddev: got %.4f, want 2.0", stddev)
	}
}
