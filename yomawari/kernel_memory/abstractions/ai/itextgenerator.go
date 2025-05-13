package ai

import (
	"context"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
)

type ITextGenerator interface {
	ITextTokenizer
	GetMaxTokenTotal() int64
	GenerateText(ctx context.Context, prompt string, options *TextGenerationOptions) <-chan GenerateTextResponse
}

type GenerateTextResponse struct {
	Content *models.GeneratedTextContent
	Err     error
}
type TextGenerationOptions struct {
	Temperature          float64
	NucleusSampling      float64
	PresencePenalty      float64
	FrequencyPenalty     float64
	MaxTokens            *int64
	StopSequences        []string
	ResultsPerPrompt     int64
	TokenSelectionBiases map[int]float32
}
