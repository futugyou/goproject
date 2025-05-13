package memorystorage

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/core/filesystem"
)

type SimpleTextDb struct {
	fileSystem filesystem.IFileSystem
}

// CreateIndex implements memorystorage.IMemoryDb.
func (s *SimpleTextDb) CreateIndex(ctx context.Context, index string, vectorSize int64) error {
	i, err := NormalizeIndexName(index)
	if err != nil {
		return err
	}
	return s.fileSystem.CreateVolume(ctx, i)
}

// Delete implements memorystorage.IMemoryDb.
func (s *SimpleTextDb) Delete(ctx context.Context, index string, record *memorystorage.MemoryRecord) error {
	index, err := NormalizeIndexName(index)
	if err != nil {
		return err
	}
	return s.fileSystem.DeleteVolume(ctx, index)
}

// DeleteIndex implements memorystorage.IMemoryDb.
func (s *SimpleTextDb) DeleteIndex(ctx context.Context, index string) error {
	i, err := NormalizeIndexName(index)
	if err != nil {
		return err
	}
	return s.fileSystem.DeleteVolume(ctx, i)
}

// GetIndexes implements memorystorage.IMemoryDb.
func (s *SimpleTextDb) GetIndexes(ctx context.Context) ([]string, error) {
	return s.fileSystem.ListVolumes(ctx)
}

// GetList implements memorystorage.IMemoryDb.
func (s *SimpleTextDb) GetList(ctx context.Context, index string, filters []models.MemoryFilter, limit int64, withEmbeddings bool) <-chan memorystorage.MemoryRecordChanResponse {
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
func (s *SimpleTextDb) GetSimilarList(ctx context.Context, index string, text string, filters []models.MemoryFilter, minRelevance float64, limit int64, withEmbeddings bool) <-chan memorystorage.MemoryRecordChanResponse {
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

	re := regexp.MustCompile(`[^a-zA-Z0-9_]+`)
	processedText := re.ReplaceAllString(text, " ")
	words := []string{}
	for _, word := range strings.Fields(processedText) {
		trimmed := strings.TrimSpace(word)
		if trimmed != "" {
			words = append(words, trimmed)
		}
	}

	similarity := map[string]float64{}
	for _, record := range records {
		similarity[record.Id] = 0
		storedText, ok := record.Payload[constant.ReservedPayloadTextField].(string)
		if !ok || len(storedText) == 0 {
			continue
		}
		for _, word := range words {
			if strings.Contains(storedText, word) {
				similarity[record.Id] += 1
			}
		}
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
func (s *SimpleTextDb) Upsert(ctx context.Context, index string, record *memorystorage.MemoryRecord) (*string, error) {
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

var s_replaceIndexNameCharsRegex *regexp.Regexp = regexp.MustCompile(`[\s|\\|/|.|_|:]`)

func NormalizeIndexName(index string) (string, error) {
	if strings.TrimSpace(index) == "" {
		return "", errors.New("the index name is empty")
	}

	index = s_replaceIndexNameCharsRegex.ReplaceAllString(strings.TrimSpace(strings.ToLower(index)), "-")
	return strings.TrimSpace(index), nil
}

func TagsMatchFilters(tags *models.TagCollection, filters []models.MemoryFilter) bool {
	if len(filters) == 0 {
		return true
	}

	for _, filter := range filters {
		match := true

		for key, value := range filter.GetData() {
			tagValues, exists := tags.Get(key)
			if !exists {
				match = false
				break
			}

			found := false
			for _, v := range value {
				if contains(tagValues, v) {
					found = true
					break
				}
			}
			if !found {
				match = false
				break
			}
		}

		if match {
			return true
		}
	}

	return false
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func EncodeId(realId string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(realId))
	return strings.ReplaceAll(encoded, "=", "_")
}

func DecodeId(encodedId string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(encodedId, "_", "="))
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func NewSimpleTextDb(config *SimpleTextDbConfig) *SimpleTextDb {
	var fs filesystem.IFileSystem
	if config == nil {
		config = NewSimpleTextDbConfig()
	}
	if config.StorageType == filesystem.Disk {
		fs = filesystem.NewDiskFileSystem(config.Directory, nil)
	} else {
		fs = filesystem.GetVolatileFileSystemInstance(config.Directory, nil)
	}
	return &SimpleTextDb{
		fileSystem: fs,
	}
}

func filterAndSort(similarity map[string]float64, minRelevance float64) []string {
	var entries []struct {
		Key   string
		Value float64
	}

	for k, v := range similarity {
		if v >= minRelevance {
			entries = append(entries, struct {
				Key   string
				Value float64
			}{k, v})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})

	result := make([]string, len(entries))
	for i, entry := range entries {
		result[i] = entry.Key
	}

	return result
}
