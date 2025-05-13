package abstractions

import (
	"fmt"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/configuration"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/dataformats"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/prompts"
)

type IKernelMemoryBuilder interface {
	Build(options *KernelMemoryBuilderBuildOptions) (IKernelMemory, error)
	WithoutDefaultHandlers() IKernelMemoryBuilder
	AddIngestionMemoryDb(service memorystorage.IMemoryDb) IKernelMemoryBuilder
	AddIngestionEmbeddingGenerator(service ai.ITextEmbeddingGenerator) IKernelMemoryBuilder
	GetOrchestrator() pipeline.IPipelineOrchestrator

	Configure(action func(IKernelMemoryBuilder)) IKernelMemoryBuilder
	ConfigureWithCondition(condition bool, actionIfTrue func(IKernelMemoryBuilder), actionIfFalse func(IKernelMemoryBuilder)) IKernelMemoryBuilder
	WithCustomIngestionQueueClientFactory(service pipeline.QueueClientFactory) IKernelMemoryBuilder
	WithCustomDocumentStorage(service documentstorage.IDocumentStorage) IKernelMemoryBuilder
	WithCustomMimeTypeDetection(service pipeline.IMimeTypeDetection) IKernelMemoryBuilder
	WithCustomEmbeddingGenerator(service ai.ITextEmbeddingGenerator, useForIngestion bool, useForRetrieval bool) IKernelMemoryBuilder
	WithCustomMemoryDb(service memorystorage.IMemoryDb, useForIngestion bool, useForRetrieval bool) IKernelMemoryBuilder
	WithCustomTextGenerator(service ai.ITextGenerator) IKernelMemoryBuilder
	WithCustomImageOcr(service dataformats.IOcrEngine) IKernelMemoryBuilder
	WithCustomPromptProvider(service prompts.IPromptProvider) IKernelMemoryBuilder
	WithCustomTextPartitioningOptions(options configuration.TextPartitioningOptions) IKernelMemoryBuilder
}

func KernelMemoryGenericsBuild[T IKernelMemory](builder IKernelMemoryBuilder, options *KernelMemoryBuilderBuildOptions) (*T, error) {
	if builder == nil {
		return nil, fmt.Errorf("builder is nil")
	}

	if m, err := builder.Build(options); err != nil {
		return nil, err
	} else {
		if v, ok := m.(T); ok {
			return &v, nil
		} else {
			return nil, fmt.Errorf("builder returned wrong type")
		}
	}
}
