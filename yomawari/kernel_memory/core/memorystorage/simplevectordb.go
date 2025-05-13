package memorystorage

import (
	"context"
	"encoding/json"
	"math"
	"strings"
	"sync"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/core/filesystem"
)

type SimpleVectorDb struct {
	fileSystem         filesystem.IFileSystem
	embeddingGenerator ai.ITextEmbeddingGenerator
}

func NewSimpleVectorDb(embeddingGenerator ai.ITextEmbeddingGenerator, config *SimpleVectorDbConfig) *SimpleVectorDb {
	var fs filesystem.IFileSystem
	if config == nil {
		config = NewSimpleVectorDbConfig()
	}
	if config.StorageType == filesystem.Disk {
		fs = filesystem.NewDiskFileSystem(config.Directory, nil)
	} else {
		fs = filesystem.GetVolatileFileSystemInstance(config.Directory, nil)
	}
	return &SimpleVectorDb{
		embeddingGenerator: embeddingGenerator,
		fileSystem:         fs,
	}
}

// CreateIndex implements memorystorage.IMemoryDb.
func (s *SimpleVectorDb) CreateIndex(ctx context.Context, index string, vectorSize int64) error {
	i, err := NormalizeIndexName(index)
	if err != nil {
		return err
	}
	return s.fileSystem.CreateVolume(ctx, i)
}

// Delete implements memorystorage.IMemoryDb.
func (s *SimpleVectorDb) Delete(ctx context.Context, index string, record *memorystorage.MemoryRecord) error {
	index, err := NormalizeIndexName(index)
	if err != nil {
		return err
	}
	return s.fileSystem.DeleteVolume(ctx, index)
}

// DeleteIndex implements memorystorage.IMemoryDb.
func (s *SimpleVectorDb) DeleteIndex(ctx context.Context, index string) error {
	i, err := NormalizeIndexName(index)
	if err != nil {
		return err
	}
	return s.fileSystem.DeleteVolume(ctx, i)
}

// GetIndexes implements memorystorage.IMemoryDb.
func (s *SimpleVectorDb) GetIndexes(ctx context.Context) ([]string, error) {
	return s.fileSystem.ListVolumes(ctx)
}

// GetList implements memorystorage.IMemoryDb.
func (s *SimpleVectorDb) GetList(ctx context.Context, index string, filters []models.MemoryFilter, limit int64, withEmbeddings bool) <-chan memorystorage.MemoryRecordChanResponse {
	recordCh := make(chan memorystorage.MemoryRecordChanResponse)
	if limit <= 0 {
		limit = math.MaxInt64
	}

	index, err := NormalizeIndexName(index)
	if err != nil {
		recordCh <- memorystorage.MemoryRecordChanResponse{
			Err: err,
		}
		close(recordCh)
		return recordCh
	}

	filts := []models.MemoryFilter{}
	for _, filter := range filters {
		if filter.Count() > 0 {
			filts = append(filts, filter)
		}
	}
	list, err := s.fileSystem.ReadAllFilesAsText(ctx, index, "")
	if err != nil {
		list = map[string]string{}
	}
	go func() {
		defer close(recordCh)
		for _, v := range list {
			var record memorystorage.MemoryRecord
			if err := json.Unmarshal([]byte(v), &record); err != nil {
				continue
			}
			if TagsMatchFilters(record.Tags, filts) {
				limit--
				if limit <= 0 {
					return
				}
				select {
				case recordCh <- memorystorage.MemoryRecordChanResponse{
					Record: &record,
					Err:    nil,
				}:
				default:
				}
			}

		}
	}()

	return recordCh
}

// GetSimilarList implements memorystorage.IMemoryDb.
func (s *SimpleVectorDb) GetSimilarList(ctx context.Context, index string, text string, filters []models.MemoryFilter, minRelevance float64, limit int64, withEmbeddings bool) <-chan memorystorage.MemoryRecordChanResponse {
	recordCh := make(chan memorystorage.MemoryRecordChanResponse)
	if limit <= 0 {
		limit = math.MaxInt64
	}

	index, err := NormalizeIndexName(index)
	if err != nil {
		recordCh <- memorystorage.MemoryRecordChanResponse{
			Err: err,
		}
		close(recordCh)
		return recordCh
	}

	listCh := s.GetList(ctx, index, filters, limit, withEmbeddings)
	records := map[string]*memorystorage.MemoryRecord{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for response := range listCh {
			if response.Err == nil {
				records[response.Record.Id] = response.Record
			}
		}
	}()
	wg.Wait()

	similarity := map[string]float64{}
	textEmbedding, err := s.embeddingGenerator.GenerateEmbedding(ctx, text)
	if err != nil {
		recordCh <- memorystorage.MemoryRecordChanResponse{
			Err: err,
		}
		close(recordCh)
		return recordCh
	}
	for _, record := range records {
		s, err := textEmbedding.CosineSimilarity(*record.Vector)
		if err != nil {
			continue
		}
		similarity[record.Id] = s
	}

	sorted := filterAndSort(similarity, minRelevance)

	go func() {
		defer close(recordCh)
		for _, id := range sorted {
			s := similarity[id]
			recordCh <- memorystorage.MemoryRecordChanResponse{
				Record:   records[id],
				Similars: &s,
				Err:      nil,
			}
		}
	}()

	return recordCh
}

// Upsert implements memorystorage.IMemoryDb.
func (s *SimpleVectorDb) Upsert(ctx context.Context, index string, record *memorystorage.MemoryRecord) (*string, error) {
	i, err := NormalizeIndexName(index)
	if err != nil {
		return nil, err
	}
	b, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(b))
	err = s.fileSystem.WriteFile(ctx, i, "", EncodeId(record.Id), reader)
	if err != nil {
		return nil, err
	}
	return &record.Id, nil
}
