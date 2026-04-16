package tagger_test

import (
	"testing"

	"portwatch/internal/tagger"
)

func TestBuiltinTCPPort(t *testing.T) {
	tg := tagger.New(nil)
	if got := tg.Tag("tcp", 22); got != "ssh" {
		t.Fatalf("expected ssh, got %s", got)
	}
}

func TestBuiltinUDPPort(t *testing.T) {
	tg := tagger.New(nil)
	if got := tg.Tag("udp", 53); got != "dns" {
		t.Fatalf("expected dns, got %s", got)
	}
}

func TestUnknownPortReturnsGeneric(t *testing.T) {
	tg := tagger.New(nil)
	if got := tg.Tag("tcp", 9999); got != "port/9999" {
		t.Fatalf("unexpected label: %s", got)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	tg := tagger.New(map[string]string{"tcp:80": "my-app"})
	if got := tg.Tag("tcp", 80); got != "my-app" {
		t.Fatalf("expected my-app, got %s", got)
	}
}

func TestOverrideForUnknownPort(t *testing.T) {
	tg := tagger.New(map[string]string{"tcp:7777": "custom"})
	if got := tg.Tag("tcp", 7777); got != "custom" {
		t.Fatalf("expected custom, got %s", got)
	}
}

func TestKnownBuiltin(t *testing.T) {
	tg := tagger.New(nil)
	if !tg.Known("tcp", 443) {
		t.Fatal("443/tcp should be known")
	}
}

func TestKnownOverride(t *testing.T) {
	tg := tagger.New(map[string]string{"tcp:9000": "custom"})
	if !tg.Known("tcp", 9000) {
		t.Fatal("9000/tcp should be known via override")
	}
}

func TestUnknownNotKnown(t *testing.T) {
	tg := tagger.New(nil)
	if tg.Known("tcp", 12345) {
		t.Fatal("12345/tcp should not be known")
	}
}
