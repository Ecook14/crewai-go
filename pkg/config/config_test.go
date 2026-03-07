package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAgents(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "agents.yaml")
	yaml := `researcher:
  role: "Senior Researcher"
  goal: "Find answers"
  backstory: "Expert investigator"
  verbose: true
writer:
  role: "Content Writer"
  goal: "Write articles"
  backstory: "Skilled wordsmith"
`
	os.WriteFile(configPath, []byte(yaml), 0644)

	agents, err := LoadAgents(configPath)
	if err != nil {
		t.Fatalf("LoadAgents failed: %v", err)
	}

	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d", len(agents))
	}

	researcher, ok := agents["researcher"]
	if !ok {
		t.Fatal("expected 'researcher' agent")
	}
	if researcher.Role != "Senior Researcher" {
		t.Errorf("expected role 'Senior Researcher', got '%s'", researcher.Role)
	}
	if !researcher.Verbose {
		t.Error("expected researcher.Verbose to be true")
	}

	writer, ok := agents["writer"]
	if !ok {
		t.Fatal("expected 'writer' agent")
	}
	if writer.Goal != "Write articles" {
		t.Errorf("expected goal 'Write articles', got '%s'", writer.Goal)
	}
}

func TestLoadTasks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create agents config
	agentsPath := filepath.Join(tmpDir, "agents.yaml")
	agentsYaml := `researcher:
  role: "Researcher"
  goal: "Research topics"
  backstory: "Expert"
`
	os.WriteFile(agentsPath, []byte(agentsYaml), 0644)

	agents, err := LoadAgents(agentsPath)
	if err != nil {
		t.Fatalf("LoadAgents failed: %v", err)
	}

	// Create tasks config
	tasksPath := filepath.Join(tmpDir, "tasks.yaml")
	tasksYaml := `research_task:
  description: "Research AI trends"
  agent: "researcher"
unbound_task:
  description: "No agent bound"
`
	os.WriteFile(tasksPath, []byte(tasksYaml), 0644)

	tasks, err := LoadTasks(tasksPath, agents)
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	rt, ok := tasks["research_task"]
	if !ok {
		t.Fatal("expected 'research_task'")
	}
	if rt.Description != "Research AI trends" {
		t.Errorf("expected description 'Research AI trends', got '%s'", rt.Description)
	}
	if rt.Agent == nil {
		t.Error("expected research_task to have an agent bound")
	}

	ut, ok := tasks["unbound_task"]
	if !ok {
		t.Fatal("expected 'unbound_task'")
	}
	if ut.Agent != nil {
		t.Error("expected unbound_task to have nil agent")
	}
}

func TestLoadAgentsInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "bad.yaml")
	os.WriteFile(configPath, []byte("{{invalid yaml"), 0644)

	_, err := LoadAgents(configPath)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoadAgentsMissingFile(t *testing.T) {
	_, err := LoadAgents("/nonexistent/path/agents.yaml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoadTasksMissingFile(t *testing.T) {
	_, err := LoadTasks("/nonexistent/path/tasks.yaml", nil)
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
