package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Ecook14/crewai-go/pkg/agents"
	"github.com/Ecook14/crewai-go/pkg/crew"
	"github.com/Ecook14/crewai-go/pkg/llm"
	"github.com/Ecook14/crewai-go/pkg/tasks"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	model := llm.NewOpenAIClient(apiKey)

	// Define an agent with a streaming callback
	writer := agents.NewAgent(
		"Poet",
		"Write a short, beautiful haiku about Go programming.",
		"A minimalist poet",
		model,
		agents.WithVerbose(true),
	)

	// Set the streaming callback to print tokens in real-time
	writer.StepStreamCallback = func(token string) {
		fmt.Print(token)
	}

	task := &tasks.Task{
		Description: "Write a haiku about Go.",
		Agent:       writer,
	}

	myCrew := crew.NewCrew(
		[]*agents.Agent{writer},
		[]*tasks.Task{task},
	)

	fmt.Println("🚀 Starting Streaming Demo (Tokens should appear one by one):")
	_, err := myCrew.Kickoff(context.Background())
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
	}
	fmt.Println("\n\n✅ Done.")
}
