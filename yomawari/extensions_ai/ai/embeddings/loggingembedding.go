package embeddings

import (
	"context"

	"github.com/futugyou/yomawari/core/logger"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
	"github.com/futugyou/yomawari/extensions_ai/ai"
)

type LoggingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding] struct {
	embeddings.DelegatingEmbeddingGenerator[TInput, TEmbedding]
	logger   logger.Logger
	metadata *embeddings.EmbeddingGeneratorMetadata
}

func NewLoggingEmbeddingGenerator[TInput any, TEmbedding embeddings.IEmbedding](
	innerGenerator embeddings.IEmbeddingGenerator[TInput, TEmbedding],
	logger logger.Logger,
	metadata *embeddings.EmbeddingGeneratorMetadata,
) *LoggingEmbeddingGenerator[TInput, TEmbedding] {
	return &LoggingEmbeddingGenerator[TInput, TEmbedding]{
		DelegatingEmbeddingGenerator: *embeddings.NewDelegatingEmbeddingGenerator[TInput, TEmbedding](innerGenerator),
		logger:                       logger,
		metadata:                     metadata,
	}
}

func (client *LoggingEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[TEmbedding], error) {
	if client.logger == nil {
		return client.DelegatingEmbeddingGenerator.Generate(ctx, values, options)
	}

	if client.logger.IsOutputLevelEnabled(logger.DebugLevel) {
		if client.logger.IsOutputLevelEnabled(logger.TraceLevel) {
			client.logger.Tracef("%s invoked: %s. Options: %s. Metadata: %s.",
				"Generate", ai.AsJson(values), ai.AsJson(options), ai.AsJson(client.metadata))
		} else {
			client.logger.Debugf("%s invoked", "Generate")
		}
	}

	response, err := client.DelegatingEmbeddingGenerator.Generate(ctx, values, options)
	if err == nil {
		if client.logger.IsOutputLevelEnabled(logger.DebugLevel) {
			if client.logger.IsOutputLevelEnabled(logger.TraceLevel) {
				client.logger.Tracef("%s completed: %s.", "Generate", ai.AsJson(response))
			} else {
				client.logger.Debugf("%s completed.", "Generate")
			}
		}
	} else {
		client.logger.Errorf("%s failed, err %s.", "Generate", err.Error())
	}

	return response, err
}
