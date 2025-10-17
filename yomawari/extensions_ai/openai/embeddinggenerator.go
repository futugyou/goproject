package openai

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
	rawopenai "github.com/openai/openai-go/v3"
)

type OpenAIEmbeddingGenerator[TEmbedding embeddings.EmbeddingT[float64]] struct {
	metadata     *embeddings.EmbeddingGeneratorMetadata
	openAIClient *rawopenai.Client
	modelId      *string
	dimensions   *int64
}

func NewOpenAIEmbeddingGenerator[TEmbedding embeddings.EmbeddingT[float64]](
	openAIClient *rawopenai.Client,
	modelId *string,
	dimensions *int64,
) *OpenAIEmbeddingGenerator[TEmbedding] {
	name := "openai"
	return &OpenAIEmbeddingGenerator[TEmbedding]{
		metadata:     &embeddings.EmbeddingGeneratorMetadata{ProviderName: &name, DefaultModelId: modelId},
		openAIClient: openAIClient,
		modelId:      modelId,
		dimensions:   dimensions,
	}
}

func (g *OpenAIEmbeddingGenerator[TEmbedding]) Generate(ctx context.Context, values []string, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[embeddings.EmbeddingT[float64]], error) {
	if options == nil {
		return nil, fmt.Errorf("no option info")
	}
	body := ToOpenAIEmbeddingParams(values, options)
	res, err := g.openAIClient.Embeddings.New(ctx, *body)
	if err != nil {
		return nil, err
	}

	return ToGeneratedEmbeddings(res), nil
}
