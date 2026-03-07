package delegation

import (
	"context"
	"testing"
)

// mockAgent implements the delegation.Agent interface for testing.
type mockAgent struct {
	role   string
	result string
}

func (m *mockAgent) GetRole() string { return m.role }
func (m *mockAgent) Execute(ctx context.Context, taskInput string, options map[string]interface{}) (interface{}, error) {
	return m.result + " (received: " + taskInput + ")", nil
}

func TestDelegateWorkTool(t *testing.T) {
	coworker := &mockAgent{role: "Researcher", result: "Research results"}
	tool := NewDelegateWorkTool([]Agent{coworker})

	if tool.Name() != "DelegateWork" {
		t.Errorf("expected name 'DelegateWork', got '%s'", tool.Name())
	}

	result, err := tool.Execute(context.Background(), map[string]interface{}{
		"coworker": "Researcher",
		"task":     "Find AI trends",
		"context":  "Focus on 2026",
	})
	if err != nil {
		t.Fatalf("DelegateWorkTool.Execute failed: %v", err)
	}

	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestDelegateWorkToolMissingCoworker(t *testing.T) {
	tool := NewDelegateWorkTool([]Agent{})
	_, err := tool.Execute(context.Background(), map[string]interface{}{
		"task": "something",
	})
	if err == nil {
		t.Error("expected error for missing coworker, got nil")
	}
}

func TestDelegateWorkToolUnknownCoworker(t *testing.T) {
	coworker := &mockAgent{role: "Writer", result: "text"}
	tool := NewDelegateWorkTool([]Agent{coworker})

	_, err := tool.Execute(context.Background(), map[string]interface{}{
		"coworker": "NonExistent",
		"task":     "something",
	})
	if err == nil {
		t.Error("expected error for unknown coworker, got nil")
	}
}

func TestAskQuestionToolDelegation(t *testing.T) {
	coworker := &mockAgent{role: "Expert", result: "Expert answer"}
	tool := NewAskQuestionTool([]Agent{coworker})

	if tool.Name() != "AskQuestion" {
		t.Errorf("expected name 'AskQuestion', got '%s'", tool.Name())
	}

	result, err := tool.Execute(context.Background(), map[string]interface{}{
		"coworker": "Expert",
		"question": "What is Go?",
	})
	if err != nil {
		t.Fatalf("AskQuestionTool.Execute failed: %v", err)
	}

	if result == "" {
		t.Error("expected non-empty result")
	}
}
