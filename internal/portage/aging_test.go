package portage

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeAgePort(port int, proto string) scanner.Port {
	return scanner.Port{Port: port, Protocol: proto, Address: "127.0.0.1"}
}

func TestCategorizeNew(t *testing.T) {
	now := time.Now()
	cat := Categorize(now.Add(-30*time.Minute), now)
	if cat != AgeNew {
		t.Fatalf("expected AgeNew, got %s", cat)
	}
}

func TestCategorizeRecent(t *testing.T) {
	now := time.Now()
	cat := Categorize(now.Add(-2*time.Hour), now)
	if cat != AgeRecent {
		t.Fatalf("expected AgeRecent, got %s", cat)
	}
}

func TestCategorizeEstablished(t *testing.T) {
	now := time.Now()
	cat := Categorize(now.Add(-3*24*time.Hour), now)
	if cat != AgeEstablished {
		t.Fatalf("expected AgeEstablished, got %s", cat)
	}
}

func TestCategorizeLongTerm(t *testing.T) {
	now := time.Now()
	cat := Categorize(now.Add(-10*24*time.Hour), now)
	if cat != AgeLongTerm {
		t.Fatalf("expected AgeLongTerm, got %s", cat)
	}
}

func TestAgeCategoryString(t *testing.T) {
	cases := []struct {
		cat  AgeCategory
		want string
	}{
		{AgeNew, "new"},
		{AgeRecent, "recent"},
		{AgeEstablished, "established"},
		{AgeLongTerm, "long-term"},
	}
	for _, tc := range cases {
		if got := tc.cat.String(); got != tc.want {
			t.Errorf("AgeCategory(%d).String() = %q, want %q", tc.cat, got, tc.want)
		}
	}
}

func TestEnrichAllUsesTrackerFirstSeen(t *testing.T) {
	dir := t.TempDir()
	tr := New(dir)
	now := time.Now()
	old := now.Add(-48 * time.Hour)

	p := makeAgePort(443, "tcp")
	// Manually inject a first-seen entry by calling Update with an old clock.
	tr.clockFn = func() time.Time { return old }
	tr.Update([]scanner.Port{p}, nil)
	tr.clockFn = func() time.Time { return now }

	agedPorts := EnrichAll(tr, []scanner.Port{p}, now)
	if len(agedPorts) != 1 {
		t.Fatalf("expected 1 aged port, got %d", len(agedPorts))
	}
	if agedPorts[0].Category != AgeEstablished {
		t.Errorf("expected AgeEstablished, got %s", agedPorts[0].Category)
	}
}

func TestEnrichAllFallsBackToNowForUnknownPort(t *testing.T) {
	dir := t.TempDir()
	tr := New(dir)
	now := time.Now()
	p := makeAgePort(9999, "tcp")

	agedPorts := EnrichAll(tr, []scanner.Port{p}, now)
	if len(agedPorts) != 1 {
		t.Fatalf("expected 1 aged port, got %d", len(agedPorts))
	}
	if agedPorts[0].Category != AgeNew {
		t.Errorf("expected AgeNew for unknown port, got %s", agedPorts[0].Category)
	}
}
