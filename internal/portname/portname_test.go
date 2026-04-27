package portname

import "testing"

func TestResolveBuiltinTCP(t *testing.T) {
	r := New(nil)
	if got := r.Resolve("tcp", 22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestResolveBuiltinUDP(t *testing.T) {
	r := New(nil)
	if got := r.Resolve("udp", 53); got != "dns" {
		t.Fatalf("expected dns, got %q", got)
	}
}

func TestResolveUnknownPortReturnsEmpty(t *testing.T) {
	r := New(nil)
	if got := r.Resolve("tcp", 9999); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	r := New(map[string]string{"tcp/80": "myproxy"})
	if got := r.Resolve("tcp", 80); got != "myproxy" {
		t.Fatalf("expected myproxy, got %q", got)
	}
}

func TestOverrideForUnknownPort(t *testing.T) {
	r := New(map[string]string{"tcp/19999": "internal-api"})
	if got := r.Resolve("tcp", 19999); got != "internal-api" {
		t.Fatalf("expected internal-api, got %q", got)
	}
}

func TestResolveOrDefaultReturnsFallback(t *testing.T) {
	r := New(nil)
	got := r.ResolveOrDefault("tcp", 9999, "unknown")
	if got != "unknown" {
		t.Fatalf("expected unknown, got %q", got)
	}
}

func TestResolveOrDefaultReturnsMappedName(t *testing.T) {
	r := New(nil)
	got := r.ResolveOrDefault("tcp", 443, "unknown")
	if got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestProtocolDistinctionInBuiltin(t *testing.T) {
	r := New(nil)
	// tcp/53 and udp/53 both resolve to dns
	if r.Resolve("tcp", 53) != "dns" {
		t.Fatal("tcp/53 should resolve to dns")
	}
	if r.Resolve("udp", 53) != "dns" {
		t.Fatal("udp/53 should resolve to dns")
	}
	// tcp/21 has no udp counterpart
	if r.Resolve("udp", 21) != "" {
		t.Fatal("udp/21 should return empty")
	}
}
