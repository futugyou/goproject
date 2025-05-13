package search

import (
	"strings"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type SearchMode int

type SearchState int

const (
	SearchModeSearch SearchMode = iota
	SearchModeAsk
)

const (
	SearchStateContinue SearchState = iota
	SearchStateSkipRecord
	SearchStateStop
)

type SearchClientResult struct {
	Mode                     SearchMode
	State                    SearchState
	RecordCount              int
	MaxRecordCount           int
	AskResult                *models.MemoryAnswer
	NoFactsResult            *models.MemoryAnswer
	NoQuestionResult         *models.MemoryAnswer
	UnsafeAnswerResult       *models.MemoryAnswer
	InsufficientTokensResult *models.MemoryAnswer
	SearchResult             *models.SearchResult
	Facts                    strings.Builder
	FactsAvailableCount      int
	FactsUsedCount           int
	TokensAvailable          int
	FactsUniqueness          map[string]struct{}
}

func NewAskResultInstance(question, emptyAnswer, moderatedAnswer string, maxGroundingFacts, tokensAvailable int) *SearchClientResult {
	r1 := "No relevant memories available"
	r2 := "No question provided"
	r3 := "Unable to use memory, max tokens reached"
	r4 := "Content moderation"
	return &SearchClientResult{
		Mode:            SearchModeAsk,
		TokensAvailable: tokensAvailable,
		MaxRecordCount:  maxGroundingFacts,
		AskResult: &models.MemoryAnswer{
			StreamState: &models.StreamStateAppend,
			Question:    &question,
			NoResult:    false,
		},
		NoFactsResult: &models.MemoryAnswer{
			StreamState:    &models.StreamStateReset,
			Question:       &question,
			NoResult:       true,
			NoResultReason: &r1,
			Result:         emptyAnswer,
		},
		NoQuestionResult: &models.MemoryAnswer{
			StreamState:    &models.StreamStateReset,
			Question:       &question,
			NoResult:       true,
			NoResultReason: &r2,
			Result:         emptyAnswer,
		},
		InsufficientTokensResult: &models.MemoryAnswer{
			StreamState:    &models.StreamStateReset,
			Question:       &question,
			NoResult:       true,
			NoResultReason: &r3,
			Result:         emptyAnswer,
		},
		UnsafeAnswerResult: &models.MemoryAnswer{
			StreamState:    &models.StreamStateReset,
			Question:       &question,
			NoResult:       true,
			NoResultReason: &r4,
			Result:         moderatedAnswer,
		},
		FactsUniqueness: make(map[string]struct{}),
	}
}

func NewSearchResultInstance(query string, maxSearchResults int) *SearchClientResult {
	return &SearchClientResult{
		Mode:           SearchModeSearch,
		MaxRecordCount: maxSearchResults,
		SearchResult: &models.SearchResult{
			Query:   query,
			Results: []models.Citation{},
		},
		FactsUniqueness: make(map[string]struct{}),
	}
}

func (s *SearchClientResult) AddSource(citation models.Citation) {
	s.SearchResult.Results = append(s.SearchResult.Results, citation)
	s.AskResult.RelevantSources = append(s.AskResult.RelevantSources, citation)
	s.InsufficientTokensResult.RelevantSources = append(s.InsufficientTokensResult.RelevantSources, citation)
	s.UnsafeAnswerResult.RelevantSources = append(s.UnsafeAnswerResult.RelevantSources, citation)
}

func (s *SearchClientResult) AddTokenUsageToStaticResults(tokenUsage models.TokenUsage) {
	s.InsufficientTokensResult.TokenUsage = []models.TokenUsage{tokenUsage}
	s.UnsafeAnswerResult.TokenUsage = []models.TokenUsage{tokenUsage}
	s.NoFactsResult.TokenUsage = []models.TokenUsage{tokenUsage}
}

func (s *SearchClientResult) SkipRecord() *SearchClientResult {
	s.State = SearchStateSkipRecord
	return s
}

func (s *SearchClientResult) Stop() *SearchClientResult {
	s.State = SearchStateStop
	return s
}

var trimChars = ".\"'`~!?@#$%^+*_-=|\\/()[]{}<>"

func ValueIsEquivalentTo(value, target string) bool {
	value = strings.Trim(value, trimChars)
	target = strings.Trim(target, trimChars)
	return strings.EqualFold(value, target)
}
