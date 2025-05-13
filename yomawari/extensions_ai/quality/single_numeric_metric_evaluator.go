package quality

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/evaluation"
)

type SingleNumericMetricEvaluator interface {
	ChatConversationEvaluator
	MetricName() string
}

type BaseSingleNumericMetricEvaluator struct {
	BaseChatConversationEvaluator
}

func (e *BaseSingleNumericMetricEvaluator) EvaluationMetricNames(evaluator SingleNumericMetricEvaluator) []string {
	return []string{evaluator.MetricName()}
}

func (e *BaseSingleNumericMetricEvaluator) SystemPrompt() string {
	return `
	You are an AI assistant. You will be given the definition of an evaluation metric for assessing the quality of
	a response in a question-answering task. Your job is to compute an accurate evaluation score for the provided
	evaluation metric based on the provided scoring guidance.

	This evaluation score should always be an integer between 1 and 5. So your response should be 1 or 2 or 3 or 4
	or 5.

	Your response should be a single character containing only the evaluation score. Do not include any other text
	in your response besides the evaluation score.
	`
}

func (e *BaseSingleNumericMetricEvaluator) InitializeResult(evaluator SingleNumericMetricEvaluator) *evaluation.EvaluationResult {
	metric := evaluation.NewNumericMetric(evaluator.MetricName(), 0, nil)
	return evaluation.NewEvaluationResult([]evaluation.IEvaluationMetric{&metric})
}

func (e *BaseSingleNumericMetricEvaluator) PerformEvaluation(evaluator SingleNumericMetricEvaluator, ctx context.Context, chatConfiguration *evaluation.ChatConfiguration, evaluationMessages []chatcompletion.ChatMessage, result *evaluation.EvaluationResult) error {
	metricName := evaluator.MetricName()

	startTime := time.Now()
	if metric, ok := result.Metrics[metricName].(*evaluation.NumericMetric); ok {
		f := chatcompletion.TextFormat
		chatOptions := &chatcompletion.ChatOptions{
			MaxOutputTokens:  toPtr(int64(1)),
			Temperature:      toPtr(float64(0)),
			TopP:             toPtr(float64(1)),
			PresencePenalty:  toPtr(float64(0)),
			FrequencyPenalty: toPtr(float64(0)),
			ResponseFormat:   &f,
		}
		evaluationResponse, err := chatConfiguration.ChatClient.GetResponse(ctx, evaluationMessages, chatOptions)
		if err != nil {
			duration := fmt.Sprintf("%f", time.Since(startTime).Seconds())
			metric.AddOrUpdateMetadata("evaluation-duration", duration)
			return err
		}

		if evaluationResponse.ModelId != nil {
			metric.AddOrUpdateMetadata("evaluation-model-used", *evaluationResponse.ModelId)
		}

		if evaluationResponse.Usage != nil {
			if evaluationResponse.Usage.InputTokenCount != nil {
				metric.AddOrUpdateMetadata("evaluation-input-tokens-used", fmt.Sprintf("%d", *evaluationResponse.Usage.InputTokenCount))
			}

			if evaluationResponse.Usage.OutputTokenCount != nil {
				metric.AddOrUpdateMetadata("evaluation-output-tokens-used", fmt.Sprintf("%d", *evaluationResponse.Usage.OutputTokenCount))
			}

			if evaluationResponse.Usage.TotalTokenCount != nil {
				metric.AddOrUpdateMetadata("evaluation-total-tokens-used", fmt.Sprintf("%d", *evaluationResponse.Usage.TotalTokenCount))
			}
		}
		evaluationResponseText := evaluationResponse.Text()

		if len(evaluationResponseText) == 0 {
			metric.AddDiagnostics([]evaluation.EvaluationDiagnostic{
				evaluation.EvaluationDiagnosticError("Evaluation failed because the model failed to produce a valid evaluation response.")},
			)

		} else if score, err := strconv.Atoi(evaluationResponseText); err != nil {
			metric.Value = float64(score)
		} else {
			metric.AddDiagnostics([]evaluation.EvaluationDiagnostic{
				evaluation.EvaluationDiagnosticError(
					fmt.Sprintf("Failed to parse '%s' as an integer score for '%s'.", evaluationResponseText, metricName),
				)},
			)
		}

		interpretation := metric.InterpretScore()
		metric.SetInterpretation(&interpretation)
		duration := fmt.Sprintf("%f", time.Since(startTime).Seconds())
		metric.AddOrUpdateMetadata("evaluation-duration", duration)
		return nil
	}

	return fmt.Errorf("metric %s not found in result", metricName)
}

func toPtr[T comparable](value T) *T {
	return &value
}
