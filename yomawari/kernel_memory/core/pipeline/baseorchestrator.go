package pipeline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	kmcontext "github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
	"github.com/futugyou/yomawari/kernel_memory/core/configuration"
	"github.com/google/uuid"
)

type BaseOrchestrator struct {
	memoryDbs                  []memorystorage.IMemoryDb
	embeddingGenerators        []ai.ITextEmbeddingGenerator
	textGenerator              ai.ITextGenerator
	defaultIngestionSteps      []string
	documentStorage            documentstorage.IDocumentStorage
	mimeTypeDetection          pipeline.IMimeTypeDetection
	defaultIndexName           *string
	embeddingGenerationEnabled bool
	cancelFunc                 context.CancelFunc
	ctx                        context.Context
	handlerNames               []string
}

// AddHandler implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) AddHandler(ctx context.Context, handler pipeline.IPipelineStepHandler) error {
	panic("unimplemented")
}

// GetEmbeddingGenerationEnabled implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) GetEmbeddingGenerationEnabled() bool {
	return b.embeddingGenerationEnabled
}

// GetEmbeddingGenerators implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) GetEmbeddingGenerators(ctx context.Context) ([]ai.ITextEmbeddingGenerator, error) {
	return b.embeddingGenerators, nil
}

// GetHandlerNames implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) GetHandlerNames() []string {
	return b.handlerNames
}

// GetMemoryDbs implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) GetMemoryDbs(ctx context.Context) ([]memorystorage.IMemoryDb, error) {
	return b.memoryDbs, nil
}

// GetTextGenerators implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) GetTextGenerators(ctx context.Context) (ai.ITextGenerator, error) {
	return b.textGenerator, nil
}

// ImportDocument implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) ImportDocument(ctx context.Context, index string, uploadRequest *models.DocumentUploadRequest, context kmcontext.IContext) (*string, error) {
	index, err := models.CleanName(&index, b.defaultIndexName)
	if err != nil {
		return nil, err
	}
	contextArgs := map[string]interface{}{}
	if context != nil {
		contextArgs = context.GetArgs()
	}
	dataPipeline, err := b.PrepareNewDocumentUpload(ctx, index, uploadRequest.DocumentId, uploadRequest.Tags, uploadRequest.Files, contextArgs)
	if err != nil {
		return nil, err
	}

	if len(uploadRequest.Steps) > 0 {
		for _, step := range uploadRequest.Steps {
			dataPipeline.Then(step)
		}
	} else {
		for _, step := range b.defaultIngestionSteps {
			dataPipeline.Then(step)
		}
	}

	dataPipeline.Build()
	b.RunPipeline(ctx, dataPipeline)
	return &dataPipeline.DocumentId, nil
}

// IsDocumentReady implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) IsDocumentReady(ctx context.Context, index string, documentId string) (bool, error) {
	index, err := models.CleanName(&index, b.defaultIndexName)
	if err != nil {
		return false, err
	}

	dataPipeline, err := b.ReadPipelineStatus(ctx, index, documentId)
	if err != nil {
		return false, err
	}

	return dataPipeline.Complete() && len(dataPipeline.Files) > 0, nil
}

// ReadFile implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) ReadFile(ctx context.Context, pipeline *pipeline.DataPipeline, fileName string) ([]byte, error) {
	content, err := b.ReadFileAsStream(ctx, pipeline, fileName)
	if err != nil {
		return nil, err
	}
	streamFunc := content.GetStream()
	if streamFunc == nil {
		return nil, fmt.Errorf("no stream function")
	}
	reader, err := streamFunc(ctx)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return data, nil
}

// ReadFileAsStream implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) ReadFileAsStream(ctx context.Context, pipeline *pipeline.DataPipeline, fileName string) (*models.StreamableFileContent, error) {
	index, err := models.CleanName(&pipeline.Index, b.defaultIndexName)
	if err != nil {
		return nil, err
	}
	pipeline.Index = index

	return b.documentStorage.ReadFile(ctx, pipeline.Index, pipeline.DocumentId, fileName, true)
}

// ReadPipelineStatus implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) ReadPipelineStatus(ctx context.Context, index string, documentId string) (*pipeline.DataPipeline, error) {
	index, err := models.CleanName(&index, b.defaultIndexName)
	if err != nil {
		return nil, err
	}
	streamableContent, err := b.documentStorage.ReadFile(ctx, index, documentId, constant.PipelineStatusFilename, false)
	if err != nil {
		return nil, err
	}
	streamFunc := streamableContent.GetStream()
	if streamFunc == nil {
		return nil, fmt.Errorf("no stream function")
	}
	reader, err := streamFunc(ctx)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var result *pipeline.DataPipeline
	err = json.Unmarshal([]byte(strings.TrimSpace(text.RemoveBOM(string(data)))), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadPipelineSummary implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) ReadPipelineSummary(ctx context.Context, index string, documentId string) (*models.DataPipelineStatus, error) {
	index, err := models.CleanName(&index, b.defaultIndexName)
	if err != nil {
		return nil, err
	}
	dataPipeline, err := b.ReadPipelineStatus(ctx, index, documentId)
	if err != nil {
		return nil, err
	}

	return dataPipeline.ToDataPipelineStatus(), nil
}

