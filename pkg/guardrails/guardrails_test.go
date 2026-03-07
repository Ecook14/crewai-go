package guardrails

import (
	"testing"
)

func TestMaxTokenGuardrail(t *testing.T) {
	g := NewMaxTokenGuardrail(5)

	if g.Name() != "MaxTokenGuardrail" {
		t.Errorf("expected name 'MaxTokenGuardrail', got '%s'", g.Name())
	}

	// Should pass: 3 words < 5 max tokens
	if err := g.Validate("hello world today"); err != nil {
		t.Errorf("expected pass for 3 words, got error: %v", err)
	}

	// Should fail: 6 words > 5 max tokens
	if err := g.Validate("one two three four five six"); err == nil {
		t.Error("expected failure for 6 words, got nil")
	}

	// Edge case: exactly at limit
	if err := g.Validate("one two three four five"); err != nil {
		t.Errorf("expected pass for exactly 5 words, got error: %v", err)
	}
}

func TestContentFilterGuardrail(t *testing.T) {
	g, err := NewContentFilterGuardrail([]string{`(?i)password`, `(?i)secret`})
	if err != nil {
		t.Fatalf("NewContentFilterGuardrail failed: %v", err)
	}

	if g.Name() != "ContentFilterGuardrail" {
		t.Errorf("expected name 'ContentFilterGuardrail', got '%s'", g.Name())
	}

	// Should pass
	if err := g.Validate("This is a safe output."); err != nil {
		t.Errorf("expected pass, got error: %v", err)
	}

	// Should fail — contains "password"
	if err := g.Validate("The password is 12345."); err == nil {
		t.Error("expected failure for 'password', got nil")
	}

	// Should fail — contains "Secret" (case insensitive)
	if err := g.Validate("This is a Secret value."); err == nil {
		t.Error("expected failure for 'Secret', got nil")
	}
}

func TestContentFilterGuardrailInvalidPattern(t *testing.T) {
	_, err := NewContentFilterGuardrail([]string{`[invalid`})
	if err == nil {
		t.Error("expected error for invalid regex pattern, got nil")
	}
}

func TestSchemaGuardrail(t *testing.T) {
	type TestSchema struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}

	schema := &TestSchema{}
	g := NewSchemaGuardrail(schema)

	if g.Name() != "SchemaGuardrail" {
		t.Errorf("expected name 'SchemaGuardrail', got '%s'", g.Name())
	}

	// Should pass — valid JSON matching schema
	validJSON := `{"name": "test", "score": 95}`
	if err := g.Validate(validJSON); err != nil {
		t.Errorf("expected pass for valid JSON, got error: %v", err)
	}

	// Should fail — not JSON at all
	if err := g.Validate("this is not json"); err == nil {
		t.Error("expected failure for non-JSON, got nil")
	}
}

func TestRunAll(t *testing.T) {
	g1 := NewMaxTokenGuardrail(100)
	g2, _ := NewContentFilterGuardrail([]string{`(?i)forbidden`})

	// Should pass all
	if err := RunAll([]Guardrail{g1, g2}, "This is a clean output."); err != nil {
		t.Errorf("expected pass, got error: %v", err)
	}

	// Should fail content filter
	if err := RunAll([]Guardrail{g1, g2}, "This contains forbidden content."); err == nil {
		t.Error("expected failure, got nil")
	}
}
