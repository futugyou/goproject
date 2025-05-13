package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/core/configuration"
)

type InProcessPipelineOrchestrator struct {
	*BaseOrchestrator
	handlers map[string]pipeline.IPipelineStepHandler
}

func NewInProcessPipelineOrchestrator(
	documentStorage documentstorage.IDocumentStorage,
	embeddingGenerators []ai.ITextEmbeddingGenerator,
	memoryDbs []memorystorage.IMemoryDb,
	textGenerator ai.ITextGenerator,
	mimeTypeDetection pipeline.IMimeTypeDetection,
	config *configuration.KernelMemoryConfig) *InProcessPipelineOrchestrator {
	return &InProcessPipelineOrchestrator{
		BaseOrchestrator: NewBaseOrchestrator(documentStorage, embeddingGenerators, memoryDbs, textGenerator, mimeTypeDetection, config),
	}
}

func (b *InProcessPipelineOrchestrator) GetHandlerNames() []string {
	keys := make([]string, 0, len(b.handlers))
	for k := range b.handlers {
		keys = append(keys, k)
	}

	return keys
}

func (b *InProcessPipelineOrchestrator) SetHandlerNames(names []string) {
}

func (b *InProcessPipelineOrchestrator) RunPipeline(ctx context.Context, dataPipeline *pipeline.DataPipeline) error {
	if err := b.UploadFiles(ctx, dataPipeline); err != nil {
		return err
	}
	if err := b.UpdatePipelineStatus(ctx, dataPipeline); err != nil {
		return err
	}

	for {
		if dataPipeline.Complete() {
			break
		}
		currentStepName := ""
		if len(dataPipeline.RemainingSteps) > 0 {
			currentStepName = dataPipeline.RemainingSteps[0]
		}

		stepHandler, ok := b.handlers[currentStepName]
		if !ok {
			return fmt.Errorf("no handler found for step '%s'", currentStepName)
		}

		returnType, updatedPipeline, err := stepHandler.Invoke(b.ctx, dataPipeline)
		if err != nil {
			return err
		}

		switch returnType {
		case pipeline.ReturnTypeSuccess:
			dataPipeline = updatedPipeline
			dataPipeline.LastUpdate = time.Now().UTC()
			dataPipeline.MoveToNextStep()
			b.UpdatePipelineStatus(ctx, dataPipeline)
		}
	}

	return b.CleanUpAfterCompletion(ctx, dataPipeline)
}

func (b *InProcessPipelineOrchestrator) TryAddHandler(ctx context.Context, handler pipeline.IPipelineStepHandler) error {
	if _, ok := b.handlers[handler.GetStepName()]; ok {
		return nil
	}
	return b.AddHandler(ctx, handler)
}

func (b *InProcessPipelineOrchestrator) AddHandler(ctx context.Context, handler pipeline.IPipelineStepHandler) error {
	if _, ok := b.handlers[handler.GetStepName()]; ok {
		return fmt.Errorf("handler with name %s already exists", handler.GetStepName())
	}

	b.handlers[handler.GetStepName()] = handler
	return nil
}
