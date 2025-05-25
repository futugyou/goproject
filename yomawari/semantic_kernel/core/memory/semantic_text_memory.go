package memory

import (
	"context"

	"github.com/futugyou/yomawari/semantic_kernel/abstractions"
)

type ITextEmbeddingGenerationService interface {
	GenerateEmbedding(ctx context.Context, text string, kernel *abstractions.Kernel) ([]float32, error)
}

var _ abstractions.ISemanticTextMemory = (*SemanticTextMemory)(nil)

type SemanticTextMemory struct {
	embeddingGenerator ITextEmbeddingGenerationService
	storage            abstractions.IMemoryStore
}

func NewSemanticTextMemory(storage abstractions.IMemoryStore, embeddingGenerator ITextEmbeddingGenerationService) *SemanticTextMemory {
	return &SemanticTextMemory{
		storage:            storage,
		embeddingGenerator: embeddingGenerator,
	}
}

func (s *SemanticTextMemory) SaveInformation(
	ctx context.Context,
	collection string,
	text string,
	id string,
	description *string,
	additionalMetadata *string,
	kernel *abstractions.Kernel,
) (string, error) {
	embedding, err := s.embeddingGenerator.GenerateEmbedding(ctx, text, kernel)
	if err != nil {
		return "", err
	}

	data := abstractions.LocalRecord(id, text, description, embedding, additionalMetadata, nil, nil)
	if f, err := s.storage.DoesCollectionExist(ctx, collection); err == nil && !f {
		s.storage.CreateCollection(ctx, collection)
	}

	return s.storage.Upsert(ctx, collection, data)
}

// Get implements memory.ISemanticTextMemory.
func (s *SemanticTextMemory) Get(ctx context.Context, collection string, key string, withEmbedding bool, kernel *abstractions.Kernel) (*abstractions.MemoryQueryResult, error) {
	record, err := s.storage.Get(ctx, collection, key, withEmbedding)
	if err != nil {
		return nil, err
	}
	r := abstractions.FromMemoryRecord(*record, 1)
	return &r, nil
}

// GetCollections implements memory.ISemanticTextMemory.
func (s *SemanticTextMemory) GetCollections(ctx context.Context, kernel *abstractions.Kernel) ([]string, error) {
	resCh, errC := s.storage.GetCollections(ctx)
	result := []string{}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errC:
		return nil, err
	case res := <-resCh:
		result = append(result, res)
	}
	return result, nil
}

// Remove implements memory.ISemanticTextMemory.
func (s *SemanticTextMemory) Remove(ctx context.Context, collection string, key string, kernel *abstractions.Kernel) error {
	return s.storage.Remove(ctx, collection, key)
}

// SaveReference implements memory.ISemanticTextMemory.
func (s *SemanticTextMemory) SaveReference(ctx context.Context, collection string, text string, externalId string, externalSourceName string, description *string, additionalMetadata *string, kernel *abstractions.Kernel) (string, error) {
	embedding, err := s.embeddingGenerator.GenerateEmbedding(ctx, text, kernel)
	if err != nil {
		return "", err
	}

	data := abstractions.ReferenceRecord(externalId, externalSourceName, description, embedding, additionalMetadata, nil, nil)

	if f, err := s.storage.DoesCollectionExist(ctx, collection); err == nil && !f {
		s.storage.CreateCollection(ctx, collection)
	}

	return s.storage.Upsert(ctx, collection, data)
}

// Search implements memory.ISemanticTextMemory.
func (s *SemanticTextMemory) Search(ctx context.Context, collection string, query string, limit int, minRelevanceScore float64, withEmbeddings bool, kernel *abstractions.Kernel) (<-chan abstractions.MemoryQueryResult, <-chan error) {
	queryEmbedding, err := s.embeddingGenerator.GenerateEmbedding(ctx, query, kernel)
	resultCh := make(chan abstractions.MemoryQueryResult)
	errorCh := make(chan error, 1)

	if err != nil {
		go func() {
			defer close(resultCh)
			defer close(errorCh)
			errorCh <- err
		}()
		return resultCh, errorCh
	}

	exists, err := s.storage.DoesCollectionExist(ctx, collection)
	if err != nil || !exists {
		go func() {
			defer close(resultCh)
			defer close(errorCh)
			if err != nil {
				errorCh <- err
			}
		}()
		return resultCh, errorCh
	}

	resCh, errCh := s.storage.GetNearestMatches(ctx, collection, queryEmbedding, limit, minRelevanceScore, withEmbeddings)

	go func() {
		defer close(resultCh)
		defer close(errorCh)

		for {
			select {
			case <-ctx.Done():
				errorCh <- ctx.Err()
				return
			case err, ok := <-errCh:
				if ok {
					errorCh <- err
				}
				return
			case res, ok := <-resCh:
				if ok {
					return
				}
				resultCh <- abstractions.FromMemoryRecord(res.Record, res.Score)
			}
		}
	}()

	return resultCh, errorCh
}
