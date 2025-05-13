package handlers

import (
	rawContext "context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/pipeline"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/prompts"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/text"
	promptsCore "github.com/futugyou/yomawari/kernel_memory/core/prompts"
	"github.com/futugyou/yomawari/kernel_memory/extensions/chunkers"
	"github.com/google/uuid"
)

var _ pipeline.IPipelineStepHandler = (*SummarizationHandler)(nil)

type SummarizationHandler struct {
	MinLength           int64
	orchestrator        pipeline.IPipelineOrchestrator
	summarizationPrompt string
	plainTextChunker    *chunkers.PlainTextChunker
	stepName            string
}

func NewSummarizationHandler(stepName string, orchestrator pipeline.IPipelineOrchestrator, promptProvider prompts.IPromptProvider) *SummarizationHandler {
	handler := &SummarizationHandler{
		MinLength:           30,
		orchestrator:        orchestrator,
		summarizationPrompt: "",
		plainTextChunker:    &chunkers.PlainTextChunker{},
		stepName:            stepName,
	}
	if promptProvider == nil {
		promptProvider = promptsCore.NewEmbeddedPromptProvider()
	}
	if prompt, err := promptProvider.ReadPrompt(rawContext.Background(), constant.PromptNamesSummarize); err == nil {
		handler.summarizationPrompt = *prompt
	}

	return handler
}

// GetStepName implements pipeline.IPipelineStepHandler.
func (s *SummarizationHandler) GetStepName() string {
	return s.stepName
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (s *SummarizationHandler) SetStepName(name string) {
	s.stepName = name
}

// Invoke implements pipeline.IPipelineStepHandler.
func (s *SummarizationHandler) Invoke(ctx rawContext.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	for i := range dataPipeline.Files {
		uploadedFile := &dataPipeline.Files[i]
		summaryFiles := map[string]pipeline.GeneratedFileDetails{}
		for _, file := range uploadedFile.GeneratedFiles {
			if file.AlreadyProcessedBy(s, nil) {
				continue
			}
			if file.ArtifactType != pipeline.ArtifactTypesExtractedText {
				continue
			}

			switch file.MimeType {
			case pipeline.MimeTypes_PlainText, pipeline.MimeTypes_MarkDown:
				content, err := s.orchestrator.ReadFile(ctx, dataPipeline, file.Name)
				if err != nil {
					continue
				}
				summary, success, err := s.Summarizes(ctx, string(content), dataPipeline.GetContext())
				if err == nil && success {
					destFile := uploadedFile.GetHandlerOutputFileName(s, 0)
					err = s.orchestrator.WriteFile(ctx, dataPipeline, destFile, []byte(summary))
					if err != nil {
						continue
					}
					tags := dataPipeline.Tags.Clone()
					tags.AddSyntheticTag(constant.TagsSyntheticSummary)
					size := len(summary)
					summaryFiles[destFile] = pipeline.GeneratedFileDetails{
						FileDetailsBase: pipeline.FileDetailsBase{
							Id:              uuid.New().String(),
							Name:            destFile,
							Size:            int64(size),
							MimeType:        pipeline.MimeTypes_PlainText,
							ArtifactType:    pipeline.ArtifactTypesSyntheticData,
							PartitionNumber: 0,
							SectionNumber:   0,
							Tags:            tags,
							ProcessedBy:     []string{},
							LogEntries:      []pipeline.PipelineLogEntry{},
						},
						ParentId:          uploadedFile.Id,
						SourcePartitionId: "",
						ContentSHA256:     models.CalculateSHA256(summary),
					}
				}
			default:
				continue
			}
		}

		for key, file := range summaryFiles {
			file.MarkProcessedBy(s, nil)
			uploadedFile.GeneratedFiles[key] = file
		}
	}

	return pipeline.ReturnTypeSuccess, dataPipeline, nil
}

func (s *SummarizationHandler) Summarizes(ctx rawContext.Context, content string, context context.IContext) (string, bool, error) {
	textGenerator, err := s.orchestrator.GetTextGenerators(ctx)
	if err != nil {
		return "", false, err
	}

	contentLength := textGenerator.CountTokens(ctx, content)
	if contentLength < s.MinLength {
		return content, true, nil
	}

	targetSummarySize := textGenerator.GetMaxTokenTotal() / 2
	customTargetSummarySize := context.GetCustomSummaryTargetTokenSizeOrDefault(-1)

	if customTargetSummarySize > 0 {
		if customTargetSummarySize > textGenerator.GetMaxTokenTotal()/2 {
			return "", false, fmt.Errorf("custom summary size is too large, the max value allowed is %d (50%% of the model capacity)", textGenerator.GetMaxTokenTotal()/2)
		}

		if customTargetSummarySize < 15 {
			return "", false, fmt.Errorf("custom summary size is too small, the min value allowed is %d", 15)
		}
		targetSummarySize = customTargetSummarySize
	}

	maxTokensPerParagraph := textGenerator.GetMaxTokenTotal() / 4
	overlappingTokens := context.GetCustomSummaryOverlappingTokensOrDefault(textGenerator.GetMaxTokenTotal() / 16)
	done := false
	summarizationPrompt := context.GetCustomSummaryPromptOrDefault(s.summarizationPrompt)
	var overlapToRemove = overlappingTokens > 0
	maxInputTokens := textGenerator.GetMaxTokenTotal() / 2
	firstRun := overlapToRemove
	previousLength := contentLength

	for {
		if done {
			break
		}
		chunks := []string{}
		if contentLength <= maxInputTokens {
			overlapToRemove = false
			chunks = append(chunks, content)
		} else {
			chunks = s.plainTextChunker.SplitWithOptions(content, chunkers.PlainTextChunkerOptions{
				MaxTokensPerChunk: int(maxTokensPerParagraph),
				Overlap:           int(overlappingTokens),
			})
		}
		newContent := &strings.Builder{}
		for index := 0; index < len(chunks); index++ {
			paragraph := chunks[index]
			var filledPrompt = strings.Replace(summarizationPrompt, "{{$input}}", paragraph, 1)
			generateTextResponse := textGenerator.GenerateText(ctx, filledPrompt, nil)
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				for response := range generateTextResponse {
					if response.Err != nil {
						newContent.WriteString(response.Content.ToString())
					}
				}
			}()
			wg.Wait()
			text.AppendLine(newContent)
		}

		content = newContent.String()
		contentLength = textGenerator.CountTokens(ctx, content)
		if !firstRun && contentLength >= previousLength {
			log.Printf("summarization stopped, the content is not getting shorter: %d tokens => %d tokens. The summary has been saved but is longer than requested.", previousLength, contentLength)
			return content, true, nil
		}
		previousLength = contentLength

		firstRun = false
		done = !overlapToRemove && (contentLength <= targetSummarySize)
	}
	return content, true, nil
}
