package handlers

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/documentstorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/google/uuid"
)

type PartitionInfo struct {
	GeneratedFile    map[string]pipeline.GeneratedFileDetails
	UploadedFile     *pipeline.FileDetails
	PartitionContent string
}

func NewPartitionInfo(generatedFile map[string]pipeline.GeneratedFileDetails,
	uploadedFile *pipeline.FileDetails,
	partitionContent string) *PartitionInfo {
	return &PartitionInfo{
		GeneratedFile:    generatedFile,
		UploadedFile:     uploadedFile,
		PartitionContent: partitionContent,
	}
}

type GenerateEmbeddingsHandlerBase struct {
	orchestrator   pipeline.IPipelineOrchestrator
	ActualInstance pipeline.IPipelineStepHandler
}

func NewGenerateEmbeddingsHandlerBase(orchestrator pipeline.IPipelineOrchestrator) *GenerateEmbeddingsHandlerBase {
	return &GenerateEmbeddingsHandlerBase{
		orchestrator: orchestrator,
	}
}

func GetEmbeddingFileName(srcFilename string, t string, embeddingName string) string {
	return fmt.Sprintf("%s.%s.%s%s", srcFilename, t, embeddingName, pipeline.FileExtensions_TextEmbeddingVector)
}

func (handler *GenerateEmbeddingsHandlerBase) GetListOfPartitionsToProcess(ctx context.Context, pipe *pipeline.DataPipeline, subStepName string) ([]PartitionInfo, error) {
	if handler == nil || handler.ActualInstance == nil {
		return nil, fmt.Errorf("handler is nil")
	}
	partitionsToProcess := []PartitionInfo{}

	for i := range pipe.Files {
		uploadedFile := &pipe.Files[i]
		for k, v := range uploadedFile.GeneratedFiles {
			if v.ArtifactType != pipeline.ArtifactTypesTextPartition && v.ArtifactType != pipeline.ArtifactTypesSyntheticData {
				continue
			}

			if v.AlreadyProcessedBy(handler.ActualInstance, &subStepName) {
				uploadedFile.GeneratedFiles[k] = v
				continue
			}
			if v.MimeType == pipeline.MimeTypes_PlainText || v.MimeType == pipeline.MimeTypes_MarkDown {
				partitionContent, err := handler.orchestrator.ReadTextFile(ctx, pipe, v.Name)
				if err != nil {
					continue
				}
				partitionsToProcess = append(partitionsToProcess, PartitionInfo{
					GeneratedFile:    map[string]pipeline.GeneratedFileDetails{k: v},
					UploadedFile:     uploadedFile,
					PartitionContent: *partitionContent,
				})
			}
		}
	}
	return partitionsToProcess, nil
}

func (handler *GenerateEmbeddingsHandlerBase) SaveEmbeddingsToDocumentStorage(ctx context.Context, pipe *pipeline.DataPipeline,
	partitions []PartitionInfo, embeddings []ai.Embedding, generatorProvider string, generatorName string) error {
	if handler == nil {
		return fmt.Errorf("handler is nil")
	}
	if len(partitions) != len(embeddings) {
		return fmt.Errorf("the list of embeddings doesn't match the list of text partitions. The two lists have different size: %d embeddings != %d text partitions", len(partitions), len(embeddings))
	}
	for i := 0; i < len(partitions); i++ {
		handler.SaveEmbeddingToDocumentStorage(ctx, pipe, partitions[i], embeddings[i], generatorProvider, generatorName)
	}
	return nil
}

func (handler *GenerateEmbeddingsHandlerBase) SaveEmbeddingToDocumentStorage(ctx context.Context, pipe *pipeline.DataPipeline,
	partition PartitionInfo, embedding ai.Embedding, generatorProvider string, generatorName string) error {

	for k, v := range partition.GeneratedFile {
		embeddingData := documentstorage.EmbeddingFileContent{
			GeneratorName:     generatorName,
			GeneratorProvider: generatorProvider,
			VectorSize:        int64(embedding.Length()),
			SourceFileName:    v.Name,
			Vector:            embedding,
			TimeStamp:         time.Now().UTC(),
		}

		embeddingDataAsJson := embeddingData.ToJson()
		embeddingDataFileName := GetEmbeddingFileName(v.Name, generatorProvider, generatorName)
		handler.orchestrator.WriteTextFile(ctx, pipe, embeddingDataFileName, embeddingDataAsJson)
		subStep := GetSubStepName2(generatorProvider, generatorName)
		v.MarkProcessedBy(handler.ActualInstance, &subStep)
		partition.GeneratedFile[k] = v
	}

	return nil
}

func (handler *GenerateEmbeddingsHandlerBase) TrackNewFileInPipelineStatus(ctx context.Context, newFileName string,
	newFileSize int, sourcePartitionFile pipeline.GeneratedFileDetails, sourceUserFile pipeline.FileDetails) error {
	newFileDetails := pipeline.GeneratedFileDetails{
		FileDetailsBase: pipeline.FileDetailsBase{
			Id:              uuid.NewString(),
			Name:            newFileName,
			Size:            int64(newFileSize),
			MimeType:        pipeline.MimeTypes_TextEmbeddingVector,
			ArtifactType:    pipeline.ArtifactTypesTextEmbeddingVector,
			PartitionNumber: sourcePartitionFile.PartitionNumber,
			SectionNumber:   sourcePartitionFile.SectionNumber,
			Tags:            sourcePartitionFile.Tags,
			ProcessedBy:     []string{},
			LogEntries:      []pipeline.PipelineLogEntry{},
		},
		ParentId:          sourceUserFile.Id,
		SourcePartitionId: sourcePartitionFile.Id,
		ContentSHA256:     "",
	}
	newFileDetails.MarkProcessedBy(handler.ActualInstance, nil)
	sourceUserFile.GeneratedFiles[newFileName] = newFileDetails
	return nil
}

func GetSubStepName(generator interface{}) string {
	return GetSubStepName2(GetEmbeddingProviderName(generator), GetEmbeddingGeneratorName(generator))
}

func GetSubStepName2(providerName string, generatorName string) string {
	return fmt.Sprintf("%s/%s", providerName, generatorName)
}

func GetEmbeddingProviderName(generator interface{}) string {
	generatorType := reflect.TypeOf(generator)
	if generatorType.Kind() == reflect.Ptr {
		generatorType = generatorType.Elem()
	}

	generatorProviderClassName := generatorType.Name()
	if generatorProviderClassName == "" {
		generatorProviderClassName = generatorType.String()
	}

	parts := strings.Split(generatorProviderClassName, ".")
	return strings.Join(parts[len(parts)-3:], ".")
}

func GetEmbeddingGeneratorName(object any) string {
	return "__"
}
