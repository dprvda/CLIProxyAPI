package registry

import (
	"encoding/json"
	"testing"
)

// The embedded catalog carries fork-local model registrations before the remote
// catalog publishes them. A remote refresh must not drop those.
func TestOverlayEmbeddedOnlyModelsKeepsForkLocalEntries(t *testing.T) {
	// Simulate a remote catalog that lacks claude-fable-5-1 (as of 2026-09-01 the
	// published remote really does) but updates metadata for claude-fable-5.
	remote := &staticModelsJSON{
		Claude: []*ModelInfo{
			{ID: "claude-fable-5", Object: "model", Type: "claude", DisplayName: "Claude Fable 5 (remote)"},
		},
	}

	overlayEmbeddedOnlyModels(remote)

	byID := make(map[string]*ModelInfo, len(remote.Claude))
	for _, m := range remote.Claude {
		byID[m.ID] = m
	}
	if _, ok := byID["claude-fable-5-1"]; !ok {
		t.Fatalf("claude-fable-5-1 missing after overlay; got %d claude models", len(remote.Claude))
	}
	// Remote entries win for IDs both sides define.
	if got := byID["claude-fable-5"].DisplayName; got != "Claude Fable 5 (remote)" {
		t.Fatalf("remote entry was replaced by embedded one: display_name=%q", got)
	}
}

// The embedded catalog itself must define claude-fable-5-1 with the same shape
// the rest of the claude family uses.
func TestEmbeddedCatalogHasClaudeFable51(t *testing.T) {
	var embedded staticModelsJSON
	if err := json.Unmarshal(embeddedModelsJSON, &embedded); err != nil {
		t.Fatalf("embedded models.json does not parse: %v", err)
	}
	for _, m := range embedded.Claude {
		if m.ID != "claude-fable-5-1" {
			continue
		}
		if m.Type != "claude" || m.OwnedBy != "anthropic" {
			t.Fatalf("claude-fable-5-1 misregistered: type=%q owned_by=%q", m.Type, m.OwnedBy)
		}
		if m.ContextLength != 1000000 || m.MaxCompletionTokens != 128000 {
			t.Fatalf("claude-fable-5-1 window wrong: context=%d max_completion=%d", m.ContextLength, m.MaxCompletionTokens)
		}
		if m.Thinking == nil || len(m.Thinking.Levels) != 5 {
			t.Fatalf("claude-fable-5-1 thinking levels wrong: %+v", m.Thinking)
		}
		return
	}
	t.Fatal("claude-fable-5-1 not found in embedded models.json claude section")
}
