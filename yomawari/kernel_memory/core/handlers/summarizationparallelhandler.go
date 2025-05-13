package handlers

import (
	rawContext "context"
	"log"
	"math"
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
	"golang.org/x/sync/errgroup"
)

var _ pipeline.IPipelineStepHandler = (*SummarizationParallelHandler)(nil)

type SummarizationParallelHandler struct {
	MinLength           int64
	orchestrator        pipeline.IPipelineOrchestrator
	summarizationPrompt string
	plainTextChunker    *chunkers.PlainTextChunker
	stepName            string
}

func NewSummarizationParallelHandler(stepName string, orchestrator pipeline.IPipelineOrchestrator, promptProvider prompts.IPromptProvider) *SummarizationParallelHandler {
	handler := &SummarizationParallelHandler{
		MinLength:           50,
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
func (s *SummarizationParallelHandler) GetStepName() string {
	return s.stepName
}

// SetStepName implements pipeline.IPipelineStepHandler.
func (s *SummarizationParallelHandler) SetStepName(name string) {
	s.stepName = name
}

// Invoke implements pipeline.IPipelineStepHandler.
func (s *SummarizationParallelHandler) Invoke(ctx rawContext.Context, dataPipeline *pipeline.DataPipeline) (pipeline.ReturnType, *pipeline.DataPipeline, error) {
	g, ctx := errgroup.WithContext(ctx)

	var mu sync.Mutex

	for i := range dataPipeline.Files {
		uploadedFile := &dataPipeline.Files[i]

		g.Go(func() error {
			summaryFiles := make(map[string]pipeline.GeneratedFileDetails)

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
					if err != nil {
						continue
					}
					if success {
						destFile := uploadedFile.GetHandlerOutputFileName(s, 0)
						err = s.orchestrator.WriteFile(ctx, dataPipeline, destFile, []byte(summary))
						if err != nil {
							continue
						}

						tags := dataPipeline.Tags.Clone()
						tags.AddSyntheticTag(constant.TagsSyntheticSummary)

						mu.Lock()
						summaryFiles[destFile] = pipeline.GeneratedFileDetails{
							FileDetailsBase: pipeline.FileDetailsBase{
								Id:              uuid.New().String(),
								Name:            destFile,
								Size:            int64(len(summary)),
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
						mu.Unlock()
					}
				default:
					continue
				}
			}

			mu.Lock()
			for key, file := range summaryFiles {
				file.MarkProcessedBy(s, nil)
				uploadedFile.GeneratedFiles[key] = file
			}
			mu.Unlock()

			return nil
		})
	}

	_ = g.Wait()

	return pipeline.ReturnTypeSuccess, dataPipeline, nil
}

func (s *SummarizationParallelHandler) Summarizes(ctx rawContext.Context, content string, context context.IContext) (string, bool, error) {
	textGenerator, err := s.orchestrator.GetTextGenerators(ctx)
	if err != nil {
		return "", false, err
	}

	summaryMaxTokens := textGenerator.GetMaxTokenTotal() / 2                            // 50% of model capacity
	maxTokensPerChunk := summaryMaxTokens / 2                                           // 25% of model capacity
	overlappingTokens := math.Min(math.Max(200, float64(maxTokensPerChunk/2)), 500) / 2 // 100...250

	contentLength := textGenerator.CountTokens(ctx, content)
	if contentLength < s.MinLength {
		return content, false, nil
	}

	done := false
	var overlapToRemove = overlappingTokens > 0
	firstRun := overlapToRemove
	previousLength := contentLength

	for {
		if done {
			break
		}
		paragraphs := []string{}
		if contentLength <= summaryMaxTokens {
			overlapToRemove = false
			paragraphs = append(paragraphs, content)
		} else {
			paragraphs = s.plainTextChunker.SplitWithOptions(content, chunkers.PlainTextChunkerOptions{
				MaxTokensPerChunk: int(maxTokensPerChunk),
				Overlap:           int(overlappingTokens),
			})
		}
		newContent := &strings.Builder{}
		for index := 0; index < len(paragraphs); index++ {
			paragraph := paragraphs[index]
			var filledPrompt = strings.Replace(s.summarizationPrompt, "{{$input}}", paragraph, 1)
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
			log.Printf("summarization failed, the content is getting longer: %d tokens => %d tokens", previousLength, contentLength)
			return content, true, nil
		}
		previousLength = contentLength

		firstRun = false
		done = !overlapToRemove && (contentLength <= summaryMaxTokens)
	}
	return content, true, nil
}
