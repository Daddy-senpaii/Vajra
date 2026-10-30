package allowlist_test

import (
	"strings"
	"sync"
	"testing"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/allowlist"
)

func TestAllowlist_AddAndContains(t *testing.T) {
	al := allowlist.New()

	if err := al.Add("  Example.COM  "); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if !al.Contains("example.com") {
		t.Errorf("expected Contains('example.com') to be true")
	}

	if al.Contains("other.com") {
		t.Errorf("expected Contains('other.com') to be false")
	}

	if al.Size() != 1 {
		t.Errorf("expected size 1, got %d", al.Size())
	}

	// Duplicate addition
	if err := al.Add("example.com"); err != nil {
		t.Fatalf("duplicate Add failed: %v", err)
	}
	if al.Size() != 1 {
		t.Errorf("expected size 1 after duplicate add, got %d", al.Size())
	}

	// Invalid addition
	if err := al.Add("invalid domain"); err == nil {
		t.Errorf("expected error adding invalid domain, got nil")
	}
}

func TestAllowlist_SubdomainMatching(t *testing.T) {
	al := allowlist.New()
	if err := al.Add("example.com"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	tests := []struct {
		query    string
		expected bool
	}{
		{"example.com", true},
		{"www.example.com", true},
		{"deep.sub.example.com", true},
		{"notexample.com", false},
		{"example.org", false},
		{"com", false},
	}

	for _, tt := range tests {
		got := al.Matches(tt.query)
		if got != tt.expected {
			t.Errorf("Matches(%q) = %v, want %v", tt.query, got, tt.expected)
		}
	}
}

func TestAllowlist_Remove(t *testing.T) {
	al := allowlist.New()
	_ = al.Add("example.com")
	_ = al.Add("test.org")

	if err := al.Remove("EXAMPLE.COM"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if al.Contains("example.com") {
		t.Errorf("expected example.com to be removed")
	}
	if !al.Contains("test.org") {
		t.Errorf("expected test.org to still be present")
	}
	if al.Size() != 1 {
		t.Errorf("expected size 1, got %d", al.Size())
	}
}

func TestAllowlist_Load(t *testing.T) {
	al := allowlist.New()
	content := `
# Sample Allowlist
example.com
www.mysite.org # with comment
duplicate.com
duplicate.com

# Malformed line to skip safely
invalid domain
`
	added, err := al.Load(strings.NewReader(content))
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if added != 3 {
		t.Errorf("expected 3 added domains, got %d", added)
	}

	if !al.Contains("example.com") || !al.Contains("www.mysite.org") || !al.Contains("duplicate.com") {
		t.Errorf("expected loaded domains to be present")
	}
}

func TestAllowlist_Concurrency(t *testing.T) {
	al := allowlist.New()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = al.Add("example.com")
			_ = al.Matches("www.example.com")
			_ = al.Contains("example.com")
		}(i)
	}
	wg.Wait()

	if al.Size() != 1 {
		t.Errorf("expected size 1 after concurrent operations, got %d", al.Size())
	}
}
