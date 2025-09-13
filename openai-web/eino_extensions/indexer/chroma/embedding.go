package chroma

import (
	"context"
	"fmt"

	"github.com/amikos-tech/chroma-go/types"
	"github.com/cloudwego/eino/components/embedding"
)

type EinoEmbedderAdapter struct {
	embedder embedding.Embedder
}

// EmbedDocuments implements the chroma.EmbeddingFunction interface.
func (e *EinoEmbedderAdapter) EmbedDocuments(ctx context.Context, texts []string) ([]*types.Embedding, error) {
	// 1. Call the original EmbedStrings method.
	vectors, err := e.embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	// 2. Convert the [][]float64 to []*chroma.Embedding.
	chromaEmbeddings := make([]*types.Embedding, len(vectors))
	for i, vec64 := range vectors {
		// Convert []float64 to []float32.
		vec32 := make([]float32, len(vec64))
		for j, val := range vec64 {
			vec32[j] = float32(val)
		}

		// 3. Assign the new []float32 to the ArrayOfFloat32 field.
		chromaEmbeddings[i] = &types.Embedding{
			ArrayOfFloat32: &vec32,
		}
	}
	return chromaEmbeddings, nil
}

// EmbedQuery implements the chroma.EmbeddingFunction interface.
func (e *EinoEmbedderAdapter) EmbedQuery(ctx context.Context, text string) (*types.Embedding, error) {
	// 1. Call the original EmbedStrings for a single text.
	vectors, err := e.embedder.EmbedStrings(ctx, []string{text})
	if err != nil {
		return nil, fmt.Errorf("embedding query failed: %w", err)
	}

	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return nil, fmt.Errorf("embedding query returned no vectors")
	}

	// 2. Convert the []float64 to []float32.
	vec64 := vectors[0]
	vec32 := make([]float32, len(vec64))
	for i, val := range vec64 {
		vec32[i] = float32(val)
	}

	// 3. Assign the new []float32 to the ArrayOfFloat32 field.
	return &types.Embedding{
		ArrayOfFloat32: &vec32,
	}, nil
}

// EmbedRecords implements the chroma.EmbeddingFunction interface.
// This is specific to Chroma's internal record structure.
// Since your eino embedder works on raw strings, you'll need to extract them from the records.
func (e *EinoEmbedderAdapter) EmbedRecords(ctx context.Context, records []*types.Record, force bool) error {
	textsToEmbed := make([]string, 0, len(records))
	for _, record := range records {
		// Only embed if the vector is not already set, or if `force` is true.
		if record.Embedding.ArrayOfFloat32 == nil || force {
			textsToEmbed = append(textsToEmbed, record.Document)
		}
	}

	if len(textsToEmbed) == 0 {
		return nil
	}

	// Get the embeddings for all the documents that need them.
	embeddings, err := e.EmbedDocuments(ctx, textsToEmbed)
	if err != nil {
		return err
	}

	// Map the new embeddings back to the original records.
	var currentEmbedIndex int
	for _, record := range records {
		if record.Embedding.ArrayOfFloat32 == nil || force {
			if currentEmbedIndex >= len(embeddings) {
				return fmt.Errorf("mismatch between number of documents and embeddings")
			}
			if embeddings[currentEmbedIndex] != nil {
				record.Embedding = *embeddings[currentEmbedIndex]
				currentEmbedIndex++
			}
		}
	}

	return nil
}
