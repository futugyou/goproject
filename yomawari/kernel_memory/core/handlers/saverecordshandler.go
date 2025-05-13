package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
	"github.com/futugyou/yomawari/kernel_memory/core/configuration"
)

var _ pipeline.IPipelineStepHandler = (*SaveRecordsHandler)(nil)

type SaveRecordsHandler struct {
	orchestrator               pipeline.IPipelineOrchestrator
	memoryDbs                  []memorystorage.IMemoryDb
	memoryDbsWithSingleUpsert  []memorystorage.IMemoryDb
	memoryDbsWithBatchUpsert   []memorystorage.IMemoryDb
	embeddingGenerationEnabled bool
	upsertBatchSize            int
	usingBatchUpsert           bool
	stepName                   string
}

func NewSaveRecordsHandler(stepName string, orchestrator pipeline.IPipelineOrchestrator, config *configuration.KernelMemoryConfig) *SaveRecordsHandler {
	upsertBatchSize := 0
	if config != nil && config.DataIngestion != nil {
		upsertBatchSize = config.DataIngestion.MemoryDbUpsertBatchSize
	}
	embeddingGenerationEnabled := orchestrator.GetEmbeddingGenerationEnabled()
	memoryDbs, _ := orchestrator.GetMemoryDbs(context.Background())

	handler := &SaveRecordsHandler{
		orchestrator:               orchestrator,
		memoryDbs:                  memoryDbs,
		memoryDbsWithSingleUpsert:  memoryDbs,
		memoryDbsWithBatchUpsert:   []memorystorage.IMemoryDb{},
		embeddingGenerationEnabled: embeddingGenerationEnabled,
		upsertBatchSize:            upsertBatchSize,
		usingBatchUpsert:           false,
		stepName:                   stepName,
	}
	if upsertBatchSize > 1 {
		memoryDbsWithSingleUpsert := []memorystorage.IMemoryDb{}
		memoryDbsWithBatchUpsert := []memorystorage.IMemoryDb{}
		for _, db := range memoryDbs {
			if _, ok := db.(memorystorage.IMemoryDbUpsertBatch); ok {
				memoryDbsWithBatchUpsert = append(memoryDbsWithBatchUpsert, db)
			} else {
				memoryDbsWithSingleUpsert = append(memoryDbsWithSingleUpsert, db)
			}
		}
		handler.memoryDbsWithSingleUpsert = memoryDbsWithSingleUpsert
		handler.memoryDbsWithBatchUpsert = memoryDbsWithBatchUpsert
		handler.usingBatchUpsert = len(memoryDbsWithBatchUpsert) > 0
	}
	return handler
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (s *SaveRecordsHandler) GetStepName() string {
	return s.stepName
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (s *SaveRecordsHandler) SetStepName(name string) {
	s.stepName = name
}

// Invoke implements pipeline.IPipelineStepHandler.
func (s *SaveRecordsHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	s.DeletePreviousRecords(ctx, dataPipeline)
	dataPipeline.PreviousExecutionsToPurge = []pipeline.DataPipeline{}

	var recordsFound = false
	createdIndexes := map[string]struct{}{}

	var sourceFiles [][]FileDetailsWithRecordId
	if s.embeddingGenerationEnabled {
		sourceFiles = Chunk(GetListOfEmbeddingFiles(dataPipeline), s.upsertBatchSize)
	} else {
		sourceFiles = Chunk(GetListOfPartitionAndSyntheticFiles(dataPipeline), s.upsertBatchSize)
	}

	for i := range sourceFiles {
		records := []*memorystorage.MemoryRecord{}
		for ii := range sourceFiles[i] {
			file := &sourceFiles[i][ii]
			if file.File == nil {
				continue
			}
			if file.File.AlreadyProcessedBy(s, nil) {
				recordsFound = true
				continue
			}

			var record *memorystorage.MemoryRecord
			fileDetails := dataPipeline.GetFile(file.File.ParentId)
			webPageUrl := s.GetSourceUrl(ctx, dataPipeline, fileDetails)
			if s.embeddingGenerationEnabled {
				recordsFound = true
				vectorJson, err := s.orchestrator.ReadTextFile(ctx, dataPipeline, file.File.Name)
				if err != nil {
					return pipeline.ReturnTypeFatalError, dataPipeline, err
				}
				var embeddingData documentstorage.EmbeddingFileContent

				err = json.Unmarshal([]byte(strings.TrimSpace(text.RemoveBOM(*vectorJson))), &embeddingData)
				if err != nil {
					return pipeline.ReturnTypeFatalError, dataPipeline, err
				}
				partitionContent, err := s.orchestrator.ReadTextFile(ctx, dataPipeline, embeddingData.SourceFileName)
				if err != nil {
					return pipeline.ReturnTypeFatalError, dataPipeline, err
				}

				record = PrepareRecord(
					dataPipeline,
					file.RecordId,
					fileDetails.Name,
					webPageUrl,
					file.File.ParentId,
					file.File.SourcePartitionId,
					*partitionContent,
					file.File.PartitionNumber,
					file.File.SectionNumber,
					&embeddingData.Vector,
					embeddingData.GeneratorProvider,
					embeddingData.GeneratorName,
					file.File.Tags)
			} else {
				switch file.File.MimeType {
				case pipeline.MimeTypes_PlainText, pipeline.MimeTypes_MarkDown:
					recordsFound = true
					partitionContent, err := s.orchestrator.ReadTextFile(ctx, dataPipeline, file.File.Name)
					if err != nil {
						return pipeline.ReturnTypeFatalError, dataPipeline, err
					}
					record = PrepareRecord(
						dataPipeline,
						file.RecordId,
						fileDetails.Name,
						webPageUrl,
						file.File.ParentId,
						file.File.Id,
						*partitionContent,
						fileDetails.PartitionNumber,
						fileDetails.SectionNumber,
						nil,
						"",
						"",
						file.File.Tags)
				default:
					continue
				}
			}
			if record == nil {
				continue
			}
			records = append(records, record)
			for _, db := range s.memoryDbsWithSingleUpsert {
				CreateIndexOnce(ctx, db, createdIndexes, dataPipeline.Index, record.Vector.Length(), false)
				s.SaveRecord(ctx, dataPipeline, db, record, createdIndexes)
			}
			if !s.usingBatchUpsert {
				file.File.MarkProcessedBy(s, nil)
			}
		}

		if s.usingBatchUpsert {
			if len(records) > 0 {
				for _, db := range s.memoryDbsWithBatchUpsert {
					CreateIndexOnce(ctx, db, createdIndexes, dataPipeline.Index, records[0].Vector.Length(), false)
					s.SaveRecordsBatch(ctx, dataPipeline, db, records, createdIndexes)
				}
			}

			for _, file := range sourceFiles[i] {
				file.File.MarkProcessedBy(s, nil)
			}
		}
	}

	if !recordsFound {
		log.Printf("Pipeline '%s/%s': step %s: no records found, cannot save, moving to next pipeline step.", dataPipeline.Index, dataPipeline.DocumentId, s.GetStepName())
	}
	return pipeline.ReturnTypeSuccess, dataPipeline, nil
}

func (s *SaveRecordsHandler) SaveRecord(ctx context.Context, pipe *pipeline.DataPipeline, db memorystorage.IMemoryDb, record *memorystorage.MemoryRecord, createdIndexes map[string]struct{}) {
	_, err := db.Upsert(ctx, pipe.Index, record)
	if err != nil {
		CreateIndexOnce(ctx, db, createdIndexes, pipe.Index, record.Vector.Length(), true)
		db.Upsert(ctx, pipe.Index, record)
	}
}

func (s *SaveRecordsHandler) SaveRecordsBatch(ctx context.Context, pipe *pipeline.DataPipeline, db memorystorage.IMemoryDb, records []*memorystorage.MemoryRecord, createdIndexes map[string]struct{}) {
	if dbBatch, ok := db.(memorystorage.IMemoryDbUpsertBatch); ok {
		streamResp := dbBatch.UpsertBatch(ctx, pipe.Index, records)
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			for response := range streamResp {
				if response.Err != nil {
					CreateIndexOnce(ctx, db, createdIndexes, pipe.Index, records[0].Vector.Length(), true)
					retryStream := dbBatch.UpsertBatch(ctx, pipe.Index, records)
					for retryResp := range retryStream {
						if retryResp.Err != nil {
							log.Printf("Retry UpsertBatch failed: %v", retryResp.Err)
						}
					}
				}
			}
		}()
		wg.Wait()
	}
}

func (s *SaveRecordsHandler) DeletePreviousRecords(ctx context.Context, pipe *pipeline.DataPipeline) {
	if len(pipe.PreviousExecutionsToPurge) == 0 {
		return
	}

	var recordsToKeep = map[string]struct{}{}

	// Decide which records not to delete, looking at the current pipeline
	ps := append(GetListOfEmbeddingFiles(pipe), GetListOfEmbeddingFiles(pipe)...)

	for _, embeddingFile := range ps {
		recordsToKeep[embeddingFile.RecordId] = struct{}{}
	}
	for _, oldPipeline := range pipe.PreviousExecutionsToPurge {
		for _, file := range append(GetListOfEmbeddingFiles(&oldPipeline), GetListOfPartitionAndSyntheticFiles(&oldPipeline)...) {
			if _, ok := recordsToKeep[file.RecordId]; ok {
				continue
			}
			for _, client := range s.memoryDbs {
				m := &memorystorage.MemoryRecord{Id: file.RecordId}
				client.Delete(ctx, pipe.Index, m)
			}
		}
	}
}

func GetListOfEmbeddingFiles(pipe *pipeline.DataPipeline) []FileDetailsWithRecordId {
	var result []FileDetailsWithRecordId
	for _, f1 := range pipe.Files {
		for _, f2 := range f1.GeneratedFiles {
			if f2.ArtifactType == pipeline.ArtifactTypesTextEmbeddingVector {
				f := NewFileDetailsWithRecordId(&f2, pipe)
				result = append(result, *f)
			}
		}
	}
	return result
}

func GetListOfPartitionAndSyntheticFiles(pipe *pipeline.DataPipeline) []FileDetailsWithRecordId {
	var result []FileDetailsWithRecordId
	for _, f1 := range pipe.Files {
		for _, f2 := range f1.GeneratedFiles {
			if f2.ArtifactType == pipeline.ArtifactTypesTextPartition || f2.ArtifactType == pipeline.ArtifactTypesSyntheticData {
				f := NewFileDetailsWithRecordId(&f2, pipe)
				result = append(result, *f)
			}
		}
	}
	return result
}

func CreateIndexOnce(ctx context.Context, client memorystorage.IMemoryDb, createdIndexes map[string]struct{}, indexName string, vectorLength int, force bool) {
	ty := reflect.TypeOf(client)
	var key = fmt.Sprintf("%s::%s::%d", ty.String(), indexName, vectorLength)

	if _, ok := createdIndexes[key]; ok && !force {
		return
	}

	err := client.CreateIndex(ctx, indexName, int64(vectorLength))
	if err != nil {
		return
	}
	createdIndexes[key] = struct{}{}
}

func (s *SaveRecordsHandler) GetSourceUrl(ctx context.Context, pipe *pipeline.DataPipeline, file *pipeline.FileDetails) string {
	if file.MimeType != pipeline.MimeTypes_WebPageUrl {
		return ""
	}

	fileContent, err := s.orchestrator.ReadFile(ctx, pipe, file.Name)
	if err != nil {
		return ""
	}
	return string(fileContent)
}

func PrepareRecord(
	pipeline *pipeline.DataPipeline,
	recordId, fileName, url, fileId, partitionFileId, partitionContent string,
	partitionNumber, sectionNumber int64,
	partitionEmbedding *ai.Embedding,
	embeddingGeneratorProvider, embeddingGeneratorName string,
	tags *models.TagCollection,
) *memorystorage.MemoryRecord {
	record := &memorystorage.MemoryRecord{
		Id:      recordId,
		Vector:  partitionEmbedding,
		Tags:    tags,
		Payload: map[string]any{},
	}

	if record.Tags == nil {
		record.Tags = models.NewTagCollection()
	}
	// DOCUMENT DETAILS
	record.Tags.AddOrAppend("ReservedDocumentIdTag", pipeline.DocumentId)

	// FILE DETAILS
	if file := pipeline.GetFile(fileId); file != nil {
		record.Tags.AddOrAppend("ReservedFileTypeTag", file.MimeType)
	}

	record.Tags.AddOrAppend("ReservedFileIdTag", fileId)
	record.Payload["ReservedPayloadFileNameField"] = fileName
	record.Payload["ReservedPayloadUrlField"] = url

	// PARTITION DETAILS
	record.Payload["ReservedPayloadTextField"] = partitionContent
	record.Payload["ReservedPayloadVectorProviderField"] = embeddingGeneratorProvider
	record.Payload["ReservedPayloadVectorGeneratorField"] = embeddingGeneratorName
	record.Tags.AddOrAppend("ReservedFilePartitionTag", partitionFileId)
	record.Tags.AddOrAppend("ReservedFilePartitionNumberTag", fmt.Sprintf("%d", partitionNumber))
	record.Tags.AddOrAppend("ReservedFileSectionNumberTag", fmt.Sprintf("%d", sectionNumber))

	// TIMESTAMP and USER TAGS
	record.Payload["ReservedPayloadLastUpdateField"] = time.Now().UTC().Format(time.RFC3339)

	tags.CopyTo(record.Tags)

	return record
}
