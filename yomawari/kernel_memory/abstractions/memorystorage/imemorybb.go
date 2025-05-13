package memorystorage

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type IMemoryDb interface {
	CreateIndex(ctx context.Context, index string, vectorSize int64) error
	GetIndexes(ctx context.Context) ([]string, error)
	DeleteIndex(ctx context.Context, index string) error
	Upsert(ctx context.Context, index string, record *MemoryRecord) (*string, error)
	GetSimilarList(ctx context.Context, index string, text string, filters []models.MemoryFilter, minRelevance float64, limit int64, withEmbeddings bool) <-chan MemoryRecordChanResponse
	GetList(ctx context.Context, index string, filters []models.MemoryFilter, limit int64, withEmbeddings bool) <-chan MemoryRecordChanResponse
	Delete(ctx context.Context, index string, record *MemoryRecord) error
}

type IMemoryDbUpsertBatch interface {
	UpsertBatch(ctx context.Context, index string, records []*MemoryRecord) <-chan MemoryUpsertBatchChanResponse
}

type MemoryUpsertBatchChanResponse struct {
	RecordID *string
	Err      error
}

type MemoryRecordChanResponse struct {
	Record   *MemoryRecord
	Similars *float64
	Err      error
}
