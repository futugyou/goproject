package handlers

import (
	"context"
	"math"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/configuration"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/extensions/chunkers"
	"github.com/futugyou/yomawari/kernel_memory/extensions/tiktoken"
	"github.com/google/uuid"
)

var _ pipeline.IPipelineStepHandler = (*TextPartitioningHandler)(nil)

type TextPartitioningHandler struct {
	orchestrator          pipeline.IPipelineOrchestrator
	options               *configuration.TextPartitioningOptions
	maxTokensPerPartition int
	plainTextChunker      *chunkers.PlainTextChunker
	markDownChunker       *chunkers.MarkDownChunker
	stepName              string
}

func NewTextPartitioningHandler(stepName string, orchestrator pipeline.IPipelineOrchestrator, options *configuration.TextPartitioningOptions) *TextPartitioningHandler {
	tokenizer, _ := tiktoken.NewCL100KTokenizer()
	gens, _ := orchestrator.GetEmbeddingGenerators(context.Background())
	maxTokensPerPartition := math.MaxInt
	for _, gen := range gens {
		maxTokensPerPartition = int(math.Min(float64(gen.GetMaxTokens()), float64(maxTokensPerPartition)))
	}

	return &TextPartitioningHandler{
		orchestrator:          orchestrator,
		options:               options,
		maxTokensPerPartition: maxTokensPerPartition,
		plainTextChunker:      chunkers.NewPlainTextChunker(tokenizer),
		markDownChunker:       chunkers.NewMarkDownChunker(tokenizer),
		stepName:              stepName,
	}
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (t *TextPartitioningHandler) GetStepName() string {
	return t.stepName
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (t *TextPartitioningHandler) SetStepName(name string) {
	t.stepName = name
}

// Invoke implements pipeline.IPipelineStepHandler.
func (t *TextPartitioningHandler) Invoke(ctx context.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	if len(dataPipeline.Files) == 0 {
		return pipeline.ReturnTypeSuccess, dataPipeline, nil
	}

	context := dataPipeline.GetContext()
	maxTokensPerChunk := context.GetCustomPartitioningMaxTokensPerChunkOrDefault(t.options.MaxTokensPerParagraph)
	if maxTokensPerChunk > int64(t.maxTokensPerPartition) {
		return pipeline.ReturnTypeFatalError, dataPipeline, nil
	}

	overlappingTokens := math.Max(0, float64(context.GetCustomPartitioningOverlappingTokensOrDefault(t.options.OverlappingTokens)))
	chunkHeader := context.GetCustomPartitioningChunkHeaderOrDefault(nil)

	for i := range dataPipeline.Files {
		uploadedFile := &dataPipeline.Files[i]
		newFiles := map[string]pipeline.GeneratedFileDetails{}
		for k, file := range uploadedFile.GeneratedFiles {
			if file.AlreadyProcessedBy(t, nil) {
				uploadedFile.GeneratedFiles[k] = file
				continue
			}
			if file.ArtifactType != pipeline.ArtifactTypesExtractedText {
				continue
			}
			content, err := t.orchestrator.ReadFile(ctx, dataPipeline, file.Name)
			if err != nil || len(content) == 0 {
				continue
			}

			var chunks []string
			chunksMimeType := pipeline.MimeTypes_PlainText
			switch file.MimeType {
			case pipeline.MimeTypes_PlainText:
				chunks = t.plainTextChunker.SplitWithOptions(string(content), chunkers.PlainTextChunkerOptions{MaxTokensPerChunk: int(maxTokensPerChunk), Overlap: int(overlappingTokens), ChunkHeader: chunkHeader})
			case pipeline.MimeTypes_MarkDown:
				chunksMimeType = pipeline.MimeTypes_MarkDown
				chunks = t.markDownChunker.SplitWithOptions(string(content), chunkers.MarkDownChunkerOptions{MaxTokensPerChunk: int(maxTokensPerChunk), Overlap: int(overlappingTokens), ChunkHeader: chunkHeader})
			default:
				continue
			}

			if len(chunks) == 0 {
				continue
			}

			for partitionNumber := 0; partitionNumber < len(chunks); partitionNumber++ {
				text := chunks[partitionNumber]
				sectionNumber := 0 // TODO: use this to store the page number (if any)
				textData := []byte(text)
				destFile := uploadedFile.GetPartitionFileName(int64(partitionNumber))
				t.orchestrator.WriteFile(ctx, dataPipeline, destFile, textData)
				destFileDetails := &pipeline.GeneratedFileDetails{
					FileDetailsBase: pipeline.FileDetailsBase{
						Id:              uuid.New().String(),
						Name:            destFile,
						Size:            int64(len(text)),
						MimeType:        chunksMimeType,
						ArtifactType:    pipeline.ArtifactTypesTextPartition,
						PartitionNumber: int64(partitionNumber),
						SectionNumber:   int64(sectionNumber),
						Tags:            dataPipeline.Tags,
						ProcessedBy:     []string{},
						LogEntries:      []pipeline.PipelineLogEntry{},
					},
					ParentId:          uploadedFile.Id,
					SourcePartitionId: "",
					ContentSHA256:     models.CalculateSHA256(string(textData)),
				}
				destFileDetails.MarkProcessedBy(t, nil)
				newFiles[destFile] = *destFileDetails
			}
			file.MarkProcessedBy(t, nil)
			uploadedFile.GeneratedFiles[k] = file
		}

		for key, file := range newFiles {
			uploadedFile.GeneratedFiles[key] = file
		}
	}

	return pipeline.ReturnTypeSuccess, dataPipeline, nil
}
