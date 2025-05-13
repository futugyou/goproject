package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/embeddings"
)

type OllamaEmbeddingGenerator[TInput string, TEmbedding embeddings.EmbeddingT[float64]] struct {
	metadata        *embeddings.EmbeddingGeneratorMetadata
	apiChatEndpoint url.URL
	httpClient      *http.Client
}

func NewOllamaEmbeddingGenerator[TInput string, TEmbedding embeddings.EmbeddingT[float64]](modelId *string, endpoint url.URL, httpClient *http.Client) *OllamaEmbeddingGenerator[TInput, TEmbedding] {
	var _ embeddings.IEmbedding = (*embeddings.EmbeddingT[float64])(nil)
	var _ embeddings.IEmbeddingGenerator[TInput, embeddings.EmbeddingT[float64]] = (*OllamaEmbeddingGenerator[TInput, embeddings.EmbeddingT[float64]])(nil)

	name := "ollama"
	apiChatEndpoint := endpoint.JoinPath("api/embed")
	if httpClient == nil {
		httpClient = SharedClient
	}

	return &OllamaEmbeddingGenerator[TInput, TEmbedding]{
		metadata: &embeddings.EmbeddingGeneratorMetadata{
			ProviderName:   &name,
			ProviderUri:    &endpoint,
			DefaultModelId: modelId,
		},
		apiChatEndpoint: *apiChatEndpoint,
		httpClient:      httpClient,
	}
}

func (g *OllamaEmbeddingGenerator[TInput, TEmbedding]) Generate(ctx context.Context, values []TInput, options *embeddings.EmbeddingGenerationOptions) (*embeddings.GeneratedEmbeddings[embeddings.EmbeddingT[float64]], error) {
	request := OllamaEmbeddingRequest{
		Input:   ToStringList(values),
		Options: &OllamaRequestOptions{},
	}

	requestModel := ""
	if options != nil {
		requestModel = *options.ModelId
		if len(requestModel) == 0 && g.metadata != nil && g.metadata.DefaultModelId != nil {
			requestModel = *g.metadata.DefaultModelId
		}
		request.Model = requestModel

		if options.AdditionalProperties != nil {
			if v, ok := options.AdditionalProperties["keep_alive"]; ok {
				if v, ok := v.(int64); ok {
					request.KeepAlive = &v
				}
			}
			if v, ok := options.AdditionalProperties["truncate"]; ok {
				if v, ok := v.(bool); ok {
					request.Truncate = v
				}
			}
		}
	}

	path := g.apiChatEndpoint.String()
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("POST", path, body)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var embeddingResponse OllamaEmbeddingResponse
	err = json.NewDecoder(resp.Body).Decode(&embeddingResponse)
	if err != nil {
		return nil, err
	}

	if len(embeddingResponse.Embeddings) != len(values) {
		return nil, fmt.Errorf("ollama generated  %d embeddings but %d were expected", len(embeddingResponse.Embeddings), len(values))
	}

	result := ToGeneratedEmbeddingsResponse(embeddingResponse, requestModel)
	return result, nil
}

func ToGeneratedEmbeddingsResponse(embeddingResponse OllamaEmbeddingResponse, requestModel string) *embeddings.GeneratedEmbeddings[embeddings.EmbeddingT[float64]] {
	metadata := map[string]int64{}
	TransferNanosecondsTime(embeddingResponse, func(response OllamaEmbeddingResponse) *int64 { return response.LoadDuration }, "load_duration", &metadata)
	TransferNanosecondsTime(embeddingResponse, func(response OllamaEmbeddingResponse) *int64 { return response.TotalDuration }, "total_duration", &metadata)

	var useage abstractions.UsageDetails
	if len(metadata) > 0 || embeddingResponse.PromptEvalCount != nil {
		useage = abstractions.UsageDetails{
			InputTokenCount:      embeddingResponse.PromptEvalCount,
			TotalTokenCount:      embeddingResponse.PromptEvalCount,
			AdditionalProperties: metadata,
		}
	}
	embeddinggs := []embeddings.EmbeddingT[float64]{}
	for _, embe := range embeddingResponse.Embeddings {
		createAt := time.Now()
		if embeddingResponse.Model != nil {
			requestModel = *embeddingResponse.Model
		}
		embeddinggs = append(embeddinggs, embeddings.EmbeddingT[float64]{
			Embedding: embeddings.Embedding{
				CreatedAt:            &createAt,
				ModelId:              &requestModel,
				AdditionalProperties: map[string]interface{}{},
			},
			Vector: embe,
		})
	}

	result := embeddings.NewGeneratedEmbeddingsFromCollection(embeddinggs)
	result.Usage = &useage
	return result
}

func ToStringList[TInput string](values []TInput) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = string(value)
	}
	return result
}
