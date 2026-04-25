package suppress

import (
	"testing"
)

func TestKeyFormat(t *testing.T) {
	got := Key("tcp", "127.0.0.1", 8080)
	want := "tcp:127.0.0.1:8080"
	if got != want {
		t.Fatalf("Key() = %q; want %q", got, want)
	}
}

func TestKeyDistinguishesProtocols(t *testing.T) {
	a := Key("tcp", "0.0.0.0", 443)
	b := Key("udp", "0.0.0.0", 443)
	if a == b {
		t.Fatal("expected tcp and udp keys to differ")
	}
}

func TestKeyDistinguishesPorts(t *testing.T) {
	a := Key("tcp", "0.0.0.0", 80)
	b := Key("tcp", "0.0.0.0", 443)
	if a == b {
		t.Fatal("expected port 80 and 443 keys to differ")
	}
}

func TestWildcardKeyUsesAsterisk(t *testing.T) {
	got := WildcardKey("tcp", 22)
	want := "tcp:*:22"
	if got != want {
		t.Fatalf("WildcardKey() = %q; want %q", got, want)
	}
}

func TestWildcardKeyDiffersFromAddressKey(t *testing.T) {
	specific := Key("tcp", "192.168.1.1", 22)
	wild := WildcardKey("tcp", 22)
	if specific == wild {
		t.Fatal("specific address key should not equal wildcard key")
	}
}