// ReadTextFile implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) ReadTextFile(ctx context.Context, pipeline *pipeline.DataPipeline, fileName string) (*string, error) {
	index, err := models.CleanName(&pipeline.Index, b.defaultIndexName)
	if err != nil {
		return nil, err
	}
	pipeline.Index = index
	data, err := b.ReadFile(ctx, pipeline, fileName)
	if err != nil {
		return nil, err
	}
	d := string(data)
	return &d, nil
}

// RunPipeline implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) RunPipeline(ctx context.Context, pipeline *pipeline.DataPipeline) error {
	panic("unimplemented")
}

// SetHandlerNames implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) SetHandlerNames(names []string) {
	b.handlerNames = names
}

// StartDocumentDeletion implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) StartDocumentDeletion(ctx context.Context, index *string, documentId string) error {
	ind, err := models.CleanName(index, b.defaultIndexName)
	if err != nil {
		return err
	}
	dataPipeline := PrepareDocumentDeletion(ind, documentId)
	return b.RunPipeline(ctx, dataPipeline)
}

// StartIndexDeletion implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) StartIndexDeletion(ctx context.Context, index *string) error {
	ind, err := models.CleanName(index, b.defaultIndexName)
	if err != nil {
		return err
	}
	dataPipeline := PrepareIndexDeletion(ind)
	return b.RunPipeline(ctx, dataPipeline)
}

// StopAllPipelines implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) StopAllPipelines(ctx context.Context) error {
	if b.cancelFunc != nil {
		b.cancelFunc()
	}
	return nil
}

// TryAddHandler implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) TryAddHandler(ctx context.Context, handler pipeline.IPipelineStepHandler) error {
	panic("unimplemented")
}

// WriteFile implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) WriteFile(ctx context.Context, pipeline *pipeline.DataPipeline, fileName string, fileContent []byte) error {
	index, err := models.CleanName(&pipeline.Index, b.defaultIndexName)
	if err != nil {
		return err
	}
	pipeline.Index = index
	return b.documentStorage.WriteFile(ctx, pipeline.Index, pipeline.DocumentId, fileName, io.NopCloser(bytes.NewReader(fileContent)))
}

// WriteTextFile implements pipeline.IPipelineOrchestrator.
func (b *BaseOrchestrator) WriteTextFile(ctx context.Context, pipeline *pipeline.DataPipeline, fileName string, fileContent string) error {
	index, err := models.CleanName(&pipeline.Index, b.defaultIndexName)
	if err != nil {
		return err
	}
	pipeline.Index = index
	return b.WriteFile(ctx, pipeline, fileName, []byte(fileContent))
}

func (b *BaseOrchestrator) PrepareNewDocumentUpload(ctx context.Context, index string, documentId string, tags *models.TagCollection,
	filesToUpload []models.UploadedFile, contextArgs map[string]interface{}) (*pipeline.DataPipeline, error) {
	index, err := models.CleanName(&index, b.defaultIndexName)
	if err != nil {
		return nil, err
	}

	if filesToUpload == nil {
		filesToUpload = []models.UploadedFile{}
	}

	if contextArgs == nil {
		contextArgs = make(map[string]interface{})
	}
	dataPipeline := &pipeline.DataPipeline{
		Index:            index,
		DocumentId:       documentId,
		Tags:             tags,
		ContextArguments: contextArgs,
		FilesToUpload:    filesToUpload,
	}

	if err = dataPipeline.Validate(); err != nil {
		return nil, err
	}

	return dataPipeline, nil
}

func PrepareIndexDeletion(index string) *pipeline.DataPipeline {
	pipeline := &pipeline.DataPipeline{
		Index:      index,
		DocumentId: "",
	}

	return pipeline.Then(constant.PipelineStepsDeleteIndex).Build()
}

func PrepareDocumentDeletion(index string, documentId string) *pipeline.DataPipeline {
	pipeline := &pipeline.DataPipeline{
		Index:      index,
		DocumentId: documentId,
	}

	return pipeline.Then(constant.PipelineStepsDeleteIndex).Build()
}

func (b *BaseOrchestrator) CleanUpAfterCompletion(ctx context.Context, pipeline *pipeline.DataPipeline) error {
	if pipeline.IsDocumentDeletionPipeline() {
		if err := b.documentStorage.DeleteDocumentDirectory(ctx, pipeline.Index, pipeline.DocumentId); err != nil {
			return err
		}
	}
	if pipeline.IsIndexDeletionPipeline() {
		if err := b.documentStorage.DeleteIndexDirectory(ctx, pipeline.Index); err != nil {
			return err
		}
	}
	return nil
}

