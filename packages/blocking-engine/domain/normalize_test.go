package domain_test

import (
	"strings"
	"testing"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
)

func TestNormalize_Valid(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "exact lowercase",
			input:    "example.com",
			expected: "example.com",
		},
		{
			name:     "uppercase and mixed case",
			input:    "WWW.Example.COM",
			expected: "www.example.com",
		},
		{
			name:     "leading and trailing whitespace",
			input:    "   example.com   \t\n",
			expected: "example.com",
		},
		{
			name:     "trailing FQDN dot",
			input:    "example.com.",
			expected: "example.com",
		},
		{
			name:     "subdomain hierarchy",
			input:    "deep.sub.example.com",
			expected: "deep.sub.example.com",
		},
		{
			name:     "http url with path",
			input:    "http://example.com/path/to/resource",
			expected: "example.com",
		},
		{
			name:     "https url with port and query",
			input:    "https://example.com:8443/search?q=test#frag",
			expected: "example.com",
		},
		{
			name:     "url without scheme but with path",
			input:    "example.com/page",
			expected: "example.com",
		},
		{
			name:     "url with userinfo",
			input:    "https://user:pass@example.com/home",
			expected: "example.com",
		},
		{
			name:     "domain with hyphens in label",
			input:    "adult-example.com",
			expected: "adult-example.com",
		},
		{
			name:     "internationalized domain name (german umlaut)",
			input:    "bücher.example",
			expected: "xn--bcher-kva.example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := domain.Normalize(tt.input)
			if err != nil {
				t.Fatalf("Normalize(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestNormalize_Invalid(t *testing.T) {
	longLabel := strings.Repeat("a", 64) + ".com"
	longDomain := strings.Repeat("a.", 130) + "com" // > 253 chars

	invalidInputs := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"whitespace only", "   \t\n"},
		{"single word / no dot", "invalid"},
		{"string with spaces", "invalid domain"},
		{"label starts with hyphen", "-example.com"},
		{"label ends with hyphen", "example-.com"},
		{"consecutive dots", "example..com"},
		{"leading dot", ".example.com"},
		{"special character @", "exam@ple.com"},
		{"underscore in hostname", "exam_ple.com"},
		{"exclamation mark", "adult!site.com"},
		{"label exceeds 63 chars", longLabel},
		{"domain exceeds 253 chars", longDomain},
		{"numeric TLD", "example.123"},
		{"single char TLD", "example.c"},
	}

	for _, tt := range invalidInputs {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.Normalize(tt.input)
			if err == nil {
				t.Errorf("Normalize(%q) expected error, got nil", tt.input)
			}
		})
	}
}
