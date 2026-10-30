package allowlist

import (
	"bufio"
	"io"
	"strings"
	"sync"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
)

// Allowlist manages explicitly permitted domains in a thread-safe manner.
type Allowlist struct {
	mu      sync.RWMutex
	domains map[string]struct{}
}

// New creates an empty, initialized Allowlist.
func New() *Allowlist {
	return &Allowlist{
		domains: make(map[string]struct{}),
	}
}

// Add normalizes and adds a domain to the allowlist.
// Returns an error if the domain is invalid.
// Prevents duplicates automatically.
func (a *Allowlist) Add(rawDomain string) error {
	norm, err := domain.Normalize(rawDomain)
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.domains[norm] = struct{}{}
	return nil
}

// Remove normalizes and removes a domain from the allowlist.
// If the domain is not in the list, it is a no-op and returns nil.
func (a *Allowlist) Remove(rawDomain string) error {
	norm, err := domain.Normalize(rawDomain)
	if err != nil {
		return err
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.domains, norm)
	return nil
}

// Contains checks if the exact normalized domain exists in the allowlist.
func (a *Allowlist) Contains(normalizedDomain string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, exists := a.domains[normalizedDomain]
	return exists
}

// Matches checks if the domain or any of its parent domains exist in the allowlist.
// For example, if "example.com" is allowlisted, "sub.example.com" and "a.b.example.com"
// will return true.
func (a *Allowlist) Matches(normalizedDomain string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Check exact match first
	if _, exists := a.domains[normalizedDomain]; exists {
		return true
	}

	// Check parent domains
	current := normalizedDomain
	for {
		dotIdx := strings.Index(current, ".")
		if dotIdx == -1 {
			break
		}
		current = current[dotIdx+1:]
		// Do not match single-label TLDs (e.g. "com") unless they have at least one dot
		if !strings.Contains(current, ".") {
			break
		}
		if _, exists := a.domains[current]; exists {
			return true
		}
	}

	return false
}

// Size returns the count of unique domains in the allowlist.
func (a *Allowlist) Size() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.domains)
}

// Clear removes all domains from the allowlist.
func (a *Allowlist) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.domains = make(map[string]struct{})
}

// Load reads lines from an io.Reader (ignoring blank lines and # comments),
// normalizes valid domains, and adds them to the allowlist.
// Returns the count of newly added domains.
func (a *Allowlist) Load(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	added := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Strip inline comments if present
		if commentIdx := strings.Index(line, "#"); commentIdx != -1 {
			line = strings.TrimSpace(line[:commentIdx])
			if line == "" {
				continue
			}
		}

		norm, err := domain.Normalize(line)
		if err != nil {
			// Skip malformed entries in list file safely
			continue
		}

		a.mu.Lock()
		if _, exists := a.domains[norm]; !exists {
			a.domains[norm] = struct{}{}
			added++
		}
		a.mu.Unlock()
	}

	if err := scanner.Err(); err != nil {
		return added, err
	}
	return added, nil
}
