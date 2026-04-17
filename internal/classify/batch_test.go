package classify_test

import (
	"testing"

	"github.com/user/portwatch/internal/classify"
)

func TestBatchCountsCritical(t *testing.T) {
	c := classify.New([]uint16{22})
	d := classify.DiffInput{
		Opened: []scanner.Port{makePort(22, "tcp"), makePort(8080, "tcp")},
	}
	br := c.Batch(d)
	if br.Critical != 1 {
		t.Fatalf("expected 1 critical, got %d", br.Critical)
	}
	if br.Warnings != 1 {
		t.Fatalf("expected 1 warning, got %d", br.Warnings)
	}
}

func TestBatchCountsInfo(t *testing.T) {
	c := classify.New(nil)
	d := classify.DiffInput{
		Closed: []scanner.Port{makePort(80, "tcp"), makePort(443, "tcp")},
	}
	br := c.Batch(d)
	if br.Info != 2 {
		t.Fatalf("expected 2 info, got %d", br.Info)
	}
}

func TestBatchResultsLength(t *testing.T) {
	c := classify.New(nil)
	d := classify.DiffInput{
		Opened: []scanner.Port{makePort(9000, "tcp")},
		Closed: []scanner.Port{makePort(9001, "tcp"), makePort(9002, "tcp")},
	}
	br := c.Batch(d)
	if len(br.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(br.Results))
	}
}

func TestBatchEmptyInputReturnsZeroCounts(t *testing.T) {
	c := classify.New(nil)
	br := c.Batch(classify.DiffInput{})
	if br.Critical != 0 || br.Warnings != 0 || br.Info != 0 {
		t.Fatal("expected all zero counts for empty input")
	}
	if len(br.Results) != 0 {
		t.Fatal("expected empty results")
	}
}
