package openai

import (
	"context"

	"github.com/futugyou/yomawari/generative-ai/abstractions/embeddings"
	rawopenai "github.com/openai/openai-go"
)

type OpenAIEmbeddingGenerator[TInput string, TEmbedding embeddings.EmbeddingT[float64]] struct {
	metadata     *embeddings.EmbeddingGeneratorMetadata
	openAIClient *rawopenai.Client
	modelId      *string
	dimensions   *int64
}

func NewOpenAIEmbeddingGenerator[TInput string, TEmbedding embeddings.EmbeddingT[float64]](
	openAIClient *rawopenai.Client,
	modelId *string,
	dimensions *int64,
) *OpenAIEmbeddingGenerator[TInput, TEmbedding] {
	name := "openai"
	return &OpenAIEmbeddingGenerator[TInput, TEmbedding]{
		metadata:     &embeddings.EmbeddingGeneratorMetadata{ProviderName: &name, ModelId: modelId},
		openAIClient: openAIClient,
		modelId:      modelId,
		dimensions:   dimensions,
	}
}

func (g *OpenAIEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[embeddings.EmbeddingT[float64]], error) {
	body := ToOpenAIEmbeddingParams[TInput](values, options)
	res ,err:=g.openAIClient.Embeddings.New(ctx, *body)
	if err != nil {
		return nil, err
	}
	 
	return ToGeneratedEmbeddings(res), nil
}
