package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/genai"

	embedding "github.com/cloudwego/eino-ext/components/embedding/gemini"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/futugyousuzu/go-openai-web/eino_extensions/graph"
	"github.com/futugyousuzu/go-openai-web/models"
)

type EinoService struct {
	db    *mongo.Database
	chain *gemini.ChatModel
	embed *embedding.Embedder
}

func (e *EinoService) Init(ctx context.Context, client *mongo.Client, chain *gemini.ChatModel, embed *embedding.Embedder) *EinoService {
	var err error
	if client == nil {
		db_name := os.Getenv("db_name")
		uri := os.Getenv("mongodb_url")
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			return e
		}

		e.db = client.Database(db_name)
	}

	if chain == nil {
		chain, err = e.getGeminiModel(ctx)
		if err != nil {
			return e
		}

		e.chain = chain
	}

	if embed == nil {
		embed, err = e.getGeminiEmbeddingModel(ctx)
		if err != nil {
			return e
		}

		e.embed = embed
	}

	return e
}

type ConversationModel struct {
	Id            string    `json:"id,omitempty" bson:"id,omitempty"`
	SystemMessage string    `json:"system_message,omitempty" bson:"system_message,omitempty"`
	UserMessage   string    `json:"user_message,omitempty" bson:"user_message,omitempty"`
	History       string    `json:"history,omitempty" bson:"history,omitempty"`
	LastUpdated   time.Time `json:"last_updated,omitempty" bson:"last_updated,omitempty"`
}

func (e *EinoService) saveModel(ctx context.Context, collection string, idKey string, idValue any, model any) error {
	coll := e.db.Collection(collection)
	opt := options.Update().SetUpsert(true)
	filter := bson.D{{Key: idKey, Value: idValue}}

	_, err := coll.UpdateOne(ctx, filter, bson.M{"$set": model}, opt)
	return err
}

func (e *EinoService) SaveConversation(ctx context.Context, model ConversationModel) error {
	return e.saveModel(ctx, "eino_conversation", "id", model.Id, model)
}

func (e *EinoService) SaveFlowGraph(ctx context.Context, model models.FlowGraph) error {
	return e.saveModel(ctx, "eino_flowgraph", "id", model.ID, model)
}

func getModel[T any](ctx context.Context, db *mongo.Database, collection string, id string) (*T, error) {
	coll := db.Collection(collection)
	filter := bson.D{{Key: "id", Value: id}}

	var model T
	if err := coll.FindOne(ctx, filter).Decode(&model); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("can not find %s with id: %s", collection, id)
		}
		return nil, err
	}
	return &model, nil
}

func (e *EinoService) GetConversation(ctx context.Context, id string) (*ConversationModel, error) {
	return getModel[ConversationModel](ctx, e.db, "eino_conversation", id)
}

func (e *EinoService) GetFlowGraph(ctx context.Context, id string) (*models.FlowGraph, error) {
	return getModel[models.FlowGraph](ctx, e.db, "eino_flowgraph", id)
}

func (e *EinoService) GeneralLLMRunner(ctx context.Context, userMsg string, systemMsg *string, useHistory bool, vs map[string]any) (*schema.Message, error) {
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
	return e.chain.Generate(ctx, messages)
}

func (e *EinoService) GraphRunner(ctx context.Context, model models.FlowGraph, input map[string]any) (*schema.Message, error) {
	g := compose.NewGraph[map[string]any, *schema.Message]()

	if err := graph.AddNodesToGraph(ctx, g, model.Nodes, e.embed, e.chain); err != nil {
		return nil, err
	}

	graph.AddEdgesToGraph(g, model.Edges)

	r, err := g.Compile(ctx, compose.WithMaxRunSteps(10))
	if err != nil {
		return nil, err
	}

	return r.Invoke(ctx, input)
}

func (e *EinoService) getGeminiModel(ctx context.Context) (*gemini.ChatModel, error) {
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

func (e *EinoService) getGeminiEmbeddingModel(ctx context.Context) (*embedding.Embedder, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	modelid := os.Getenv("GEMINI_EMBEDDING_MODEL_ID")
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return embedding.NewEmbedder(ctx, &embedding.EmbeddingConfig{
		Client:   client,
		Model:    modelid,
		TaskType: "RETRIEVAL_QUERY",
	})
}
