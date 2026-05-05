package portage

import (
	"strings"
	"testing"
	"time"
)

func makeSummaryEntry(bucket string) Entry {
	return Entry{
		Protocol:  "tcp",
		Address:   "127.0.0.1",
		Port:      8080,
		FirstSeen: time.Now().Add(-24 * time.Hour),
		LastSeen:  time.Now(),
		Bucket:    bucket,
	}
}

func TestSummarizeEmpty(t *testing.T) {
	s := Summarize(nil)
	if s.Total != 0 {
		t.Fatalf("expected 0 total, got %d", s.Total)
	}
}

func TestSummarizeCountsEachBucket(t *testing.T) {
	entries := []Entry{
		makeSummaryEntry(BucketNew),
		makeSummaryEntry(BucketNew),
		makeSummaryEntry(BucketRecent),
		makeSummaryEntry(BucketEstablished),
		makeSummaryEntry(BucketLongTerm),
		makeSummaryEntry(BucketLongTerm),
	}
	s := Summarize(entries)
	if s.Total != 6 {
		t.Fatalf("expected total=6, got %d", s.Total)
	}
	if s.New != 2 {
		t.Errorf("expected new=2, got %d", s.New)
	}
	if s.Recent != 1 {
		t.Errorf("expected recent=1, got %d", s.Recent)
	}
	if s.Established != 1 {
		t.Errorf("expected established=1, got %d", s.Established)
	}
	if s.LongTerm != 2 {
		t.Errorf("expected long-term=2, got %d", s.LongTerm)
	}
}

func TestSummarizeUnknownBucketCountsOnlyTotal(t *testing.T) {
	entries := []Entry{makeSummaryEntry("unknown")}
	s := Summarize(entries)
	if s.Total != 1 {
		t.Fatalf("expected total=1, got %d", s.Total)
	}
	if s.New+s.Recent+s.Established+s.LongTerm != 0 {
		t.Error("unknown bucket should not increment named fields")
	}
}

func TestAgeSummaryString(t *testing.T) {
	s := AgeSummary{Total: 3, New: 1, Recent: 1, Established: 1}
	out := s.String()
	for _, want := range []string{"total=3", "new=1", "recent=1", "established=1", "long-term=0"} {
		if !strings.Contains(out, want) {
			t.Errorf("String() missing %q in %q", want, out)
		}
	}
}
