package memory

import (
	"fmt"
	"math"
)

// CosineSimilarity calculates the cosine similarity between two float slices.
// Returns a value between -1.0 and 1.0 (closer to 1.0 means higher similarity).
func CosineSimilarity(a, b []float32) (float32, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("vector lengths do not match: %d != %d", len(a), len(b))
	}

	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0, nil // Handle null vectors linearly
	}

	sim := dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
	return sim, nil
}
