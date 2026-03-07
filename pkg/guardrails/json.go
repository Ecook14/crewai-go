package guardrails

import (
	"encoding/json"
	"fmt"
)

// JSONValidGuardrail ensures the output is a valid JSON string.
type JSONValidGuardrail struct{}

func (g *JSONValidGuardrail) Validate(output string) error {
	var js json.RawMessage
	if err := json.Unmarshal([]byte(output), &js); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func (g *JSONValidGuardrail) Name() string {
	return "JSONValidator"
}

// NewJSONValidator returns a new JSON validator guardrail.
func NewJSONValidator() *JSONValidGuardrail {
	return &JSONValidGuardrail{}
}
