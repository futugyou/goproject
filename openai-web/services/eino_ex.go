package services

import (
	"context"
	"os"

	"google.golang.org/genai"

	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func GeneralLLMRunner(ctx context.Context, chatModel *gemini.ChatModel, userMsg string, systemMsg *string, useHistory bool, vs map[string]any) (*schema.Message, error) {
	templates := []schema.MessagesTemplate{}
	if systemMsg != nil {
		templates = append(templates, schema.SystemMessage(*systemMsg))
	}
	if useHistory {
		templates = append(templates, schema.MessagesPlaceholder("chat_history", true))
	}
	templates = append(templates, schema.UserMessage(userMsg))
	template := prompt.FromMessages(schema.FString, templates...)
	messages, err := template.Format(ctx, vs)
	if err != nil {
		return nil, err
	}
	return chatModel.Generate(ctx, messages)
}

func GetGeminiModel(ctx context.Context) (*gemini.ChatModel, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	modelid := os.Getenv("GEMINI_MODEL_ID")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return gemini.NewChatModel(ctx, &gemini.Config{
		Client: client,
		Model:  modelid,
		ThinkingConfig: &genai.ThinkingConfig{
			IncludeThoughts: true,
			ThinkingBudget:  nil,
		},
	})
}
