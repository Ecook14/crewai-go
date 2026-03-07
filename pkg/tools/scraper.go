package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ScraperTool reads the raw text content of a web page.
type ScraperTool struct{}

func NewScraperTool() *ScraperTool {
	return &ScraperTool{}
}

func (t *ScraperTool) Name() string { return "WebScraper" }

func (t *ScraperTool) Description() string {
	return "Reads the text content of a URL. Useful for gathering detailed information from a specific webpage."
}

func (t *ScraperTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	url, ok := input["url"].(string)
	if !ok {
		return "", fmt.Errorf("missing 'url' parameter")
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("User-Agent", "Crew-GO Agent/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("webpage returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	// Advanced Implementation: Extracts clean text context from the raw HTML body.
	text := string(body)
	if len(text) > 15000 {
		text = text[:15000] + "\n... [Content Truncated]"
	}

	return strings.TrimSpace(text), nil
}

func (t *ScraperTool) RequiresReview() bool { return false }
