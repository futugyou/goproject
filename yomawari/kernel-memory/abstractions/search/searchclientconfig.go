package search

import (
	"errors"
)

type SearchClientConfig struct {
	MaxAskPromptSize      int             `json:"max_ask_prompt_size"`
	MaxMatchesCount       int             `json:"max_matches_count"`
	AnswerTokens          int             `json:"answer_tokens"`
	EmptyAnswer           string          `json:"empty_answer"`
	FactTemplate          string          `json:"fact_template"`
	IncludeDuplicateFacts bool            `json:"include_duplicate_facts"`
	Temperature           float64         `json:"temperature"`
	TopP                  float64         `json:"top_p"`
	PresencePenalty       float64         `json:"presence_penalty"`
	FrequencyPenalty      float64         `json:"frequency_penalty"`
	StopSequences         []string        `json:"stop_sequences"`
	TokenSelectionBiases  map[int]float32 `json:"token_selection_biases"`
	UseContentModeration  bool            `json:"use_content_moderation"`
	ModeratedAnswer       string          `json:"moderated_answer"`
}

func NewSearchClientConfig() *SearchClientConfig {
	return &SearchClientConfig{
		MaxAskPromptSize: -1,
		MaxMatchesCount:  100,
		AnswerTokens:     300,
		EmptyAnswer:      "INFO NOT FOUND",
		// "==== [File:{{$source}};Relevance:{{$relevance}}]:\n{{$content}}",
		FactTemplate:          "==== [File:{{$source}};Relevance:{{$relevance}}]:\n{{$content}}",
		IncludeDuplicateFacts: false,
		Temperature:           0,
		TopP:                  0,
		PresencePenalty:       0,
		FrequencyPenalty:      0,
		StopSequences:         []string{},
		TokenSelectionBiases:  make(map[int]float32),
		UseContentModeration:  true,
		ModeratedAnswer:       "Sorry, the generated content contains unsafe or inappropriate information.",
	}
}

func (c *SearchClientConfig) Validate() error {
	if c.MaxAskPromptSize > 0 && c.MaxAskPromptSize < 1024 {
		return errors.New("SearchClient: MaxAskPromptSize cannot be less than 1024")
	}
	if c.MaxMatchesCount < 1 {
		return errors.New("SearchClient: MaxMatchesCount cannot be less than 1")
	}
	if c.AnswerTokens < 1 {
		return errors.New("SearchClient: AnswerTokens cannot be less than 1")
	}
	if len(c.EmptyAnswer) > 256 {
		return errors.New("SearchClient: EmptyAnswer is too long, consider something shorter")
	}
	if c.Temperature < 0 || c.Temperature > 2 {
		return errors.New("SearchClient: Temperature must be between 0 and 2")
	}
	if c.TopP < 0 || c.TopP > 2 {
		return errors.New("SearchClient: TopP must be between 0 and 2")
	}
	if c.PresencePenalty < -2 || c.PresencePenalty > 2 {
		return errors.New("SearchClient: PresencePenalty must be between -2 and 2")
	}
	if c.FrequencyPenalty < -2 || c.FrequencyPenalty > 2 {
		return errors.New("SearchClient: FrequencyPenalty must be between -2 and 2")
	}
	return nil
}
