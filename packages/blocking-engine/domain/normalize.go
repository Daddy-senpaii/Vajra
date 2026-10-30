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

	// If a scheme is present (e.g. "http://", "https://"), parse as URL to extract authority.
	if strings.Contains(s, "://") {
		u, err := url.Parse(s)
		if err != nil || u.Host == "" {
			return "", ErrInvalidDomain
		}
		s = u.Host
	} else if strings.Contains(s, "/") {
		// Contains path without scheme (e.g. "example.com/path")
		parts := strings.SplitN(s, "/", 2)
		s = parts[0]
	}

	// Strip a valid numeric port, but reject malformed authorities instead of
	// silently turning inputs such as "example.com:abc" into valid domains.
	var err error
	s, err = stripPort(s)
	if err != nil {
		return "", ErrInvalidDomain
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

func stripPort(authority string) (string, error) {
	if strings.HasPrefix(authority, "[") {
		end := strings.IndexByte(authority, ']')
		if end < 0 {
			return "", ErrInvalidDomain
		}
		host, suffix := authority[1:end], authority[end+1:]
		if suffix != "" {
			if !strings.HasPrefix(suffix, ":") || validatePort(suffix[1:]) != nil {
				return "", ErrInvalidDomain
			}
		}
		return host, nil
	}

	if strings.Count(authority, ":") > 1 {
		return "", ErrInvalidDomain
	}
	if i := strings.LastIndexByte(authority, ':'); i >= 0 {
		if err := validatePort(authority[i+1:]); err != nil {
			return "", err
		}
		return authority[:i], nil
	}
	return authority, nil
}

func validatePort(port string) error {
	if port == "" {
		return ErrInvalidDomain
	}
	value := 0
	for _, digit := range port {
		if digit < '0' || digit > '9' {
			return ErrInvalidDomain
		}
		value = value*10 + int(digit-'0')
		if value > 65535 {
			return ErrInvalidDomain
		}
	}
	if value == 0 {
		return ErrInvalidDomain
	}
	return nil
}
