package abstractions

import (
	"context"
)

type ISemanticTextMemory interface {
	SaveInformation(ctx context.Context, collection string, text string, id string, description *string, additionalMetadata *string, kernel *Kernel) (string, error)
	SaveReference(ctx context.Context, collection string, text string, externalId string, externalSourceName string, description *string, additionalMetadata *string, kernel *Kernel) (string, error)
	Get(ctx context.Context, collection string, key string, withEmbedding bool, kernel *Kernel) (*MemoryQueryResult, error)
	Remove(ctx context.Context, collection string, key string, kernel *Kernel) error
	// limit = 1, minRelevanceScore = 0.7,withEmbeddings = false,
	Search(ctx context.Context, collection string, query string, limit int, minRelevanceScore float64, withEmbeddings bool, kernel *Kernel) (<-chan MemoryQueryResult, <-chan error)
	GetCollections(ctx context.Context, kernel *Kernel) ([]string, error)
}

var _ ISemanticTextMemory = (*NullMemory)(nil)

type NullMemory struct {
}

// Get implements ISemanticTextMemory.
func (n *NullMemory) Get(ctx context.Context, collection string, key string, withEmbedding bool, kernel *Kernel) (*MemoryQueryResult, error) {
	return nil, nil
}

// GetCollections implements ISemanticTextMemory.
func (n *NullMemory) GetCollections(ctx context.Context, kernel *Kernel) ([]string, error) {
	return nil, nil
}

// Remove implements ISemanticTextMemory.
func (n *NullMemory) Remove(ctx context.Context, collection string, key string, kernel *Kernel) error {
	return nil
}

// SaveInformation implements ISemanticTextMemory.
func (n *NullMemory) SaveInformation(ctx context.Context, collection string, text string, id string, description *string, additionalMetadata *string, kernel *Kernel) (string, error) {
	return "", nil
}

// SaveReference implements ISemanticTextMemory.
func (n *NullMemory) SaveReference(ctx context.Context, collection string, text string, externalId string, externalSourceName string, description *string, additionalMetadata *string, kernel *Kernel) (string, error) {
	return "", nil
}

// Search implements ISemanticTextMemory.
func (n *NullMemory) Search(ctx context.Context, collection string, query string, limit int, minRelevanceScore float64, withEmbeddings bool, kernel *Kernel) (<-chan MemoryQueryResult, <-chan error) {
	results := make(chan MemoryQueryResult)
	errs := make(chan error)

	close(results)
	close(errs)

	return results, errs
}
