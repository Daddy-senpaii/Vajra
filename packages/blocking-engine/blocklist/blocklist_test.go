package blocklist_test

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/blocklist"
)

func TestBlocklist_AddAndContains(t *testing.T) {
	bl := blocklist.New()

	if bl.Size() != 0 {
		t.Errorf("expected empty blocklist size 0, got %d", bl.Size())
	}

	if err := bl.Add("Adult-Example.com"); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if !bl.Contains("adult-example.com") {
		t.Errorf("expected Contains('adult-example.com') to be true")
	}

	if bl.Contains("other.com") {
		t.Errorf("expected Contains('other.com') to be false")
	}

	if bl.Size() != 1 {
		t.Errorf("expected size 1, got %d", bl.Size())
	}

	// Duplicate add
	if err := bl.Add("adult-example.com"); err != nil {
		t.Fatalf("duplicate add failed: %v", err)
	}
	if bl.Size() != 1 {
		t.Errorf("expected size 1 after duplicate, got %d", bl.Size())
	}
}

func TestBlocklist_SubdomainMatching(t *testing.T) {
	bl := blocklist.New()
	_ = bl.Add("adult-example.com")

	tests := []struct {
		domain   string
		expected bool
	}{
		{"adult-example.com", true},
		{"www.adult-example.com", true},
		{"sub.deep.adult-example.com", true},
		// Similar but unrelated domains must NOT match
		{"notadult-example.com", false},
		{"adult-example.org", false},
		{"fakeadult-example.com", false},
		{"example.com", false},
		{"com", false},
	}

	for _, tt := range tests {
		got := bl.Matches(tt.domain)
		if got != tt.expected {
			t.Errorf("Matches(%q) = %v, want %v", tt.domain, got, tt.expected)
		}
	}
}

func TestBlocklist_Remove(t *testing.T) {
	bl := blocklist.New()
	_ = bl.Add("bad1.com")
	_ = bl.Add("bad2.com")

	if err := bl.Remove("bad1.com"); err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	if bl.Contains("bad1.com") {
		t.Errorf("expected bad1.com to be removed")
	}
	if !bl.Contains("bad2.com") {
		t.Errorf("expected bad2.com to remain")
	}
	if bl.Size() != 1 {
		t.Errorf("expected size 1, got %d", bl.Size())
	}
}

func TestBlocklist_Loader(t *testing.T) {
	bl := blocklist.New()
	content := `
# Adult Content Blocklist
adult-site1.com
adult-site2.com # with comment
duplicate.com
duplicate.com

# Malformed entry to safely skip
-bad..domain
`
	added, err := blocklist.LoadFromReader(strings.NewReader(content), bl)
	if err != nil {
		t.Fatalf("LoadFromReader failed: %v", err)
	}

	if added != 3 {
		t.Errorf("expected 3 added domains, got %d", added)
	}

	if !bl.Contains("adult-site1.com") || !bl.Contains("adult-site2.com") || !bl.Contains("duplicate.com") {
		t.Errorf("expected loaded domains to be present")
	}

	// Test LoadFromFile
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "blocklist.txt")
	if err := os.WriteFile(filePath, []byte("adult-file.com\n# comment\n"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	fileAdded, err := blocklist.LoadFromFile(filePath, bl)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}
	if fileAdded != 1 {
		t.Errorf("expected 1 domain added from file, got %d", fileAdded)
	}
	if !bl.Contains("adult-file.com") {
		t.Errorf("expected adult-file.com to be in blocklist")
	}
}

func TestBlocklist_LoaderCountsOnlyNewEntries(t *testing.T) {
	bl := blocklist.New()
	if err := bl.Add("existing.com"); err != nil {
		t.Fatal(err)
	}

	added, err := blocklist.LoadFromReader(strings.NewReader("existing.com\nnew.com\nnew.com\n"), bl)
	if err != nil {
		t.Fatalf("LoadFromReader failed: %v", err)
	}
	if added != 1 {
		t.Fatalf("LoadFromReader added %d domains, want 1", added)
	}
}

func TestBlocklist_Concurrency(t *testing.T) {
	bl := blocklist.New()
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = bl.Add("concurrent-site.com")
			_ = bl.Matches("sub.concurrent-site.com")
			_ = bl.Contains("concurrent-site.com")
		}(i)
	}
	wg.Wait()

	if bl.Size() != 1 {
		t.Errorf("expected size 1, got %d", bl.Size())
	}
}
