package pipeline

import (
	rawContext "context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type IPipelineOrchestrator interface {
	GetHandlerNames() []string
	SetHandlerNames(names []string)
	AddHandler(ctx rawContext.Context, handler IPipelineStepHandler) error
	TryAddHandler(ctx rawContext.Context, handler IPipelineStepHandler) error
	ImportDocument(ctx rawContext.Context, text string, uploadRequest *models.DocumentUploadRequest, context context.IContext) (*string, error)
	PrepareNewDocumentUpload(ctx rawContext.Context, index string, documentId string, tags *models.TagCollection,
		filesToUpload []models.UploadedFile, contextArgs map[string]interface{}) (*DataPipeline, error)
	RunPipeline(ctx rawContext.Context, pipeline *DataPipeline) error
	ReadPipelineStatus(ctx rawContext.Context, index string, documentId string) (*DataPipeline, error)
	ReadPipelineSummary(ctx rawContext.Context, index string, documentId string) (*models.DataPipelineStatus, error)
	IsDocumentReady(ctx rawContext.Context, index string, documentId string) (bool, error)
	StopAllPipelines(ctx rawContext.Context) error
	ReadFileAsStream(ctx rawContext.Context, pipeline *DataPipeline, fileName string) (*models.StreamableFileContent, error)
	ReadFile(ctx rawContext.Context, pipeline *DataPipeline, fileName string) ([]byte, error)
	ReadTextFile(ctx rawContext.Context, pipeline *DataPipeline, fileName string) (*string, error)
	WriteTextFile(ctx rawContext.Context, pipeline *DataPipeline, fileName string, fileContent string) error
	WriteFile(ctx rawContext.Context, pipeline *DataPipeline, fileName string, fileContent []byte) error
	GetEmbeddingGenerationEnabled() bool
	GetEmbeddingGenerators(ctx rawContext.Context) ([]ai.ITextEmbeddingGenerator, error)
	GetMemoryDbs(ctx rawContext.Context) ([]memorystorage.IMemoryDb, error)
	GetTextGenerators(ctx rawContext.Context) (ai.ITextGenerator, error)
	StartIndexDeletion(ctx rawContext.Context, index *string) error
	StartDocumentDeletion(ctx rawContext.Context, index *string, documentId string) error
}
