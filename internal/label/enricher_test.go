package label_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/label"
	"github.com/user/portwatch/internal/scanner"
)

func hasLabelTag(p scanner.Port, substr string) bool {
	for _, t := range p.Tags {
		if strings.Contains(t, substr) {
			return true
		}
	}
	return false
}

func TestEnrichKnownPortAddsTag(t *testing.T) {
	e := label.NewEnricher(label.New(nil))
	p := makePort("tcp", 22)
	got := e.Enrich(p)
	if !hasLabelTag(got, "ssh") {
		t.Fatalf("expected label:ssh tag, got %v", got.Tags)
	}
}

func TestEnrichUnknownPortNoTag(t *testing.T) {
	e := label.NewEnricher(label.New(nil))
	p := makePort("tcp", 9999)
	got := e.Enrich(p)
	if len(got.Tags) != 0 {
		t.Fatalf("expected no tags, got %v", got.Tags)
	}
}

func TestEnrichDoesNotDuplicateTag(t *testing.T) {
	e := label.NewEnricher(label.New(nil))
	p := makePort("tcp", 80)
	p.Tags = []string{"label:http"}
	got := e.Enrich(p)
	count := 0
	for _, t := range got.Tags {
		if t == "label:http" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly one label:http tag, got %d", count)
	}
}

func TestEnrichAllAppliesLabels(t *testing.T) {
	e := label.NewEnricher(label.New(nil))
	ports := []scanner.Port{
		makePort("tcp", 80),
		makePort("tcp", 443),
		makePort("tcp", 9999),
	}
	result := e.EnrichAll(ports)
	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
	if !hasLabelTag(result[0], "http") {
		t.Error("port 80 should have label:http")
	}
	if !hasLabelTag(result[1], "https") {
		t.Error("port 443 should have label:https")
	}
	if hasLabelTag(result[2], "label:") {
		t.Error("port 9999 should have no label tag")
	}
}

func TestEnrichWithOverrideRule(t *testing.T) {
	rules := []label.Rule{
		{Port: 2222, Protocol: "tcp", Label: "DevSSH"},
	}
	e := label.NewEnricher(label.New(rules))
	p := makePort("tcp", 2222)
	got := e.Enrich(p)
	if !hasLabelTag(got, "devssh") {
		t.Fatalf("expected label:devssh, got %v", got.Tags)
	}
}
