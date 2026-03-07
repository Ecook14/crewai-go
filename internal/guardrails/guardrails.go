// Package guardrails provides input/output validation middleware for agent execution.
// Guardrails intercept agent outputs to enforce quality, safety, and structural constraints.
package guardrails

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Guardrail defines the interface for all validation middleware.
// Implementations validate agent output and return an error if the output fails validation.
type Guardrail interface {
	// Name returns a human-readable identifier for this guardrail.
	Name() string
	// Validate checks the output string and returns nil if valid, or an error describing the failure.
	Validate(output string) error
}

// MaxTokenGuardrail rejects outputs exceeding a specified word count.
// Uses word-level approximation (1 word ≈ 1 token) for fast, dependency-free checking.
type MaxTokenGuardrail struct {
	MaxTokens int
}

func NewMaxTokenGuardrail(maxTokens int) *MaxTokenGuardrail {
	return &MaxTokenGuardrail{MaxTokens: maxTokens}
}

func (g *MaxTokenGuardrail) Name() string { return "MaxTokenGuardrail" }

func (g *MaxTokenGuardrail) Validate(output string) error {
	words := strings.Fields(output)
	if len(words) > g.MaxTokens {
		return fmt.Errorf("output has %d tokens, exceeds maximum of %d", len(words), g.MaxTokens)
	}
	return nil
}

// ContentFilterGuardrail blocks outputs containing any of the forbidden patterns.
// Patterns are compiled as regular expressions for flexible matching.
type ContentFilterGuardrail struct {
	ForbiddenPatterns []*regexp.Regexp
	RawPatterns       []string // stored for error messages
}

func NewContentFilterGuardrail(patterns []string) (*ContentFilterGuardrail, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern '%s': %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &ContentFilterGuardrail{
		ForbiddenPatterns: compiled,
		RawPatterns:       patterns,
	}, nil
}

func (g *ContentFilterGuardrail) Name() string { return "ContentFilterGuardrail" }

func (g *ContentFilterGuardrail) Validate(output string) error {
	for i, re := range g.ForbiddenPatterns {
		if re.MatchString(output) {
			return fmt.Errorf("output matches forbidden pattern '%s'", g.RawPatterns[i])
		}
	}
	return nil
}

// SchemaGuardrail validates that the output is valid JSON and can be unmarshalled
// into the provided schema template. The schema is a pointer to a Go struct.
type SchemaGuardrail struct {
	SchemaTemplate interface{}
}

func NewSchemaGuardrail(schema interface{}) *SchemaGuardrail {
	return &SchemaGuardrail{SchemaTemplate: schema}
}

func (g *SchemaGuardrail) Name() string { return "SchemaGuardrail" }

func (g *SchemaGuardrail) Validate(output string) error {
	trimmed := strings.TrimSpace(output)
	if !json.Valid([]byte(trimmed)) {
		return fmt.Errorf("output is not valid JSON")
	}
	// Attempt unmarshal to verify structural compatibility
	if err := json.Unmarshal([]byte(trimmed), g.SchemaTemplate); err != nil {
		return fmt.Errorf("output does not match expected schema: %w", err)
	}
	return nil
}

// RunAll executes all guardrails against the given output.
// Returns nil if all pass, or the first encountered error.
func RunAll(guardrails []Guardrail, output string) error {
	for _, g := range guardrails {
		if err := g.Validate(output); err != nil {
			return fmt.Errorf("guardrail '%s' failed: %w", g.Name(), err)
		}
	}
	return nil
}
