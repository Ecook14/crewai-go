package knowledge

import (
	"testing"
)

func TestTokenSplitterBasic(t *testing.T) {
	splitter := NewTokenSplitter(5, 1)
	text := "one two three four five six seven eight nine ten"

	chunks := splitter.SplitText(text)

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	// With chunk size 5 and overlap 1, stepping by 4
	// "one two three four five" (5 words)
	// "five six seven eight nine" (5 words, overlapping "five")
	// "nine ten" (2 words, final partial chunk)
	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(chunks))
	}
}

func TestTokenSplitterSmallDocument(t *testing.T) {
	splitter := NewTokenSplitter(100, 10)
	text := "short document"

	chunks := splitter.SplitText(text)
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk for small document, got %d", len(chunks))
	}
	if chunks[0] != "short document" {
		t.Errorf("expected 'short document', got '%s'", chunks[0])
	}
}

func TestTokenSplitterEmptyInput(t *testing.T) {
	splitter := NewTokenSplitter(10, 2)
	chunks := splitter.SplitText("")
	if chunks != nil {
		t.Errorf("expected nil for empty input, got %v", chunks)
	}
}

func TestTokenSplitterDefaults(t *testing.T) {
	// Test that negative/zero values get default values
	splitter := NewTokenSplitter(0, -1)
	if splitter.ChunkSize != 1000 {
		t.Errorf("expected default ChunkSize 1000, got %d", splitter.ChunkSize)
	}
	if splitter.ChunkOverlap != 100 {
		t.Errorf("expected default ChunkOverlap 100, got %d", splitter.ChunkOverlap)
	}
}

func TestTokenSplitterOverlap(t *testing.T) {
	splitter := NewTokenSplitter(3, 1)
	text := "a b c d e f"

	chunks := splitter.SplitText(text)

	// With chunk size 3 and overlap 1, step = 2
	// Chunk 1: "a b c"
	// Chunk 2: "c d e"
	// Chunk 3: "e f"
	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks with overlap, got %d", len(chunks))
	}
}
