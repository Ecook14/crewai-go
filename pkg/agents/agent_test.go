package agents

import (
	"context"
	"errors"
	"testing"

	"github.com/Ecook14/crewai-go/pkg/llm"
	"github.com/Ecook14/crewai-go/pkg/tools"
)

// mockLLM is a simple mock for the llm.Client interface
type mockLLM struct {
	llm.Client
	generateFunc func(messages []llm.Message) (string, error)
}

func (m *mockLLM) Generate(ctx context.Context, messages []llm.Message, options map[string]interface{}) (string, error) {
	return m.generateFunc(messages)
}

func TestAgentExecute_Basic(t *testing.T) {
	mock := &mockLLM{
		generateFunc: func(messages []llm.Message) (string, error) {
			return "Success", nil
		},
	}

	agent := &Agent{
		Role: "Tester",
		Goal: "Test the agent",
		LLM:  mock,
	}

	result, err := agent.Execute(context.Background(), "Hello", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != "Success" {
		t.Errorf("Expected 'Success', got %v", result)
	}
}

func TestAgentExecute_SelfHealing(t *testing.T) {
	callCount := 0
	mock := &mockLLM{
		generateFunc: func(messages []llm.Message) (string, error) {
			callCount++
			if callCount == 1 {
				// Simulate tool call that will fail
				return `{"tool": "FailingTool", "input": {}}`, nil
			}
			return "Fixed after failure", nil
		},
	}

	failingTool := &mockTool{
		name: "FailingTool",
		executeFunc: func(input map[string]interface{}) (string, error) {
			return "", errors.New("tool injection failure")
		},
	}

	agent := &Agent{
		Role:        "Healer",
		Goal:        "Test self-healing",
		LLM:         mock,
		Tools:       []tools.Tool{failingTool},
		SelfHealing: true,
		MaxIter:     3,
	}

	result, err := agent.Execute(context.Background(), "Heal me", nil)
	if err != nil {
		t.Fatalf("Expected no error with self-healing, got %v", err)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls to LLM (fail then fix), got %d", callCount)
	}

	if result != "Fixed after failure" {
		t.Errorf("Expected fixed result, got %v", result)
	}
}

type mockTool struct {
	name           string
	executeFunc    func(input map[string]interface{}) (string, error)
	requiresReview bool
}

func (m *mockTool) RequiresReview() bool { return m.requiresReview }

func (m *mockTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	return m.executeFunc(input)
}

func TestAgentExecute_HITL(t *testing.T) {
	mock := &mockLLM{
		generateFunc: func(messages []llm.Message) (string, error) {
			return `{"tool": "ReviewedTool", "input": {}}`, nil
		},
	}

	tool := &mockTool{
		name:           "ReviewedTool",
		requiresReview: true,
		executeFunc: func(input map[string]interface{}) (string, error) {
			return "Tool Executed", nil
		},
	}

	// 1. Test Approval
	agent := &Agent{
		Role:  "HITLTester",
		LLM:   mock,
		Tools: []tools.Tool{tool},
		UsageMetrics: make(map[string]int),
		StepReview: func(toolName string, input map[string]interface{}) bool {
			return true // Approved
		},
	}

	res, err := agent.Execute(context.Background(), "Run tool", nil)
	if err != nil {
		t.Fatalf("Expected no error on approval, got %v", err)
	}
	if res != "Tool Executed" {
		t.Errorf("Expected 'Tool Executed', got %v", res)
	}

	// 2. Test Denial
	agent.StepReview = func(toolName string, input map[string]interface{}) bool {
		return false // Denied
	}

	// Re-mock LLM to return a final answer after denial is fed back as an observation
	denialCount := 0
	agent.LLM = &mockLLM{
		generateFunc: func(messages []llm.Message) (string, error) {
			for _, m := range messages {
				if strings.Contains(m.Content, "Tool Execution Denied") {
					return "Stopping because human denied.", nil
				}
			}
			denialCount++
			return `{"tool": "ReviewedTool", "input": {}}`, nil
		},
	}

	res, err = agent.Execute(context.Background(), "Run tool", nil)
	if err != nil {
		t.Fatalf("Expected no error on denial, got %v", err)
	}
	if !strings.Contains(res, "Stopping because human denied") {
		t.Errorf("Expected final answer after denial, got %v", res)
	}
}
