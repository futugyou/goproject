package search

import (
	rawContext "context"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/memorystorage"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/prompts"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/search"
	promptsCore "github.com/futugyou/yomawari/kernel_memory/core/prompts"
)

type SearchClient struct {
	memoryDb        memorystorage.IMemoryDb
	textGenerator   ai.ITextGenerator
	config          *search.SearchClientConfig
	answerGenerator *AnswerGenerator
	answerPrompt    string
}

func NewSearchClient(memoryDb memorystorage.IMemoryDb, contentModeration ai.IContentModeration, textGenerator ai.ITextGenerator, promptProvider prompts.IPromptProvider, config *search.SearchClientConfig) *SearchClient {
	if promptProvider == nil {
		promptProvider = promptsCore.NewEmbeddedPromptProvider()
	}
	s := &SearchClient{
		memoryDb:        memoryDb,
		textGenerator:   textGenerator,
		config:          config,
		answerGenerator: NewAnswerGenerator(config, contentModeration, textGenerator, promptProvider),
		answerPrompt:    "",
	}
	if prompt, err := promptProvider.ReadPrompt(rawContext.Background(), constant.PromptNamesAnswerWithFacts); err == nil {
		s.answerPrompt = *prompt
	}

	return s
}

// Ask implements search.ISearchClient.
func (s *SearchClient) Ask(ctx rawContext.Context, index string, question string, filters []models.MemoryFilter, minRelevance float64, context context.IContext) (*models.MemoryAnswer, error) {
	result := &models.MemoryAnswer{}

	stream := s.AskStreaming(ctx, index, question, filters, minRelevance, context)

	done := false
	text := strings.Builder{}
	text.WriteString(result.Result)

	for {
		if done {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case part := <-stream:
			if part.Answer != nil {
				result.TokenUsage = part.Answer.TokenUsage

				switch part.Answer.StreamState {
				case &models.StreamStateError:
					text = strings.Builder{}
					result = part.Answer
					done = true
				case &models.StreamStateReset:
					text = strings.Builder{}
					text.WriteString(part.Answer.Result)
					result = part.Answer
				case &models.StreamStateAppend:
					result.NoResult = part.Answer.NoResult
					result.NoResultReason = part.Answer.NoResultReason
					text.WriteString(part.Answer.Result)
					result.Result = text.String()
					if result.RelevantSources != nil && part.Answer.RelevantSources != nil {
						result.RelevantSources = append(result.RelevantSources, part.Answer.RelevantSources...)
					}
				case &models.StreamStateLast:
					result.NoResult = part.Answer.NoResult
					result.NoResultReason = part.Answer.NoResultReason
					text.WriteString(part.Answer.Result)
					result.Result = text.String()
					if result.RelevantSources != nil && part.Answer.RelevantSources != nil {
						result.RelevantSources = append(result.RelevantSources, part.Answer.RelevantSources...)
					}

					done = true
				}
			}

		}

	}

	result.Question = &question
	result.StreamState = nil
	return result, nil
}

// AskStreaming implements search.ISearchClient.
func (s *SearchClient) AskStreaming(ctx rawContext.Context, index string, question string, filters []models.MemoryFilter, minRelevance float64, context context.IContext) <-chan search.AskStreamingStreamResponse {
	out := make(chan search.AskStreamingStreamResponse)
	go func() {
		defer close(out)

		emptyAnswer := s.config.EmptyAnswer
		answerPrompt := s.answerPrompt
		limit := s.config.MaxMatchesCount
		includeDuplicateFacts := s.config.IncludeDuplicateFacts

		maxTokens := s.config.MaxAskPromptSize
		if maxTokens <= 0 {
			maxTokens = int(s.textGenerator.GetMaxTokenTotal())
		}

		result := NewAskResultInstance(
			question, emptyAnswer, s.config.ModeratedAnswer, limit,
			maxTokens-int(s.textGenerator.CountTokens(ctx, answerPrompt))-int(s.textGenerator.CountTokens(ctx, question))-s.config.AnswerTokens,
		)

		if question == "" {
			out <- search.AskStreamingStreamResponse{
				Answer: result.NoQuestionResult,
			}
			return
		}

		matches := s.memoryDb.GetSimilarList(
			ctx, index, question, filters, minRelevance, int64(limit), false,
		)

		factTemplate := s.config.FactTemplate
		if !strings.HasSuffix(factTemplate, "\n") {
			factTemplate += "\n"
		}

		for match := range matches {
			if match.Err != nil || match.Record == nil {
				continue
			}

			memoryRecord := match.Record
			var recordRelevance float64 = 0
			if match.Similars != nil {
				recordRelevance = *match.Similars
			}
			result.State = SearchStateContinue
			result = s.processMemoryRecord(ctx, result, index, memoryRecord, recordRelevance, includeDuplicateFacts, &factTemplate)

			if result.State == SearchStateSkipRecord {
				continue
			}
			if result.State == SearchStateStop {
				break
			}
		}

		first := true
		for answer := range s.answerGenerator.GenerateAnswer(ctx, question, result, context) {
			if first {
				first = false
			} else {
				result.AskResult.RelevantSources = nil
				result.AskResult.Question = nil
			}
			out <- search.AskStreamingStreamResponse{
				Answer: &answer,
			}
		}
	}()
	return out
}

