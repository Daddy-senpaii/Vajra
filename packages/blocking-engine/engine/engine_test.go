package engine_test

import (
	"testing"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/allowlist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/blocklist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/engine"
)

func setupTestEngine(t *testing.T) *engine.Engine {
	t.Helper()
	al := allowlist.New()
	bl := blocklist.New()

	// Seed blocklist
	if err := bl.Add("adult-example.com"); err != nil {
		t.Fatalf("failed to add to blocklist: %v", err)
	}
	if err := bl.Add("restricted-hub.net"); err != nil {
		t.Fatalf("failed to add to blocklist: %v", err)
	}

	// Seed allowlist
	if err := al.Add("allowed-site.com"); err != nil {
		t.Fatalf("failed to add to allowlist: %v", err)
	}
	// Subdomain override on an otherwise blocked root
	if err := al.Add("safe.adult-example.com"); err != nil {
		t.Fatalf("failed to add to allowlist: %v", err)
	}

	return engine.New(al, bl)
}

func TestEngine_Section6_IllustrativeCases(t *testing.T) {
	eng := setupTestEngine(t)

	tests := []struct {
		name         string
		input        string
		wantDecision engine.Decision
	}{
		{
			name:         "exact adult domain is blocked",
			input:        "adult-example.com",
			wantDecision: engine.DecisionBlock,
		},
		{
			name:         "subdomain of adult domain is blocked",
			input:        "www.adult-example.com",
			wantDecision: engine.DecisionBlock,
		},
		{
			name:         "unlisted regular domain is allowed",
			input:        "example.com",
			wantDecision: engine.DecisionAllow,
		},
		{
			name:         "similar but unrelated domain is allowed",
			input:        "notadult-example.com",
			wantDecision: engine.DecisionAllow,
		},
		{
			name:         "invalid domain string returns INVALID",
			input:        "Invalid domain",
			wantDecision: engine.DecisionInvalid,
		},
		{
			name:         "explicitly allowed domain is allowed",
			input:        "allowed-site.com",
			wantDecision: engine.DecisionAllow,
		},
		{
			name:         "allowlist precedence over blocklist (subdomain override)",
			input:        "safe.adult-example.com",
			wantDecision: engine.DecisionAllow,
		},
		{
			name:         "subdomain of explicitly allowed subdomain is allowed",
			input:        "api.safe.adult-example.com",
			wantDecision: engine.DecisionAllow,
		},
		{
			name:         "other subdomain of blocked domain remains blocked",
			input:        "unsafe.adult-example.com",
			wantDecision: engine.DecisionBlock,
		},
		{
			name:         "case-insensitive matching (uppercase blocked domain)",
			input:        "WWW.ADULT-EXAMPLE.COM",
			wantDecision: engine.DecisionBlock,
		},
		{
			name:         "casing and whitespace normalization",
			input:        "   Adult-Example.Com   \n",
			wantDecision: engine.DecisionBlock,
		},
		{
			name:         "URL scheme stripping during evaluation",
			input:        "https://www.adult-example.com/videos/123",
			wantDecision: engine.DecisionBlock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := eng.Evaluate(tt.input)
			if res.Decision != tt.wantDecision {
				t.Errorf("Evaluate(%q) = %v (reason: %q), want %v", tt.input, res.Decision, res.Reason, tt.wantDecision)
			}
		})
	}
}

func TestEngine_EmptyBlocklist(t *testing.T) {
	eng := engine.New(nil, nil)

	res := eng.Evaluate("adult-example.com")
	if res.Decision != engine.DecisionAllow {
		t.Errorf("expected ALLOW on empty engine, got %v", res.Decision)
	}

	resInv := eng.Evaluate("")
	if resInv.Decision != engine.DecisionInvalid {
		t.Errorf("expected INVALID on empty domain, got %v", resInv.Decision)
	}
}

func TestEngine_DynamicUpdates(t *testing.T) {
	eng := engine.New(nil, nil)

	// Initially allowed
	if res := eng.Evaluate("new-bad-site.com"); res.Decision != engine.DecisionAllow {
		t.Fatalf("expected initial ALLOW, got %v", res.Decision)
	}

	// Add to blocklist
	if err := eng.Blocklist().Add("new-bad-site.com"); err != nil {
		t.Fatalf("failed to add: %v", err)
	}
	if res := eng.Evaluate("new-bad-site.com"); res.Decision != engine.DecisionBlock {
		t.Fatalf("expected BLOCK after addition, got %v", res.Decision)
	}

	// Add to allowlist -> overrides blocklist
	if err := eng.Allowlist().Add("new-bad-site.com"); err != nil {
		t.Fatalf("failed to add to allowlist: %v", err)
	}
	if res := eng.Evaluate("new-bad-site.com"); res.Decision != engine.DecisionAllow {
		t.Fatalf("expected ALLOW after allowlist override, got %v", res.Decision)
	}

	// Remove from allowlist -> back to blocked
	if err := eng.Allowlist().Remove("new-bad-site.com"); err != nil {
		t.Fatalf("failed to remove from allowlist: %v", err)
	}
	if res := eng.Evaluate("new-bad-site.com"); res.Decision != engine.DecisionBlock {
		t.Fatalf("expected BLOCK after allowlist removal, got %v", res.Decision)
	}

	// Remove from blocklist -> back to allowed
	if err := eng.Blocklist().Remove("new-bad-site.com"); err != nil {
		t.Fatalf("failed to remove from blocklist: %v", err)
	}
	if res := eng.Evaluate("new-bad-site.com"); res.Decision != engine.DecisionAllow {
		t.Fatalf("expected ALLOW after blocklist removal, got %v", res.Decision)
	}
}
