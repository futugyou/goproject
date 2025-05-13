package handlers

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"

	"golang.org/x/sync/errgroup"
)

type GenerateEmbeddingsParallelHandler struct {
	*GenerateEmbeddingsHandlerBase
	embeddingGenerators        []ai.ITextEmbeddingGenerator
	embeddingGenerationEnabled bool
	stepName                   string
}

func NewGenerateEmbeddingsParallelHandler(stepName string, orchestrator pipeline.IPipelineOrchestrator) *GenerateEmbeddingsParallelHandler {
	list, _ := orchestrator.GetEmbeddingGenerators(context.Background())
	return &GenerateEmbeddingsParallelHandler{
		GenerateEmbeddingsHandlerBase: NewGenerateEmbeddingsHandlerBase(orchestrator),
		embeddingGenerators:           list,
		embeddingGenerationEnabled:    orchestrator.GetEmbeddingGenerationEnabled(),
		stepName:                      stepName,
	}
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (d *GenerateEmbeddingsParallelHandler) GetStepName() string {
	return d.stepName
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (d *GenerateEmbeddingsParallelHandler) SetStepName(name string) {
	d.stepName = name
}

var _ pipeline.IPipelineStepHandler = (*GenerateEmbeddingsParallelHandler)(nil)

// Invoke implements pipeline.IPipelineStepHandler.
func (d *GenerateEmbeddingsParallelHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
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

func (d *GenerateEmbeddingsParallelHandler) GenerateEmbeddingsWithBatching(ctx context.Context, dataPipeline *pipeline.DataPipeline,
	generator ai.ITextEmbeddingBatchGenerator, batchSize int,
	partitions []PartitionInfo,
) error {
	g, ctx := errgroup.WithContext(ctx)
	batches := Chunk(partitions, batchSize)

	for _, batch := range batches {
		batch := batch
		g.Go(func() error {
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
				return err
			}
			return d.SaveEmbeddingsToDocumentStorage(ctx, dataPipeline, partitions, embeddings, GetEmbeddingProviderName(generator), GetEmbeddingGeneratorName(generator))
		})
	}

	return g.Wait()
}

func (d *GenerateEmbeddingsParallelHandler) GenerateEmbeddingsOneAtATime(ctx context.Context, dataPipeline *pipeline.DataPipeline,
	generator ai.ITextEmbeddingGenerator,
	partitions []PartitionInfo,
) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, partition := range partitions {
		partition := partition
		g.Go(func() error {
			embedding, err := generator.GenerateEmbedding(ctx, partition.PartitionContent)
			if err != nil {
				return err
			}
			return d.SaveEmbeddingToDocumentStorage(ctx, dataPipeline, partition, embedding, GetEmbeddingProviderName(generator), GetEmbeddingGeneratorName(generator))
		})
	}
	return g.Wait()
}
