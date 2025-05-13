package evaluation

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

type IEvaluator interface {
	EvaluationMetricNames() []string
	Evaluate(ctx context.Context, messages []chatcompletion.ChatMessage, modelResponse chatcompletion.ChatResponse, chatConfiguration *ChatConfiguration, additionalContext []EvaluationContext) (*EvaluationResult, error)
}