// ListIndexes implements search.ISearchClient.
func (s *SearchClient) ListIndexes(ctx rawContext.Context) ([]string, error) {
	return s.memoryDb.GetIndexes(ctx)
}

// Search implements search.ISearchClient.
func (s *SearchClient) Search(ctx rawContext.Context, index string, query string, filters []models.MemoryFilter, minRelevance float64, limit int, context context.IContext) (*models.SearchResult, error) {
	if limit <= 0 {
		limit = s.config.MaxMatchesCount
	}

	result := NewSearchResultInstance(query, limit)

	if len(query) == 0 && len(filters) == 0 {
		return result.SearchResult, nil
	}

	var match <-chan memorystorage.MemoryRecordChanResponse

	if len(query) == 0 {
		match = s.memoryDb.GetList(ctx, index, filters, int64(limit), false)
	} else {
		match = s.memoryDb.GetSimilarList(ctx, index, query, filters, minRelevance, int64(limit), false)
	}
	for data := range match {
		if data.Err != nil || data.Record == nil {
			continue
		}
		result.State = SearchStateContinue
		var recordRelevance float64 = 0
		if data.Similars != nil {
			recordRelevance = *data.Similars
		}
		result = s.processMemoryRecord(ctx, result, index, data.Record, recordRelevance, true, nil)

		if result.State == SearchStateSkipRecord {
			continue
		}

		if result.State == SearchStateStop {
			break
		}
	}

	return result.SearchResult, nil
}

func (s *SearchClient) processMemoryRecord(ctx rawContext.Context, result *SearchClientResult, index string, record *memorystorage.MemoryRecord,
	recordRelevance float64, includeDupes bool, factTemplate *string) *SearchClientResult {
	partitionText := strings.TrimSpace(record.GetPartitionText())
	if len(partitionText) == 0 {
		return result.SkipRecord()
	}

	result.RecordCount++
	documentId := record.GetDocumentId()
	fileId := record.GetFileId()
	linkToFile := fmt.Sprintf("%s/%s/%s", index, documentId, fileId)
	fileName := record.GetFileName()
	fileDownloadUrl := record.GetWebPageUrl(index)
	fileNameForLLM := fileName
	if fileName == "content.url" {
		fileNameForLLM = fileDownloadUrl
	}
	isDupe := false
	if _, ok := result.FactsUniqueness[partitionText]; ok {
		isDupe = true
	}
	skipFactInPrompt := (isDupe && !includeDupes)

	if result.Mode == SearchModeSearch {
		if recordRelevance > math.MinInt16 {
			log.Printf("adding result with relevance %f/n", recordRelevance)
		}
	} else if result.Mode == SearchModeAsk {
		result.FactsAvailableCount++

		if !skipFactInPrompt {
			fact := promptsCore.RenderFactTemplate(
				*factTemplate,
				partitionText,
				fileNameForLLM,
				fmt.Sprintf("%.1f%%", recordRelevance*100),
				record.Id,
				record.Tags,
				record.Payload)

			// Use the partition/chunk only if there's room for it
			factSizeInTokens := s.textGenerator.CountTokens(ctx, fact)
			if factSizeInTokens >= int64(result.TokensAvailable) {
				// Stop after reaching the max number of tokens
				return result.Stop()
			}

			result.Facts.WriteString(fact)
			result.FactsUsedCount++
			result.TokensAvailable -= int(factSizeInTokens)
		} else {
			result.FactsUsedCount++
		}
	}

	var citation models.Citation = models.Citation{}
	switch result.Mode {
	case SearchModeSearch:
		for _, v := range result.SearchResult.Results {
			if v.Link == linkToFile {
				citation = v
				break
			}
		}
	case SearchModeAsk:
		for _, v := range result.AskResult.RelevantSources {
			if v.Link == linkToFile {
				citation = v
				break
			}
		}

	}
	citation.Index = index
	citation.DocumentId = documentId
	citation.FileId = fileId
	citation.Link = linkToFile
	citation.SourceContentType = record.GetFileContentType()
	citation.SourceName = fileName
	citation.SourceUrl = &fileDownloadUrl
	citation.Partitions = append(citation.Partitions, models.Partition{
		Text:            partitionText,
		Relevance:       recordRelevance,
		PartitionNumber: int64(record.GetPartitionNumber()),
		SectionNumber:   int64(record.GetSectionNumber()),
		LastUpdate:      record.GetLastUpdate(),
		Tags:            record.Tags,
	})

	result.AddSource(citation)
	// Stop when reaching the max number of results or facts. This acts also as
	// a protection against storage connectors disregarding 'limit' and returning too many records.
	if (result.Mode == SearchModeSearch && len(result.SearchResult.Results) >= result.MaxRecordCount) || (result.Mode == SearchModeAsk && result.FactsUsedCount >= result.MaxRecordCount) {
		return result.Stop()
	}

	return result
}
