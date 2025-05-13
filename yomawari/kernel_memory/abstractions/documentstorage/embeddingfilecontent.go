package documentstorage

import (
	"encoding/json"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
)

type EmbeddingFileContent struct {
	GeneratorName     string       `json:"generator_name"`
	GeneratorProvider string       `json:"generator_provider"`
	VectorSize        int64        `json:"vector_size"`
	SourceFileName    string       `json:"source_file_name"`
	Vector            ai.Embedding `json:"vector"`
	TimeStamp         time.Time    `json:"timestamp"`
}

func (e *EmbeddingFileContent) ToJson() string {
	if e == nil {
		return "{}"
	}

	datas, err := json.Marshal(e)
	if err != nil {
		return "{}"
	}
	return string(datas)
}
