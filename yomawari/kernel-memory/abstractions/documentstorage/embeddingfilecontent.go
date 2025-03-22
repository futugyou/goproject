package documentstorage

import (
	"time"

	"github.com/futugyou/yomawari/kernel-memory/abstractions/ai"
)

type EmbeddingFileContent struct {
	GeneratorName     string       `json:"generator_name"`
	GeneratorProvider string       `json:"generator_provider"`
	VectorSize        int64        `json:"vector_size"`
	SourceFileName    string       `json:"source_file_name"`
	Vector            ai.Embedding `json:"vector"`
	TimeStamp         time.Time    `json:"timestamp"`
}
