package abstractions

import "context"

type IMemoryStore interface {
	CreateCollection(ctx context.Context, collectionName string) error
	GetCollections(ctx context.Context) (<-chan string, <-chan error)
	DoesCollectionExist(ctx context.Context, collectionName string) (bool, error)
	DeleteCollection(ctx context.Context, collectionName string) error
	Upsert(ctx context.Context, collectionName string, record MemoryRecord) (string, error)
	UpsertBatch(ctx context.Context, collectionName string, records []MemoryRecord) (<-chan string, <-chan error)
	Get(ctx context.Context, collectionName string, key string, withEmbedding bool) (*MemoryRecord, error)
	GetBatch(ctx context.Context, collectionName string, keys []string, withEmbeddings bool) (<-chan MemoryRecord, <-chan error)
	Remove(ctx context.Context, collectionName string, key string) error
	RemoveBatch(ctx context.Context, collectionName string, keys []string) error
	GetNearestMatches(ctx context.Context, collectionName string, embedding []float32, limit int, minRelevanceScore float64, withEmbeddings bool) (<-chan NearestMatchesResult, <-chan error)
	GetNearestMatch(ctx context.Context, collectionName string, embedding []float32, limit int, minRelevanceScore float64, withEmbeddings bool) (*MemoryRecord, float64, error)
}

type NearestMatchesResult struct {
	Record MemoryRecord
	Score  float64
}
