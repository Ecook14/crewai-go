package knowledge

import (
	"context"
	"os"
	"testing"

	"github.com/Ecook14/crewai-go/pkg/llm"
	"github.com/Ecook14/crewai-go/pkg/memory"
)

func TestTokenSplitter(t *testing.T) {
	splitter := NewTokenSplitter(5, 1)
	text := "one two three four five six seven eight nine ten"
	chunks := splitter.SplitText(text)

	if len(chunks) != 3 {
		t.Errorf("Expected 3 chunks, got %d", len(chunks))
	}

	// Chunk 1: one two three four five
	// Chunk 2: five six seven eight nine (starts at index 5-1=4)
	// Chunk 3: nine ten (starts at index 9-1=8)
	
	if chunks[0] != "one two three four five" {
		t.Errorf("Unexpected chunk 0: %s", chunks[0])
	}
}

type mockLLM struct {
	llm.Client
	embedFunc func(text string) ([]float32, error)
}

func (m *mockLLM) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	return m.embedFunc(text)
}

type mockStore struct {
	memory.Store
	items []*memory.MemoryItem
}

func (m *mockStore) Add(ctx context.Context, item *memory.MemoryItem) error {
	m.items = append(m.items, item)
	return nil
}

func TestIngestionEngine(t *testing.T) {
	tmpFile := "test_ingest.txt"
	content := "This is a test file for RAG ingestion. It should be split into multiple chunks."
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	mockLLM := &mockLLM{
		embedFunc: func(text string) ([]float32, error) {
			return []float32{0.1, 0.2}, nil
		},
	}
	mockStore := &mockStore{}
	splitter := NewTokenSplitter(5, 0)
	
	ie := &IngestionEngine{
		Store:    mockStore,
		LLM:      mockLLM,
		Splitter: splitter,
	}

	err := ie.IngestFile(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("IngestFile failed: %v", err)
	}

	if len(mockStore.items) == 0 {
		t.Errorf("No items ingested into store")
	}
}
