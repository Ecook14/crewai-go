package memory

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore is a persistent implementation of the Store interface using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore initializes a new SQLite database for persistent memory.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// Create table for memory items
	createTable := `
	CREATE TABLE IF NOT EXISTS memory_items (
		id TEXT PRIMARY KEY,
		text TEXT,
		vector TEXT,
		metadata TEXT
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("failed to create memory_items table: %w", err)
	}

	// Performance and Concurrency tuning
	db.Exec("PRAGMA journal_mode=WAL")
	db.Exec("PRAGMA synchronous=NORMAL")

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(ctx context.Context, item *MemoryItem) error {
	vectorJSON, err := json.Marshal(item.Vector)
	if err != nil {
		return fmt.Errorf("failed to marshal vector: %w", err)
	}

	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `INSERT OR REPLACE INTO memory_items (id, text, vector, metadata) VALUES (?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query, item.ID, item.Text, string(vectorJSON), string(metadataJSON))
	if err != nil {
		return fmt.Errorf("failed to insert memory item: %w", err)
	}

	return nil
}

func (s *SQLiteStore) Search(ctx context.Context, queryVector []float32, limit int) ([]*MemoryItem, error) {
	// Elite Implementation: Brute-force Cosine Similarity scan for high-precision retrieval
	// on standard persistent datasets. For billion-scale, use Chroma/LanceDB.
	
	rows, err := s.db.QueryContext(ctx, `SELECT id, text, vector, metadata FROM memory_items`)
	if err != nil {
		return nil, fmt.Errorf("failed to query memory items: %w", err)
	}
	defer rows.Close()

	var items []*MemoryItem
	for rows.Next() {
		var item MemoryItem
		var vectorJSON, metadataJSON string
		if err := rows.Scan(&item.ID, &item.Text, &vectorJSON, &metadataJSON); err != nil {
			return nil, fmt.Errorf("failed to scan memory item: %w", err)
		}

		if err := json.Unmarshal([]byte(vectorJSON), &item.Vector); err != nil {
			return nil, fmt.Errorf("failed to unmarshal vector: %w", err)
		}
		if err := json.Unmarshal([]byte(metadataJSON), &item.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}

		items = append(items, &item)
	}

	type scoredItem struct {
		item  *MemoryItem
		score float32
	}

	var results []scoredItem
	for _, item := range items {
		if len(item.Vector) != len(queryVector) {
			continue
		}
		sim, err := CosineSimilarity(queryVector, item.Vector)
		if err != nil {
			return nil, err
		}
		results = append(results, scoredItem{item: item, score: sim})
	}

	// Sort highest score first
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	var out []*MemoryItem
	for i := 0; i < limit && i < len(results); i++ {
		out = append(out, results[i].item)
	}

	return out, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
