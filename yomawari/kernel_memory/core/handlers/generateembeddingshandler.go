package handlers

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
)

type GenerateEmbeddingsHandler struct {
	*GenerateEmbeddingsHandlerBase
	embeddingGenerators        []ai.ITextEmbeddingGenerator
	embeddingGenerationEnabled bool
	stepName                   string
}

func NewGenerateEmbeddingsHandler(stepName string, orchestrator pipeline.IPipelineOrchestrator) *GenerateEmbeddingsHandler {
	list, _ := orchestrator.GetEmbeddingGenerators(context.Background())
	return &GenerateEmbeddingsHandler{
		GenerateEmbeddingsHandlerBase: NewGenerateEmbeddingsHandlerBase(orchestrator),
		embeddingGenerators:           list,
		embeddingGenerationEnabled:    orchestrator.GetEmbeddingGenerationEnabled(),
		stepName:                      stepName,
	}
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (d *GenerateEmbeddingsHandler) GetStepName() string {
	return d.stepName
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (d *GenerateEmbeddingsHandler) SetStepName(name string) {
	d.stepName = name
}

var _ pipeline.IPipelineStepHandler = (*GenerateEmbeddingsHandler)(nil)

// Invoke implements pipeline.IPipelineStepHandler.
func (d *GenerateEmbeddingsHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	if !d.embeddingGenerationEnabled {
		return pipeline.ReturnTypeSuccess, dataPipeline, nil
	}

	for _, generator := range d.embeddingGenerators {
		var subStepName = GetSubStepName(generator)
		partitions, err := d.GetListOfPartitionsToProcess(ctx, dataPipeline, subStepName)
		if err != nil {
			return pipeline.ReturnTypeFatalError, dataPipeline, err
		}
		if batchGenerator, ok := generator.(ai.ITextEmbeddingBatchGenerator); ok {
			maxBatchSize := batchGenerator.GetMaxBatchSize()
			err = d.GenerateEmbeddingsWithBatching(ctx, dataPipeline, batchGenerator, int(maxBatchSize), partitions)
		} else {
			err = d.GenerateEmbeddingsOneAtATime(ctx, dataPipeline, generator, partitions)
		}
		if err != nil {
			return pipeline.ReturnTypeFatalError, dataPipeline, err
		}
	}
	return pipeline.ReturnTypeSuccess, dataPipeline, nil
}

func (d *GenerateEmbeddingsHandler) GenerateEmbeddingsWithBatching(ctx context.Context, dataPipeline *pipeline.DataPipeline,
	generator ai.ITextEmbeddingBatchGenerator, batchSize int,
	partitions []PartitionInfo,
) error {
	batches := Chunk(partitions, batchSize)
	for _, batch := range batches {
		strings := []string{}
		var totalTokens int64 = 0
		for _, v := range batch {
			strings = append(strings, v.PartitionContent)
			if generator, ok := generator.(ai.ITextEmbeddingGenerator); ok {
				totalTokens = totalTokens + generator.CountTokens(ctx, v.PartitionContent)
			}
		}
		embeddings, err := generator.GenerateEmbeddingBatch(ctx, strings)
		if err != nil {
			continue
		}
		d.SaveEmbeddingsToDocumentStorage(ctx, dataPipeline, partitions, embeddings, GetEmbeddingProviderName(generator), GetEmbeddingGeneratorName(generator))
	}
	return nil
}

func (d *GenerateEmbeddingsHandler) GenerateEmbeddingsOneAtATime(ctx context.Context, dataPipeline *pipeline.DataPipeline,
	generator ai.ITextEmbeddingGenerator,
	partitions []PartitionInfo,
) error {
	for _, partition := range partitions {
		embedding, err := generator.GenerateEmbedding(ctx, partition.PartitionContent)
		if err != nil {
			continue
		}
		d.SaveEmbeddingToDocumentStorage(ctx, dataPipeline, partition, embedding, GetEmbeddingProviderName(generator), GetEmbeddingGeneratorName(generator))
	}
	return nil
}

func Chunk[T any](data []T, size int) [][]T {
	if size <= 0 {
		panic("size must be greater than 0")
	}

	var result [][]T
	for i := 0; i < len(data); i += size {
		end := i + size
		if end > len(data) {
			end = len(data)
		}
		result = append(result, data[i:end])
	}
	return result
}
