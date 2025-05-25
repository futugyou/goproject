package embeddings

import (
	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
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

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) ConfigureOptions(configure func(*embeddings.EmbeddingGenerationOptions)) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	b.Use(func(innerClient embeddings.IEmbeddingGenerator[TInput, TEmbedding], sp core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
		return NewConfigureOptionsEmbeddingGenerator(innerClient, configure)
	})
	return b
}

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) UseDistributedCache(
	storage core.IDistributedCache,
	configure func(*DistributedCachingEmbeddingGenerator[TInput, TEmbedding]),
) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	b.Use(func(innerClient embeddings.IEmbeddingGenerator[TInput, TEmbedding], sp core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
		if storage == nil {
			storage, _ = core.GetService[core.IDistributedCache](sp)
		}

		var chatClient = NewDistributedCachingEmbeddingGenerator[TInput, TEmbedding](innerClient, storage)
		if configure != nil {
			configure(chatClient)
		}

		return chatClient
	})
	return b
}

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) UseLogging(configure func(*LoggingEmbeddingGenerator[TInput, TEmbedding])) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	return b.Use(func(client embeddings.IEmbeddingGenerator[TInput, TEmbedding], services core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
		logger, _ := core.GetService[logger.Logger](services)
		metadata, _ := core.GetService[*embeddings.EmbeddingGeneratorMetadata](services)
		logclient := NewLoggingEmbeddingGenerator[TInput, TEmbedding](
			client,
			logger,
			metadata,
		)

		if configure != nil {
			configure(logclient)
		}

		return logclient
	})
}

func (b *EmbeddingGeneratorBuilder[TInput, TEmbedding]) UseOpenTelemetry(configure func(*OpenTelemetryEmbeddingGenerator[TInput, TEmbedding])) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	return b.Use(func(client embeddings.IEmbeddingGenerator[TInput, TEmbedding], services core.IServiceProvider) embeddings.IEmbeddingGenerator[TInput, TEmbedding] {
		metadata, _ := core.GetService[*embeddings.EmbeddingGeneratorMetadata](services)
		otelclient := NewOpenTelemetryEmbeddingGenerator[TInput, TEmbedding](
			client,
			metadata,
		)

		if configure != nil {
			configure(otelclient)
		}

		return otelclient
	})
}

func IEmbeddingGeneratorAsBuilder[TInput any, TEmbedding embeddings.IEmbedding](
	innerClient embeddings.IEmbeddingGenerator[TInput, TEmbedding],
) *EmbeddingGeneratorBuilder[TInput, TEmbedding] {
	return NewEmbeddingGeneratorBuilder(innerClient)
}
