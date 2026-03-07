package tools

import (
	"context"
	"time"
)

// DateTimeTool provides current date and time information.
type DateTimeTool struct{}

func NewDateTimeTool() *DateTimeTool {
	return &DateTimeTool{}
}

func (t *DateTimeTool) Name() string { return "DateTime" }

func (t *DateTimeTool) Description() string {
	return "Returns the current date and time."
}

func (t *DateTimeTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	return time.Now().Format(time.RFC3339), nil
}

func (t *DateTimeTool) RequiresReview() bool { return false }
