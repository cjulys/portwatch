package portage

import (
	"testing"
)

func TestAgeKeyFormat(t *testing.T) {
	got := AgeKey("tcp", 443)
	want := "tcp:443"
	if got != want {
		t.Fatalf("AgeKey = %q, want %q", got, want)
	}
}

func TestAgeKeyDistinguishesProtocols(t *testing.T) {
	tcp := AgeKey("tcp", 53)
	udp := AgeKey("udp", 53)
	if tcp == udp {
		t.Fatal("expected tcp and udp keys to differ")
	}
}

func TestBucketKeyNew(t *testing.T) {
	key := BucketKey("tcp", 80, 60) // 60 seconds = new
	want := "tcp:80:new"
	if key != want {
		t.Fatalf("BucketKey = %q, want %q", key, want)
	}
}

func TestBucketKeyRecent(t *testing.T) {
	key := BucketKey("tcp", 80, 7200) // 2 hours = recent
	want := "tcp:80:recent"
	if key != want {
		t.Fatalf("BucketKey = %q, want %q", key, want)
	}
}

func TestBucketKeyEstablished(t *testing.T) {
	key := BucketKey("tcp", 22, 172800) // 2 days = established
	want := "tcp:22:established"
	if key != want {
		t.Fatalf("BucketKey = %q, want %q", key, want)
	}
}

func TestBucketKeyOld(t *testing.T) {
	key := BucketKey("udp", 161, 700000) // >7 days = old
	want := "udp:161:old"
	if key != want {
		t.Fatalf("BucketKey = %q, want %q", key, want)
	}
}

func TestBucketLabelBoundaries(t *testing.T) {
	cases := []struct {
		secs  int64
		want  string
	}{
		{0, "new"},
		{3599, "new"},
		{3600, "recent"},
		{86399, "recent"},
		{86400, "established"},
		{604799, "established"},
		{604800, "old"},
		{999999, "old"},
	}
	for _, tc := range cases {
		got := bucketLabel(tc.secs)
		if got != tc.want {
			t.Errorf("bucketLabel(%d) = %q, want %q", tc.secs, got, tc.want)
		}
	}
}
