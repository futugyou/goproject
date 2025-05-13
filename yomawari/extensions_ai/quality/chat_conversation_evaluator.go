package quality

import (
	"context"
	"fmt"
	"strings"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/evaluation"
)

type ChatConversationEvaluator interface {
	evaluation.IEvaluator
	IgnoresHistory() bool
	SystemPrompt() string
	RenderChatResponse(ctx context.Context, response *chatcompletion.ChatResponse) string
	RenderChatMessage(ctx context.Context, message *chatcompletion.ChatMessage) string
	RenderEvaluationPrompt(ctx context.Context, userRequest *chatcompletion.ChatMessage, modelResponse chatcompletion.ChatResponse, conversationHistory []chatcompletion.ChatMessage, additionalContext []evaluation.EvaluationContext) (*string, error)
	InitializeResult() *evaluation.EvaluationResult
	PerformEvaluation(ctx context.Context, chatConfiguration *evaluation.ChatConfiguration, evaluationMessages []chatcompletion.ChatMessage, result *evaluation.EvaluationResult) error
}

type BaseChatConversationEvaluator struct {
}

func (e *BaseChatConversationEvaluator) getUserRequestAndConversationHistory(evaluator ChatConversationEvaluator, messages []chatcompletion.ChatMessage) (*chatcompletion.ChatMessage, []chatcompletion.ChatMessage) {
	var userRequest *chatcompletion.ChatMessage
	var conversationHistory []chatcompletion.ChatMessage

	if evaluator.IgnoresHistory() {
		if len(messages) > 0 {
			lastMessage := messages[len(messages)-1]
			if lastMessage.Role == chatcompletion.RoleUser {
				userRequest = &lastMessage
			}
		}

		conversationHistory = []chatcompletion.ChatMessage{}
	} else {
		conversationHistory = messages
		lastMessageIndex := len(conversationHistory) - 1

		if lastMessageIndex >= 0 && conversationHistory[lastMessageIndex].Role == chatcompletion.RoleUser {
			userRequest = &conversationHistory[lastMessageIndex]
			conversationHistory = append(conversationHistory[:lastMessageIndex], conversationHistory[lastMessageIndex+1:]...)
		}
	}

	return userRequest, conversationHistory
}

func (e *BaseChatConversationEvaluator) Evaluate(evaluator ChatConversationEvaluator, ctx context.Context, messages []chatcompletion.ChatMessage, modelResponse chatcompletion.ChatResponse, chatConfiguration *evaluation.ChatConfiguration, additionalContext []evaluation.EvaluationContext) (*evaluation.EvaluationResult, error) {
	result := evaluator.InitializeResult()
	if len(modelResponse.Text()) == 0 {
		result.AddDiagnosticsToAllMetrics([]evaluation.EvaluationDiagnostic{
			evaluation.EvaluationDiagnosticError("Evaluation failed because the model response supplied for evaluation was null or empty."),
		})

		return result, nil
	}

	evaluationMessages := []chatcompletion.ChatMessage{}
	systemPrompt := evaluator.SystemPrompt()
	if len(systemPrompt) > 0 {
		evaluationMessages = append(evaluationMessages, *chatcompletion.NewChatMessageWithText(chatcompletion.RoleSystem, systemPrompt))
	}

	userRequest, conversationHistory := e.getUserRequestAndConversationHistory(evaluator, messages)
	evaluationPrompt, err := evaluator.RenderEvaluationPrompt(ctx, userRequest, modelResponse, conversationHistory, additionalContext)

	if err == nil {
		evaluationMessages = append(evaluationMessages, *chatcompletion.NewChatMessageWithText(chatcompletion.RoleUser, *evaluationPrompt))
	}

	evaluator.PerformEvaluation(ctx, chatConfiguration, evaluationMessages, result)

	return result, nil
}

func (e *BaseChatConversationEvaluator) RenderChatResponse(ctx context.Context, response *chatcompletion.ChatResponse) string {
	sb := strings.Builder{}
	for _, message := range response.Messages {
		msg := e.RenderChatMessage(ctx, &message)
		if len(msg) > 0 {
			sb.WriteString(msg)
		}

	}

	return sb.String()
}

func (e *BaseChatConversationEvaluator) RenderChatMessage(ctx context.Context, message *chatcompletion.ChatMessage) string {
	if message == nil {
		return ""
	}
	author := message.AuthorName
	role := message.Role
	content := message.Text()

	if author == nil {
		return fmt.Sprintf("[%s] %s\n", role, content)
	}

	return fmt.Sprintf("[%s] (%s)] %s\n", *author, role, content)
}
