package domain

import (
	"errors"
	"strings"
	"unicode"
)

// Validate checks whether a normalized domain string conforms to RFC 1035/1123
// and standard domain name syntax requirements.
func Validate(domain string) error {
	if domain == "" {
		return ErrEmptyDomain
	}

	// RFC 1035: Total domain length cannot exceed 253 characters.
	if len(domain) > 253 {
		return errors.New("domain exceeds maximum length of 253 characters")
	}

	// A valid internet domain must have at least one dot (e.g. "example.com")
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return errors.New("domain must contain at least two labels separated by a dot")
	}

	for i, label := range labels {
		// RFC 1035: Label length must be between 1 and 63 characters
		if len(label) == 0 {
			return errors.New("domain label cannot be empty (consecutive dots)")
		}
		if len(label) > 63 {
			return errors.New("domain label exceeds maximum length of 63 characters")
		}

		// Labels must not start or end with a hyphen
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return errors.New("domain label cannot start or end with a hyphen")
		}

		// Check valid characters: only lowercase alphanumeric and hyphen allowed
		for _, ch := range label {
			if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-') {
				return errors.New("domain contains invalid characters")
			}
		}

		// Top-level domain (last label) checks
		if i == len(labels)-1 {
			// TLD must be at least 2 characters
			if len(label) < 2 {
				return errors.New("top-level domain must be at least 2 characters")
			}

			// TLD cannot be all digits
			allDigits := true
			for _, ch := range label {
				if !unicode.IsDigit(ch) {
					allDigits = false
					break
				}
			}
			if allDigits {
				return errors.New("top-level domain cannot be all digits")
			}
		}
	}

	return nil
}
