// Package store provides an in-memory DuckDB vector index over vault
// documents, using the VSS extension for cosine-similarity search.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ardanlabs/kronk/sdk/kronk"
	"github.com/ardanlabs/kronk/sdk/kronk/model"
	_ "github.com/duckdb/duckdb-go/v2"

	"github.com/jeff/mcp-vault/internal/vault"
)

// Result is a single vault document returned by Search, ranked by
// similarity to the query embedding.
type Result struct {
	Path       string  `json:"path"`
	Text       string  `json:"text"`
	Similarity float64 `json:"similarity"`
}

// LoadVault creates an in-memory DuckDB database, embeds every document in
// docs using krnEmbed, and indexes the embeddings with an HNSW index for
// cosine-similarity search. dims must match the native embedding dimension
// of krnEmbed.
func LoadVault(ctx context.Context, krnEmbed *kronk.Kronk, docs []vault.Document, dims int) (*sql.DB, error) {
	db, err := sql.Open("duckdb", "")
	if err != nil {
		return nil, fmt.Errorf("open duckdb: %w", err)
	}

	if _, err := db.Exec("INSTALL vss; LOAD vss;"); err != nil {
		return nil, fmt.Errorf("load vss extension: %w", err)
	}

	if _, err := db.Exec("SET hnsw_enable_experimental_persistence = true;"); err != nil {
		return nil, fmt.Errorf("set hnsw persistence: %w", err)
	}

	createTable := fmt.Sprintf(`
		CREATE TABLE items (
			id        INTEGER PRIMARY KEY,
			path      VARCHAR,
			text      VARCHAR,
			embedding FLOAT[%d]
		);
	`, dims)

	if _, err := db.Exec(createTable); err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}

	for i, doc := range docs {
		vec, err := embed(ctx, krnEmbed, doc.Content)
		if err != nil {
			return nil, fmt.Errorf("embed %s: %w", doc.Path, err)
		}

		if err := insert(db, i, doc.Path, doc.Content, vec); err != nil {
			return nil, fmt.Errorf("insert %s: %w", doc.Path, err)
		}
	}

	createIndex := `
		CREATE INDEX idx_embedding ON items
		USING HNSW (embedding)
		WITH (metric = 'cosine');
	`

	if _, err := db.Exec(createIndex); err != nil {
		return nil, fmt.Errorf("create hnsw index: %w", err)
	}

	return db, nil
}

func embed(ctx context.Context, krnEmbed *kronk.Kronk, text string) ([]float32, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	d := model.D{
		"input":              text,
		"truncate":           true,
		"truncate_direction": "right",
	}

	resp, err := krnEmbed.Embeddings(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("embeddings: %w", err)
	}

	if len(resp.Data) == 0 || len(resp.Data[0].Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding response")
	}

	return resp.Data[0].Embedding, nil
}

func insert(db *sql.DB, id int, path, text string, vec []float32) error {
	text = strings.ReplaceAll(text, "'", "''")
	path = strings.ReplaceAll(path, "'", "''")
	vecStr := strings.ReplaceAll(fmt.Sprintf("%v", vec), " ", ",")

	sqlStr := fmt.Sprintf("INSERT INTO items (id, path, text, embedding) VALUES (%d, '%s', '%s', %s);",
		id, path, text, vecStr)

	if _, err := db.Exec(sqlStr); err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	return nil
}

// Search returns the topK documents most similar to queryVector, ordered by
// descending cosine similarity.
func Search(db *sql.DB, queryVector []float32, topK int) ([]Result, error) {
	query := fmt.Sprintf(`
		SELECT
			path,
			text,
			array_cosine_similarity(embedding, ?::FLOAT[%d]) AS similarity
		FROM items
		ORDER BY similarity DESC
		LIMIT %d;
	`, len(queryVector), topK)

	rows, err := db.Query(query, queryVector)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var results []Result

	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.Path, &r.Text, &r.Similarity); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		results = append(results, r)
	}

	return results, rows.Err()
}
