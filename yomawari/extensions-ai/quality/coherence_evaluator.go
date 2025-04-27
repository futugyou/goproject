package quality

import (
	"context"
	"fmt"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions-ai/evaluation"
)

var _ SingleNumericMetricEvaluator = (*CoherenceEvaluator)(nil)

type CoherenceEvaluator struct {
	BaseSingleNumericMetricEvaluator
}

// Evaluate implements SingleNumericMetricEvaluator.
// Subtle: this method shadows the method (BaseSingleNumericMetricEvaluator).Evaluate of CoherenceEvaluator.BaseSingleNumericMetricEvaluator.
func (c *CoherenceEvaluator) Evaluate(ctx context.Context, messages []chatcompletion.ChatMessage, modelResponse chatcompletion.ChatResponse, chatConfiguration *evaluation.ChatConfiguration, additionalContext []evaluation.EvaluationContext) (*evaluation.EvaluationResult, error) {
	return c.BaseSingleNumericMetricEvaluator.Evaluate(c, ctx, messages, modelResponse, chatConfiguration, additionalContext)
}

// PerformEvaluation implements SingleNumericMetricEvaluator.
// Subtle: this method shadows the method (BaseSingleNumericMetricEvaluator).PerformEvaluation of CoherenceEvaluator.BaseSingleNumericMetricEvaluator.
func (c *CoherenceEvaluator) PerformEvaluation(ctx context.Context, chatConfiguration *evaluation.ChatConfiguration, evaluationMessages []chatcompletion.ChatMessage, result *evaluation.EvaluationResult) error {
	return c.BaseSingleNumericMetricEvaluator.PerformEvaluation(c, ctx, chatConfiguration, evaluationMessages, result)
}

// EvaluationMetricNames implements SingleNumericMetricEvaluator.
// Subtle: this method shadows the method (BaseSingleNumericMetricEvaluator).EvaluationMetricNames of CoherenceEvaluator.BaseSingleNumericMetricEvaluator.
func (c *CoherenceEvaluator) EvaluationMetricNames() []string {
	return c.BaseSingleNumericMetricEvaluator.EvaluationMetricNames(c)
}

// IgnoresHistory implements SingleNumericMetricEvaluator.
func (c *CoherenceEvaluator) IgnoresHistory() bool {
	return true
}

// InitializeResult implements SingleNumericMetricEvaluator.
// Subtle: this method shadows the method (BaseSingleNumericMetricEvaluator).InitializeResult of CoherenceEvaluator.BaseSingleNumericMetricEvaluator.
func (c *CoherenceEvaluator) InitializeResult() *evaluation.EvaluationResult {
	return c.BaseSingleNumericMetricEvaluator.InitializeResult(c)
}

// MetricName implements SingleNumericMetricEvaluator.
func (c *CoherenceEvaluator) MetricName() string {
	return "Coherence"
}

// RenderEvaluationPrompt implements SingleNumericMetricEvaluator.
func (c *CoherenceEvaluator) RenderEvaluationPrompt(ctx context.Context, userRequest *chatcompletion.ChatMessage, modelResponse chatcompletion.ChatResponse, conversationHistory []chatcompletion.ChatMessage, additionalContext []evaluation.EvaluationContext) (*string, error) {
	renderedModelResponse := c.RenderChatResponse(ctx, &modelResponse)
	renderedUserRequest := ""
	if userRequest != nil {
		renderedUserRequest = c.RenderChatMessage(ctx, userRequest)
	}

	prompt := fmt.Sprintf(`
		Coherence of an answer is measured by how well all the sentences fit together and sound naturally as a
		whole. Consider the overall quality of the answer when evaluating coherence.

		Given the question and answer, score the coherence of the answer between one to five stars using the
		following rating scale:
		One star: the answer completely lacks coherence
		Two stars: the answer mostly lacks coherence
		Three stars: the answer is partially coherent
		Four stars: the answer is mostly coherent
		Five stars: the answer has perfect coherency

		The rating value should always be an integer between 1 and 5. So the rating produced should be 1 or 2 or 3
		or 4 or 5.

		question: What is your favorite indoor activity and why do you enjoy it?
		answer: I like pizza. The sun is shining.
		stars: 1

		question: Can you describe your favorite movie without giving away any spoilers?
		answer: It is a science fiction movie. There are dinosaurs. The actors eat cake. People must stop the
		villain.
		stars: 2

		question: What are some benefits of regular exercise?
		answer: Regular exercise improves your mood. A good workout also helps you sleep better. Trees are green.
		stars: 3

		question: How do you cope with stress in your daily life?
		answer: I usually go for a walk to clear my head. Listening to music helps me relax as well. Stress is a
		part of life, but we can manage it through some activities.
		stars: 4

		question: What can you tell me about climate change and its effects on the environment?
		answer: Climate change has far-reaching effects on the environment. Rising temperatures result in the
		melting of polar ice caps, contributing to sea-level rise. Additionally, more frequent and severe weather
		events, such as hurricanes and heatwaves, can cause disruption to ecosystems and human societies alike.
		stars: 5

		question: %s
		answer: %s
		stars:
		`, renderedUserRequest, renderedModelResponse)

	return &prompt, nil
}
