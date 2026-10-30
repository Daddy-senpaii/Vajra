package engine

import (
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/allowlist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/blocklist"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/classifier"
	"github.com/Daddy-senpaii/Vajra/packages/blocking-engine/domain"
)

// Decision represents the final blocking evaluation outcome.
type Decision string

const (
	// DecisionAllow permits navigation to the domain.
	DecisionAllow Decision = "ALLOW"
	// DecisionBlock prevents navigation to the domain.
	DecisionBlock Decision = "BLOCK"
	// DecisionInvalid indicates the domain input was malformed or failed validation.
	DecisionInvalid Decision = "INVALID"
)

// Result contains the evaluation decision and accompanying metadata.
type Result struct {
	Decision         Decision `json:"decision"`
	NormalizedDomain string   `json:"normalized_domain,omitempty"`
	Reason           string   `json:"reason"`
	Error            error    `json:"-"`
}

// Engine coordinates domain normalization, validation, allowlist check,
// and blocklist classification into a single unified decision pipeline.
type Engine struct {
	allowlist  *allowlist.Allowlist
	blocklist  *blocklist.Blocklist
	classifier *classifier.Classifier
}

// New creates an Engine initialized with the given allowlist and blocklist.
// If allowlist or blocklist is nil, a new empty one is instantiated.
func New(al *allowlist.Allowlist, bl *blocklist.Blocklist) *Engine {
	if al == nil {
		al = allowlist.New()
	}
	if bl == nil {
		bl = blocklist.New()
	}
	cls := classifier.New(bl)

	return &Engine{
		allowlist:  al,
		blocklist:  bl,
		classifier: cls,
	}
}

// NewDefault creates an Engine preloaded with the bundled offline adult-domain
// blocklist. It returns an error if the embedded list cannot be loaded.
func NewDefault() (*Engine, error) {
	bl, _, err := blocklist.NewDefault()
	if err != nil {
		return nil, err
	}
	return New(nil, bl), nil
}

// Evaluate processes an input domain according to the specified decision pipeline:
// 1. Normalize and validate the domain input.
// 2. Check allowlist (takes precedence).
// 3. Check blocklist classifier.
// 4. Return decision (ALLOW, BLOCK, or INVALID).
func (e *Engine) Evaluate(rawDomain string) Result {
	norm, err := domain.Normalize(rawDomain)
	if err != nil {
		return Result{
			Decision: DecisionInvalid,
			Reason:   "domain normalization or validation failed: " + err.Error(),
			Error:    err,
		}
	}

	// 1. Check Allowlist (strict precedence over Blocklist)
	if e.allowlist.Matches(norm) {
		return Result{
			Decision:         DecisionAllow,
			NormalizedDomain: norm,
			Reason:           "domain matched allowlist",
		}
	}

	// 2. Check Blocklist / Classifier
	if e.classifier.IsAdult(norm) {
		return Result{
			Decision:         DecisionBlock,
			NormalizedDomain: norm,
			Reason:           "domain matched adult blocklist",
		}
	}

	// 3. Default to Allow
	return Result{
		Decision:         DecisionAllow,
		NormalizedDomain: norm,
		Reason:           "domain not in blocklist",
	}
}

// Allowlist returns the underlying Allowlist instance.
func (e *Engine) Allowlist() *allowlist.Allowlist {
	return e.allowlist
}

// Blocklist returns the underlying Blocklist instance.
func (e *Engine) Blocklist() *blocklist.Blocklist {
	return e.blocklist
}

// Classifier returns the underlying Classifier instance.
func (e *Engine) Classifier() *classifier.Classifier {
	return e.classifier
}
