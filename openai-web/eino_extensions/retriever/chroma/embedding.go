package chroma

import (
	"context"
	"fmt"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/cloudwego/eino/components/embedding"
)

type EinoEmbedderAdapter struct {
	embedder embedding.Embedder
}

func (e *EinoEmbedderAdapter) EmbedDocuments(ctx context.Context, documents []string) ([]embeddings.Embedding, error) {
	if len(documents) == 0 {
		return embeddings.NewEmptyEmbeddings(), nil
	}
	vectors, err := e.embedder.EmbedStrings(ctx, documents)
	if err != nil {
		return nil, fmt.Errorf("embedding docs failed: %w", err)
	}

	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return nil, fmt.Errorf("embedding docs returned no vectors")
	}

	return embeddings.NewEmbeddingsFromFloat32(float64ToFloat32(vectors))
}

func (e *EinoEmbedderAdapter) EmbedQuery(ctx context.Context, document string) (embeddings.Embedding, error) {
	vectors, err := e.embedder.EmbedStrings(ctx, []string{document})
	if err != nil {
		return nil, fmt.Errorf("embedding query failed: %w", err)
	}

	if len(vectors) == 0 || len(vectors[0]) == 0 {
		return nil, fmt.Errorf("embedding query returned no vectors")
	}
	vec64 := vectors[0]
	vec32 := make([]float32, len(vec64))
	for i, val := range vec64 {
		vec32[i] = float32(val)
	}
	return embeddings.NewEmbeddingFromFloat32(vec32), nil
}

func float64ToFloat32(vectors [][]float64) [][]float32 {
	result := make([][]float32, len(vectors))
	for i, row := range vectors {
		r := make([]float32, len(row))
		for j, val := range row {
			r[j] = float32(val)
		}
		result[i] = r
	}
	return result
}
