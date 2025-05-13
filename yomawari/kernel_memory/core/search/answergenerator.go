package search

import (
	"context"
	"strings"
	"time"

	"github.com/futugyou/yomawari/kernel_memory/abstractions/ai"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/constant"
	aicontext "github.com/futugyou/yomawari/kernel_memory/abstractions/context"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/models"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/prompts"
	"github.com/futugyou/yomawari/kernel_memory/abstractions/search"
)

type AnswerGenerator struct {
	contentModeration ai.IContentModeration
	config            *search.SearchClientConfig
	answerPrompt      string
	textGenerator     ai.ITextGenerator
}

func NewAnswerGenerator(config *search.SearchClientConfig, contentModeration ai.IContentModeration, textGenerator ai.ITextGenerator, promptProvider prompts.IPromptProvider) *AnswerGenerator {
	g := &AnswerGenerator{
		contentModeration: contentModeration,
		config:            config,
		answerPrompt:      config.FactTemplate,
		textGenerator:     textGenerator,
	}
	answerPrompt, err := promptProvider.ReadPrompt(context.Background(), constant.PromptNamesAnswerWithFacts)
	if err == nil {
		g.answerPrompt = *answerPrompt
	}
	return g
}

func (a *AnswerGenerator) generateAnswerTokens(ctx context.Context, prompt string, context aicontext.IContext) <-chan ai.GenerateTextResponse {
	maxTokens := context.GetCustomRagMaxTokensOrDefault(int64(a.config.AnswerTokens))
	temperature := context.GetCustomRagTemperatureOrDefault(a.config.Temperature)
	nucleusSampling := context.GetCustomRagNucleusSamplingOrDefault(a.config.TopP)

	var options = &ai.TextGenerationOptions{
		MaxTokens:            &maxTokens,
		Temperature:          temperature,
		NucleusSampling:      nucleusSampling,
		PresencePenalty:      a.config.PresencePenalty,
		FrequencyPenalty:     a.config.FrequencyPenalty,
		StopSequences:        a.config.StopSequences,
		TokenSelectionBiases: a.config.TokenSelectionBiases,
	}

	return a.textGenerator.GenerateText(ctx, prompt, options)
}

func (a *AnswerGenerator) preparePrompt(question string, facts string, context aicontext.IContext) string {
	prompt := context.GetCustomRagPromptOrDefault(a.answerPrompt)
	emptyAnswer := context.GetCustomEmptyAnswerTextOrDefault(a.config.EmptyAnswer)

	question = strings.TrimSpace(question)
	if !strings.HasSuffix(question, "?") {
		question = question + "?"
	}

	prompt = strings.ReplaceAll(prompt, "{{$facts}}", strings.TrimSpace(facts))
	prompt = strings.ReplaceAll(prompt, "{{$input}}", question)
	prompt = strings.ReplaceAll(prompt, "{{$notFound}}", emptyAnswer)

	return prompt
}

func (s *AnswerGenerator) GenerateAnswer(ctx context.Context, question string, result *SearchClientResult, context aicontext.IContext) <-chan models.MemoryAnswer {
	out := make(chan models.MemoryAnswer)
	go func() {
		defer close(out)

		prompt := s.preparePrompt(question, result.Facts.String(), context)
		promptSize := int(s.textGenerator.CountTokens(ctx, prompt))
		modelType := constant.TextGeneration
		tokenUsage := models.TokenUsage{
			Timestamp:         time.Now().UTC(),
			ModelType:         &modelType,
			TokenizerTokensIn: &promptSize,
		}
		result.AddTokenUsageToStaticResults(tokenUsage)

		if result.FactsAvailableCount > 0 && result.FactsUsedCount == 0 {
			out <- *result.InsufficientTokensResult
			return
		}

		if result.FactsUsedCount == 0 {
			out <- *result.NoFactsResult
			return
		}

		var completeAnswerTokens strings.Builder
		for answerToken := range s.generateAnswerTokens(ctx, prompt, context) {
			if answerToken.Err != nil {
				continue
			}
			completeAnswerTokens.WriteString(answerToken.Content.Text)
			tokenUsage.Merge(answerToken.Content.TokenUsage)
			result.AskResult.Result = answerToken.Content.Text

			out <- *result.AskResult
		}

		completeAnswer := completeAnswerTokens.String()
		if strings.TrimSpace(completeAnswer) == "" || ValueIsEquivalentTo(completeAnswer, s.config.EmptyAnswer) {
			out <- *result.NoFactsResult
			return
		}

		if s.config.UseContentModeration && s.contentModeration != nil {
			isSafe := s.contentModeration.IsSafe(ctx, completeAnswer)
			if !isSafe {
				out <- *result.UnsafeAnswerResult
				return
			}
		}

		to := int(s.textGenerator.CountTokens(ctx, completeAnswer))
		result.AskResult.Result = ""
		tokenUsage.TokenizerTokensOut = &to
		result.AskResult.TokenUsage = []models.TokenUsage{tokenUsage}
		out <- *result.AskResult
	}()
	return out
}
