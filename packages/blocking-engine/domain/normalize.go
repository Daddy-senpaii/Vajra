package domain

import (
	"errors"
	"net/url"
	"strings"

	"golang.org/x/net/idna"
)

var (
	// ErrEmptyDomain indicates that the input domain string is empty.
	ErrEmptyDomain = errors.New("domain is empty")
	// ErrInvalidDomain indicates that the domain failed syntax validation.
	ErrInvalidDomain = errors.New("domain is invalid")
)

// Normalize cleans, sanitizes, converts IDNs to ASCII Punycode,
// lowercases, and validates an input domain or URL string.
func Normalize(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", ErrEmptyDomain
	}

	// If a scheme is present (e.g. "http://", "https://"), parse as URL to extract host.
	if strings.Contains(s, "://") {
		u, err := url.Parse(s)
		if err != nil {
			return "", ErrInvalidDomain
		}
		s = u.Host
	} else if strings.Contains(s, "/") {
		// Contains path without scheme (e.g. "example.com/path")
		parts := strings.SplitN(s, "/", 2)
		s = parts[0]
	}

	// Strip port if present (e.g. "example.com:8080" or "[::1]:8080")
	if strings.HasPrefix(s, "[") {
		// IPv6 literal - strip closing bracket and port
		closeBracket := strings.Index(s, "]")
		if closeBracket != -1 {
			s = s[1:closeBracket]
		}
	} else if colonIdx := strings.LastIndex(s, ":"); colonIdx != -1 {
		// Only strip port if there is only one colon (not IPv6 without brackets)
		if strings.Count(s, ":") == 1 {
			s = s[:colonIdx]
		}
	}

	// Strip trailing FQDN dot (e.g. "example.com." -> "example.com")
	s = strings.TrimSuffix(s, ".")

	s = strings.TrimSpace(s)
	if s == "" {
		return "", ErrEmptyDomain
	}

	// Convert Internationalized Domain Names (IDN) to ASCII punycode and lowercase
	asciiDomain, err := idna.Lookup.ToASCII(s)
	if err != nil {
		return "", ErrInvalidDomain
	}

	// Ensure lowercase
	asciiDomain = strings.ToLower(asciiDomain)

	// Validate the normalized domain syntax
	if err := Validate(asciiDomain); err != nil {
		return "", err
	}

	return asciiDomain, nil
}
