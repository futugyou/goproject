package embeddings

import (
	"github.com/futugyou/ai-extension/abstractions/embeddings"
	"github.com/futugyou/ai-extension/core"
)

type EmbeddingGeneratorBuilder[TInput any, TEmbedding embeddings.IEmbedding] struct {
	innerGeneratorFactory func(core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding]
	generatorFactories    []func(embeddings.IEmbeddingGenerator[TInput, TEmbedding], core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding]
}

func NewEmbeddingGeneratorBuilder[TInput any, TEmbedding embeddings.IEmbedding](
	innerGenerator embeddings.IEmbeddingGenerator[TInput, TEmbedding],
) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	return &EmbeddingGeneratorBuilder[TInput, TEmbedding]{
		innerGeneratorFactory: func(_ core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
			return innerGenerator
		},
	}
}

func NewEmbeddingGeneratorBuilderWithIServiceProvider[TInput any, TEmbedding embeddings.IEmbedding](
	innerGeneratorFactory func(core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding],
) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	return &EmbeddingGeneratorBuilder[TInput, TEmbedding]{
		innerGeneratorFactory: innerGeneratorFactory,
	}
}

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) Build(
	services core.IServiceProvider,
) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
	if services == nil {
		services = &core.ServiceProvider{}
	}
	embeddingGenerator := b.innerGeneratorFactory(services)

	for i := len(b.generatorFactories) - 1; i >= 0; i-- {
		factory := b.generatorFactories[i]
		embeddingGenerator = factory(embeddingGenerator, services)
	}

	return embeddingGenerator
}

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) Use(
	generatorFactory func(embeddings.IEmbeddingGenerator[TInput, TEmbedding], core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding],
) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	if b.generatorFactories == nil {
		b.generatorFactories = []func(embeddings.IEmbeddingGenerator[TInput, TEmbedding], core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding]{}
	}

	b.generatorFactories = append(b.generatorFactories, generatorFactory)

	return b
}

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) UseFactory(
	generatorFactory func(embeddings.IEmbeddingGenerator[TInput, TEmbedding]) embeddings.IEmbeddingGenerator[TInput, TEmbedding],
) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	return b.Use(func(generator embeddings.IEmbeddingGenerator[TInput, TEmbedding], _ core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
		return generatorFactory(generator)
	})
}

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) UseAnonymousGenerator(
	enerateFunc GenerateFunc[TInput, TEmbedding],
) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	return b.Use(func(generator embeddings.IEmbeddingGenerator[TInput, TEmbedding], _ core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
		return NewAnonymousDelegatingEmbeddingGenerator(generator, enerateFunc)
	})
}