func (b *BaseOrchestrator) UpdatePipelineStatus(ctx context.Context, pipeline *pipeline.DataPipeline) error {
	data, err := json.Marshal(pipeline)
	if err != nil {
		return err
	}
	return b.documentStorage.WriteFile(
		ctx,
		pipeline.Index,
		pipeline.DocumentId,
		constant.PipelineStatusFilename,
		io.NopCloser(bytes.NewReader(data)),
	)
}

func (b *BaseOrchestrator) UploadFiles(ctx context.Context, currentPipeline *pipeline.DataPipeline) error {
	if currentPipeline.Complete() {
		return nil
	}
	previousPipeline, err := b.ReadPipelineStatus(ctx, currentPipeline.Index, currentPipeline.DocumentId)
	if err != nil {
		log.Println(err.Error())
	}
	if previousPipeline != nil && previousPipeline.ExecutionId != currentPipeline.ExecutionId {
		var dedupe = map[string]struct{}{}

		for _, oldExecution := range currentPipeline.PreviousExecutionsToPurge {
			dedupe[oldExecution.ExecutionId] = struct{}{}
		}

		for _, oldExecution := range previousPipeline.PreviousExecutionsToPurge {
			if _, ok := dedupe[oldExecution.ExecutionId]; ok {
				continue
			}
			oldExecution.PreviousExecutionsToPurge = []pipeline.DataPipeline{}
			currentPipeline.PreviousExecutionsToPurge = append(currentPipeline.PreviousExecutionsToPurge, oldExecution)
			dedupe[oldExecution.ExecutionId] = struct{}{}
		}

		// Reset the list to avoid wasting space with nested trees
		previousPipeline.PreviousExecutionsToPurge = []pipeline.DataPipeline{}
		currentPipeline.PreviousExecutionsToPurge = append(currentPipeline.PreviousExecutionsToPurge, *previousPipeline)
	}

	return b.uploadFormFiles(ctx, currentPipeline)
}

func (b *BaseOrchestrator) uploadFormFiles(ctx context.Context, dataPipeline *pipeline.DataPipeline) error {
	err := b.documentStorage.CreateIndexDirectory(ctx, dataPipeline.Index)
	if err != nil {
		return err
	}
	err = b.documentStorage.CreateDocumentDirectory(ctx, dataPipeline.Index, dataPipeline.DocumentId)
	if err != nil {
		return err
	}

	for _, file := range dataPipeline.FilesToUpload {
		if file.FileName == constant.PipelineStatusFilename {
			continue
		}
		fileSize, err := io.ReadAll(file.FileContent)
		if err != nil {
			continue
		}
		err = b.documentStorage.WriteFile(ctx, dataPipeline.Index, dataPipeline.DocumentId, file.FileName, file.FileContent)
		if err != nil {
			continue
		}

		mimeType := ""
		if ft, err := b.mimeTypeDetection.GetFileType(file.FileName); err == nil {
			mimeType = *ft
		}

		dataPipeline.Files = append(dataPipeline.Files, pipeline.FileDetails{
			FileDetailsBase: pipeline.FileDetailsBase{
				Id:          uuid.New().String(),
				Name:        file.FileName,
				Size:        int64(len(fileSize)),
				MimeType:    mimeType,
				Tags:        dataPipeline.Tags,
				ProcessedBy: []string{},
				LogEntries:  []pipeline.PipelineLogEntry{},
			},
			GeneratedFiles: map[string]pipeline.GeneratedFileDetails{},
		})
		dataPipeline.LastUpdate = time.Now().UTC()
	}
	return nil
}

func NewBaseOrchestrator(
	documentStorage documentstorage.IDocumentStorage,
	embeddingGenerators []ai.ITextEmbeddingGenerator,
	memoryDbs []memorystorage.IMemoryDb,
	textGenerator ai.ITextGenerator,
	mimeTypeDetection pipeline.IMimeTypeDetection,
	config *configuration.KernelMemoryConfig,
) *BaseOrchestrator {
	if config == nil {
		config = &configuration.KernelMemoryConfig{}
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &BaseOrchestrator{
		memoryDbs:                  memoryDbs,
		embeddingGenerators:        embeddingGenerators,
		textGenerator:              textGenerator,
		defaultIngestionSteps:      config.DataIngestion.GetDefaultStepsOrDefaults(),
		documentStorage:            documentStorage,
		mimeTypeDetection:          mimeTypeDetection,
		defaultIndexName:           &config.DefaultIndexName,
		embeddingGenerationEnabled: config.DataIngestion.EmbeddingGenerationEnabled,
		cancelFunc:                 cancel,
		ctx:                        ctx,
		handlerNames:               []string{},
	}
}
