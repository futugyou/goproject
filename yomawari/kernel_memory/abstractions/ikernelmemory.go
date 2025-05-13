package abstractions

import (
	contentRaw "context"
	"io"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/search"
)

type IKernelMemory interface {
	ImportDocument(
		ctx contentRaw.Context,
		document models.Document,
		index *string,
		steps []string,
		context context.IContext) (*string, error)

	ImportDocumentWithPath(
		ctx contentRaw.Context,
		filePath string,
		documentId *string,
		tags *models.TagCollection,
		index *string,
		steps []string,
		context context.IContext) (*string, error)

	ImportDocumentWithDocumentUpload(
		ctx contentRaw.Context,
		uploadRequest models.DocumentUploadRequest) (*string, error)

	ImportDocumentWithStream(
		ctx contentRaw.Context,
		content io.ReadCloser,
		fileName *string,
		documentId *string,
		tags *models.TagCollection,
		index *string,
		steps []string,
		context context.IContext) (*string, error)

	ImportText(
		ctx contentRaw.Context,
		text string,
		documentId *string,
		tags *models.TagCollection,
		index *string,
		steps []string,
		context context.IContext) (*string, error)

	ImportWebPage(
		ctx contentRaw.Context,
		url string,
		documentId *string,
		tags *models.TagCollection,
		index *string,
		steps []string,
		context context.IContext) (*string, error)

	ListIndexes(ctx contentRaw.Context) ([]models.IndexDetails, error)

	DeleteIndex(ctx contentRaw.Context, index *string) error

	DeleteDocument(ctx contentRaw.Context, documentId string, index *string) error

	IsDocumentReady(ctx contentRaw.Context, documentId string, index *string) bool

	GetDocumentStatus(ctx contentRaw.Context, documentId string, index *string) (*models.DataPipelineStatus, error)

	ExportFile(ctx contentRaw.Context, documentId string, fileName string, index *string) (*models.StreamableFileContent, error)

	Search(ctx contentRaw.Context,
		query string,
		index *string,
		filter *models.MemoryFilter,
		filters []models.MemoryFilter,
		minRelevance float64,
		limit int,
		context context.IContext) (*models.SearchResult, error)

	AskStreaming(ctx contentRaw.Context,
		question string,
		index *string,
		filter *models.MemoryFilter,
		filters []models.MemoryFilter,
		minRelevance float64,
		options *search.SearchOptions,
		context context.IContext) <-chan search.AskStreamingStreamResponse

	Ask(ctx contentRaw.Context,
		question string,
		index *string,
		filter *models.MemoryFilter,
		filters []models.MemoryFilter,
		minRelevance float64,
		options *search.SearchOptions,
		context context.IContext) (*models.MemoryAnswer, error)

	SearchSynthetics(
		ctx contentRaw.Context,
		syntheticType string,
		index *string,
		filter *models.MemoryFilter,
		filters []models.MemoryFilter) ([]models.Citation, error)

	SearchSummaries(
		ctx contentRaw.Context,
		index *string,
		filter *models.MemoryFilter,
		filters []models.MemoryFilter) ([]models.Citation, error)
}
