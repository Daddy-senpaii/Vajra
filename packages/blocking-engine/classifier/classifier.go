package classifier

// Category represents a classification category for domains.
type Category string

const (
	// CategoryNone represents domains with no restricted content.
	CategoryNone Category = "NONE"
	// CategoryAdult represents domains classified as adult content.
	CategoryAdult Category = "ADULT"
)

// Matcher defines the interface required to check whether a domain belongs to a category.
type Matcher interface {
	Matches(normalizedDomain string) bool
}

// Classifier determines whether a given normalized domain belongs to a restricted category.
// It decouples the classification logic from the underlying storage mechanism.
type Classifier struct {
	matcher Matcher
}

// New creates a new Classifier with the given Matcher.
func New(matcher Matcher) *Classifier {
	return &Classifier{
		matcher: matcher,
	}
}

// Classify evaluates a normalized domain and returns its Category.
func (c *Classifier) Classify(normalizedDomain string) Category {
	if c.matcher != nil && c.matcher.Matches(normalizedDomain) {
		return CategoryAdult
	}
	return CategoryNone
}

// IsAdult returns true if the normalized domain is classified as CategoryAdult.
func (c *Classifier) IsAdult(normalizedDomain string) bool {
	return c.Classify(normalizedDomain) == CategoryAdult
}
