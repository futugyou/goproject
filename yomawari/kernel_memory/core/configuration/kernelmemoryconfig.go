package configuration

import (
	"github.com/futugyou/yomawari/kernel_memory/abstractions/configuration"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/search"
)

const OrchestrationTypeInProcess string = "InProcess"
const OrchestrationTypeDistributed string = "Distributed"

type KernelMemoryConfig struct {
	Service               *ServiceConfig
	DocumentStorageType   string
	TextGeneratorType     string
	ContentModerationType string
	DefaultIndexName      string
	ServiceAuthorization  *ServiceAuthorizationConfig
	DataIngestion         *DataIngestionConfig
	Retrieval             *RetrievalConfig
	Services              map[string]map[string]interface{}
}

type DataIngestionConfig struct {
	OrchestrationType          string
	DistributedOrchestration   *DistributedOrchestrationConfig
	EmbeddingGenerationEnabled bool
	EmbeddingGeneratorTypes    []string
	MemoryDbTypes              []string
	MemoryDbUpsertBatchSize    int
	ImageOcrType               string
	TextPartitioning           *configuration.TextPartitioningOptions
	DefaultSteps               []string
}

func (d *DataIngestionConfig) GetDefaultStepsOrDefaults() []string {
	if d == nil || len(d.DefaultSteps) == 0 {
		return constant.DefaultPipeline
	}
	return d.DefaultSteps
}

type DistributedOrchestrationConfig struct {
	QueueType string
}

type RetrievalConfig struct {
	MemoryDbType           string
	EmbeddingGeneratorType string
	SearchClient           search.SearchClientConfig
}
