package search

import (
	contentRaw "context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type ISearchClient interface {
	Search(
		ctx contentRaw.Context,
		index string,
		query string,
		filters []models.MemoryFilter,
		minRelevance float64,
		limit int,
		context context.IContext) (*models.SearchResult, error)

	Ask(
		ctx contentRaw.Context,
		index string,
		question string,
		filters []models.MemoryFilter,
		minRelevance float64,
		context context.IContext) (*models.MemoryAnswer, error)

	AskStreaming(
		ctx contentRaw.Context,
		index string,
		question string,
		filters []models.MemoryFilter,
		minRelevance float64,
		context context.IContext) <-chan AskStreamingStreamResponse

	ListIndexes(ctx contentRaw.Context) ([]string, error)
}

type AskStreamingStreamResponse struct {
	Answer *models.MemoryAnswer
	Err    error
}
