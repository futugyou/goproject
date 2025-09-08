package services

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"

	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type ChatHistoryModel struct {
	Id            string `json:"id,omitempty" bson:"id,omitempty"`
	SystemMessage string `json:"system_message,omitempty" bson:"system_message,omitempty"`
	UserMessage   string `json:"user_message,omitempty" bson:"user_message,omitempty"`
	History       string `json:"history,omitempty" bson:"history,omitempty"`
}

func SaveMessages(ctx context.Context, model ChatHistoryModel) error {
	db_name := os.Getenv("db_name")
	uri := os.Getenv("mongodb_url")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	db := client.Database(db_name)
	coll := db.Collection("eino_chat_history")
	opt := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: model.Id}}

	if _, err := coll.UpdateOne(ctx, filter, bson.M{
		"$set": model,
	}, opt); err != nil {
		return err
	}

	return nil
}

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
