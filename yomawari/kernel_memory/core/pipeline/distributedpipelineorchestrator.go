package pipeline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/core/configuration"
)

type DistributedPipelineOrchestrator struct {
	*BaseOrchestrator
	queueClientFactory *pipeline.QueueClientFactory
	queues             map[string]pipeline.IQueue
}

func NewDistributedPipelineOrchestrator(
	queueClientFactory *pipeline.QueueClientFactory,
	documentStorage documentstorage.IDocumentStorage,
	embeddingGenerators []ai.ITextEmbeddingGenerator,
	memoryDbs []memorystorage.IMemoryDb,
	textGenerator ai.ITextGenerator,
	mimeTypeDetection pipeline.IMimeTypeDetection,
	config *configuration.KernelMemoryConfig) *DistributedPipelineOrchestrator {
	return &DistributedPipelineOrchestrator{
		BaseOrchestrator:   NewBaseOrchestrator(documentStorage, embeddingGenerators, memoryDbs, textGenerator, mimeTypeDetection, config),
		queueClientFactory: queueClientFactory,
		queues:             map[string]pipeline.IQueue{},
	}
}

func (b *DistributedPipelineOrchestrator) GetHandlerNames() []string {
	keys := make([]string, 0, len(b.queues))
	for k := range b.queues {
		keys = append(keys, k)
	}

	return keys
}

func (b *DistributedPipelineOrchestrator) SetHandlerNames(names []string) {
}

func (b *DistributedPipelineOrchestrator) runPipelineStep(ctx context.Context, dataPipeline *pipeline.DataPipeline, handler pipeline.IPipelineStepHandler) pipeline.ReturnType {
	if dataPipeline.Complete() {
		return pipeline.ReturnTypeSuccess
	}
	returnType, updatedPipeline, err := handler.Invoke(ctx, dataPipeline)
	if err != nil {
		return pipeline.ReturnTypeFatalError
	}

	switch returnType {
	case pipeline.ReturnTypeSuccess:
		dataPipeline = updatedPipeline
		dataPipeline.LastUpdate = time.Now().UTC()
		dataPipeline.MoveToNextStep()
		b.moveForward(ctx, dataPipeline)
	}

	return returnType
}

func (b *DistributedPipelineOrchestrator) moveForward(ctx context.Context, dataPipeline *pipeline.DataPipeline) (err error) {
	if dataPipeline.Complete() {
		err = b.UpdatePipelineStatus(ctx, dataPipeline)
		if err != nil {
			return
		}
		err = b.CleanUpAfterCompletion(ctx, dataPipeline)
	} else {
		if len(dataPipeline.RemainingSteps) > 0 {
			var queue pipeline.IQueue
			nextStepName := dataPipeline.RemainingSteps[0]
			queue, err = b.queueClientFactory.Build(ctx)
			if err != nil {
				return
			}
			queue, err = queue.ConnectToQueue(ctx, nextStepName, pipeline.PublishOnly)
			if err != nil {
				return
			}
			err = b.UpdatePipelineStatus(ctx, dataPipeline)
			if err != nil {
				return
			}
			pointer := pipeline.NewDataPipelinePointer(*dataPipeline)
			var data []byte
			data, err = json.Marshal(pointer)
			if err != nil {
				return
			}
			err = queue.Enqueue(ctx, string(data))
		}
	}

	return
}

func (b *DistributedPipelineOrchestrator) RunPipeline(ctx context.Context, dataPipeline *pipeline.DataPipeline) error {
	if err := b.UploadFiles(ctx, dataPipeline); err != nil {
		return err
	}
	if dataPipeline.Complete() {
		return nil
	}
	return b.moveForward(ctx, dataPipeline)
}

func (b *DistributedPipelineOrchestrator) TryAddHandler(ctx context.Context, handler pipeline.IPipelineStepHandler) error {
	if _, ok := b.queues[handler.GetStepName()]; ok {
		return nil
	}
	return b.AddHandler(ctx, handler)
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func (b *DistributedPipelineOrchestrator) AddHandler(ctx context.Context, handler pipeline.IPipelineStepHandler) error {
	if _, ok := b.queues[handler.GetStepName()]; ok {
		return fmt.Errorf("handler with name %s already exists", handler.GetStepName())
	}
	q, err := b.queueClientFactory.Build(ctx)
	if err != nil {
		return err
	}
	err = q.OnDequeue(ctx, func(ctx context.Context, message string) pipeline.ReturnType {
		var pipelinePointer pipeline.DataPipelinePointer
		err := json.Unmarshal([]byte(message), &pipelinePointer)
		if err != nil {
			return pipeline.ReturnTypeFatalError
		}
		dataPipeline, err := b.ReadPipelineStatus(ctx, pipelinePointer.Index, pipelinePointer.DocumentId)
		if err != nil {
			deletingIndex := handler.GetStepName() == constant.PipelineStepsDeleteIndex && contains(pipelinePointer.Steps, constant.PipelineStepsDeleteIndex)
			if deletingIndex {
				dataPipeline = &pipeline.DataPipeline{
					Index:                     pipelinePointer.Index,
					DocumentId:                pipelinePointer.DocumentId,
					ExecutionId:               pipelinePointer.ExecutionId,
					Steps:                     pipelinePointer.Steps,
					RemainingSteps:            []string{},
					CompletedSteps:            []string{},
					Tags:                      &models.TagCollection{},
					Creation:                  time.Time{},
					LastUpdate:                time.Time{},
					Files:                     []pipeline.FileDetails{},
					ContextArguments:          map[string]any{},
					PreviousExecutionsToPurge: []pipeline.DataPipeline{},
					FilesToUpload:             []models.UploadedFile{},
					UploadComplete:            false,
				}
				return b.runPipelineStep(b.ctx, dataPipeline, handler)
			}
		}
		if dataPipeline == nil {
			return pipeline.ReturnTypeTransientError
		}
		if pipelinePointer.ExecutionId != dataPipeline.ExecutionId {
			return pipeline.ReturnTypeSuccess
		}
		currentStepName := ""
		if len(dataPipeline.RemainingSteps) > 0 {
			currentStepName = dataPipeline.RemainingSteps[0]
		}
		if currentStepName != handler.GetStepName() {
			dataPipeline.RollbackToPreviousStep()
			b.UpdatePipelineStatus(ctx, dataPipeline)
		}
		return b.runPipelineStep(b.ctx, dataPipeline, handler)
	})
	if err != nil {
		return err
	}
	q, err = q.ConnectToQueue(ctx, handler.GetStepName(), pipeline.PubSub)
	if err != nil {
		return err
	}
	b.queues[handler.GetStepName()] = q
	return err
}
