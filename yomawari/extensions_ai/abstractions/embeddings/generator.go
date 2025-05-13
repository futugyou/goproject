package embeddings

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type IEmbeddingGenerator[TInput any, TEmbedding IEmbedding] interface {
	Generate(ctx context.Context, values []TInput, options *EmbeddingGenerationOptions) (*GeneratedEmbeddings[TEmbedding], error)
}

func GenerateVector[TInput any, TEmbeddingElement any](
	generator IEmbeddingGenerator[TInput, EmbeddingT[TEmbeddingElement]],
	ctx context.Context,
	value TInput,
	options *EmbeddingGenerationOptions,
) ([]TEmbeddingElement, error) {
	embedding, err := Generate(generator, ctx, value, options)
	if err != nil {
		return nil, err
	}
	return embedding.Vector, nil
}

func Generate[TInput any, TEmbedding IEmbedding](
	generator IEmbeddingGenerator[TInput, TEmbedding],
	ctx context.Context,
	value TInput,
	options *EmbeddingGenerationOptions,
) (*TEmbedding, error) {
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

	emb := embeddings.Get(0)
	return &emb, nil
}

func isNil(value interface{}) bool {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		return val.IsNil()
	}

	return false
}

type EmbeddingWithInput[TInput any, TEmbedding IEmbedding] struct {
	Value     TInput
	Embedding TEmbedding
}

func GenerateAndZip[TInput any, TEmbedding IEmbedding](
	generator IEmbeddingGenerator[TInput, TEmbedding],
	ctx context.Context,
	values []TInput,
	options *EmbeddingGenerationOptions,
) ([]EmbeddingWithInput[TInput, TEmbedding], error) {
	if generator == nil {
		return nil, errors.New("generator cannot be nil")
	}

	if len(values) == 0 {
		return nil, nil
	}

	embeddings, err := generator.Generate(ctx, values, options)
	if err != nil {
		return nil, err
	}

	if embeddings.Count() != len(values) {
		return nil, fmt.Errorf("expected the number of embeddings (%d) to match the number of inputs (%d)", embeddings.Count(), len(values))
	}

	var results []EmbeddingWithInput[TInput, TEmbedding]
	for i := 0; i < len(values); i++ {
		results = append(results, EmbeddingWithInput[TInput, TEmbedding]{
			Value:     values[i],
			Embedding: embeddings.Get(i),
		})
	}

	return results, nil
}
