package blocklist

import (
	"strings"
	"sync"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
)

// Blocklist manages adult-content domain entries in a thread-safe manner.
type Blocklist struct {
	mu      sync.RWMutex
	domains map[string]struct{}
}

// New creates an empty, initialized Blocklist.
func New() *Blocklist {
	return &Blocklist{
		domains: make(map[string]struct{}),
	}
}

// Add normalizes and adds a domain to the blocklist.
// Duplicate additions are automatically prevented.
func (b *Blocklist) Add(rawDomain string) error {
	norm, err := domain.Normalize(rawDomain)
	if err != nil {
		return err
	}

	b.addNormalized(norm)
	return nil
}

// addNormalized inserts a normalized domain and reports whether it was new.
func (b *Blocklist) addNormalized(normalizedDomain string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.domains[normalizedDomain]; exists {
		return false
	}
	b.domains[normalizedDomain] = struct{}{}
	return true
}

// Remove normalizes and removes a domain from the blocklist.
// If the domain does not exist in the blocklist, it is a no-op and returns nil.
func (b *Blocklist) Remove(rawDomain string) error {
	norm, err := domain.Normalize(rawDomain)
	if err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.domains, norm)
	return nil
}

// Contains checks whether the exact normalized domain exists in the blocklist.
func (b *Blocklist) Contains(normalizedDomain string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, exists := b.domains[normalizedDomain]
	return exists
}

// Matches checks whether the given normalized domain or any of its parent
// domains exist in the blocklist (hierarchical subdomain matching).
// E.g., if "adult-example.com" is blocked, both "adult-example.com" and
// "sub.adult-example.com" match.
func (b *Blocklist) Matches(normalizedDomain string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Check exact match first
	if _, exists := b.domains[normalizedDomain]; exists {
		return true
	}

	// Check parent domain hierarchy
	current := normalizedDomain
	for {
		dotIdx := strings.Index(current, ".")
		if dotIdx == -1 {
			break
		}
		current = current[dotIdx+1:]
		// Stop if reached single-label TLD
		if !strings.Contains(current, ".") {
			break
		}
		if _, exists := b.domains[current]; exists {
			return true
		}
	}

	return false
}

// Size returns the total number of unique domains in the blocklist.
func (b *Blocklist) Size() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.domains)
}

// Clear removes all domains from the blocklist.
func (b *Blocklist) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.domains = make(map[string]struct{})
}
