package evaluation

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

var _ IEvaluator = (*CompositeEvaluator)(nil)

type CompositeEvaluator struct {
	evaluationMetricNames []string
	evaluators            []IEvaluator
}

func NewCompositeEvaluator(evaluators []IEvaluator) *CompositeEvaluator {
	metricNames := map[string]struct{}{}
	for _, evaluator := range evaluators {
		for _, metricName := range evaluator.EvaluationMetricNames() {
			if _, ok := metricNames[metricName]; !ok {
				metricNames[metricName] = struct{}{}
			}
		}
	}
	c := &CompositeEvaluator{
		evaluationMetricNames: make([]string, 0),
		evaluators:            evaluators,
	}

	for metricName := range metricNames {
		c.evaluationMetricNames = append(c.evaluationMetricNames, metricName)
	}
	return c
}

// Evaluate implements IEvaluator.
func (c *CompositeEvaluator) Evaluate(ctx context.Context, messages []chatcompletion.ChatMessage, modelResponse chatcompletion.ChatResponse, chatConfiguration *ChatConfiguration, additionalContext []EvaluationContext) (*EvaluationResult, error) {
	metrics := []IEvaluationMetric{}
	resultsStream := c.evaluateAndStreamResults(ctx, messages, modelResponse, chatConfiguration, additionalContext)

	for result := range resultsStream {
		for _, metric := range result.Metrics {
			metrics = append(metrics, metric)
		}
	}

	return NewEvaluationResult(metrics), nil
}

// EvaluationMetricNames implements IEvaluator.
func (c *CompositeEvaluator) EvaluationMetricNames() []string {
	return c.evaluationMetricNames
}

func (c *CompositeEvaluator) evaluateAndStreamResults(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	modelResponse chatcompletion.ChatResponse,
	chatConfiguration *ChatConfiguration,
	additionalContext []EvaluationContext,
) chan EvaluationResult {
	result := make(chan EvaluationResult)

	go func() {
		defer close(result)
		for _, evaluator := range c.evaluators {
			if res, err := evaluator.Evaluate(ctx, messages, modelResponse, chatConfiguration, additionalContext); err != nil {
				message := err.Error()
				res := &EvaluationResult{}
				for _, metricName := range evaluator.EvaluationMetricNames() {
					metric := NewEvaluationMetric(metricName, nil)
					(&metric).AddDiagnostics([]EvaluationDiagnostic{EvaluationDiagnosticError(message)})
					res.Metrics[metricName] = &metric
				}
			} else {
				result <- *res
			}
		}
	}()
	return result
}
