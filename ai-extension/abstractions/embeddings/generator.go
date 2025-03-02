package embeddings

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type IEmbeddingGenerator[TInput any, TEmbedding IEmbedding] interface {
	Generate(ctx context.Context, values []TInput, options *EmbeddingGenerationOptions) (GeneratedEmbeddings[TEmbedding], error)
}

func GenerateEmbeddingVector[TInput, TEmbeddingElement any, TEmbedding Embedding](
	generator IEmbeddingGenerator[TInput, *EmbeddingT[TEmbeddingElement]],
	ctx context.Context,
	value TInput,
	options *EmbeddingGenerationOptions,
) ([]TEmbeddingElement, error) {
	embedding, err := GenerateEmbedding(generator, ctx, value, options)
	if err != nil {
		return nil, err
	}
	return embedding.Vector, nil
}

func GenerateEmbedding[TInput, TEmbeddingElement any, TEmbedding Embedding](
	generator IEmbeddingGenerator[TInput, *EmbeddingT[TEmbeddingElement]],
	ctx context.Context,
	value TInput,
	options *EmbeddingGenerationOptions,
) (*EmbeddingT[TEmbeddingElement], error) {
	if generator == nil {
		return nil, errors.New("generator cannot be nil")
	}

	if isNil(value) {
		return nil, errors.New("value cannot be nil")
	}

	embeddings, err := generator.Generate(ctx, []TInput{value}, options)
	if err != nil {
		return nil, err
	}

	if embeddings.Count() != 1 {
		return nil, fmt.Errorf("expected the number of embeddings (%d) to match the number of inputs (1)", embeddings.Count())
	}

	return embeddings.Get(0), nil
}

func isNil(value interface{}) bool {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return false
}
