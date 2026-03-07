package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ArxivTool allows agents to search for academic papers.
type ArxivTool struct{}

func NewArxivTool() *ArxivTool {
	return &ArxivTool{}
}

func (t *ArxivTool) Name() string { return "ArxivTool" }

func (t *ArxivTool) Description() string {
	return "Searches arXiv.org for academic papers. Action: search."
}

func (t *ArxivTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	query, ok := input["query"].(string)
	if !ok {
		return "", fmt.Errorf("missing 'query'")
	}

	apiURL := fmt.Sprintf("http://export.arxiv.org/api/query?search_query=all:%s&start=0&max_results=3", url.QueryEscape(query))
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// Actual Implementation: Returns raw XML response from arXiv API for parsing.
	return string(body), nil
}

func (t *ArxivTool) RequiresReview() bool { return false }
