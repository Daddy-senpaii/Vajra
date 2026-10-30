package classifier_test

import (
	"testing"

	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/classifier"
)

type mockMatcher struct {
	blocked map[string]bool
}

func (m *mockMatcher) Matches(domain string) bool {
	return m.blocked[domain]
}

func TestClassifier_Classify(t *testing.T) {
	mock := &mockMatcher{
		blocked: map[string]bool{
			"adult-example.com": true,
		},
	}

	c := classifier.New(mock)

	tests := []struct {
		domain       string
		wantCategory classifier.Category
		wantAdult    bool
	}{
		{"adult-example.com", classifier.CategoryAdult, true},
		{"example.com", classifier.CategoryNone, false},
		{"notadult-example.com", classifier.CategoryNone, false},
	}

	for _, tt := range tests {
		cat := c.Classify(tt.domain)
		if cat != tt.wantCategory {
			t.Errorf("Classify(%q) = %v, want %v", tt.domain, cat, tt.wantCategory)
		}
		if got := c.IsAdult(tt.domain); got != tt.wantAdult {
			t.Errorf("IsAdult(%q) = %v, want %v", tt.domain, got, tt.wantAdult)
		}
	}
}

func TestClassifier_NilMatcher(t *testing.T) {
	c := classifier.New(nil)
	if c.IsAdult("any.com") {
		t.Errorf("expected false with nil matcher")
	}
	if c.Classify("any.com") != classifier.CategoryNone {
		t.Errorf("expected CategoryNone with nil matcher")
	}
}
