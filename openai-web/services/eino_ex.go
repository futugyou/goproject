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

	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/futugyousuzu/go-openai-web/models"
)

type EinoService struct {
	db    *mongo.Database
	chain *gemini.ChatModel
}

func (e *EinoService) Init(ctx context.Context, client *mongo.Client, chain *gemini.ChatModel) *EinoService {
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
		chain, err = getGeminiModel(ctx)
		if err != nil {
			return e
		}

		e.chain = chain
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

func (e *EinoService) SaveConversation(ctx context.Context, model ConversationModel) error {
	db := e.db
	coll := db.Collection("eino_conversation")
	opt := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "id", Value: model.Id}}

	if _, err := coll.UpdateOne(ctx, filter, bson.M{
		"$set": model,
	}, opt); err != nil {
		return err
	}

	return nil
}

func (e *EinoService) GetConversation(ctx context.Context, id string) (*ConversationModel, error) {
	model := &ConversationModel{}
	db := e.db
	coll := db.Collection("eino_conversation")
	filter := bson.D{{Key: "id", Value: id}}
	opts := &options.FindOneOptions{}
	if err := coll.FindOne(ctx, filter, opts).Decode(model); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, fmt.Errorf("can not find conversation with id: %s", id)
		} else {
			return nil, err
		}
	}

	return model, nil
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
	// TODO：Need to fill slowly
	for _, node := range model.Nodes {
		switch node.Type {
		case "branch":
		case "model":
		case "template":
		case "doc":
		case "embed":
		case "graph":
		case "indexer":
		case "lambda":
		case "loader":
		case "passthrough":
		case "retriever":
		case "tools":
		}
	}

	for _, edge := range model.Edges {
		g.AddEdge(edge.Source, edge.Target)
	}

	starts, ends := findStartAndEnd(model.Edges)
	for _, edge := range starts {
		g.AddEdge(compose.START, edge)
	}

	for _, edge := range ends {
		g.AddEdge(edge, compose.END)
	}

	r, err := g.Compile(ctx, compose.WithMaxRunSteps(10))
	if err != nil {
		return nil, err
	}

	return r.Invoke(ctx, input)
}

// Currently there can only be one start and one end
func findStartAndEnd(edges []models.Edge) (starts []string, ends []string) {
	sourceSet := make(map[string]struct{})
	targetSet := make(map[string]struct{})

	for _, e := range edges {
		sourceSet[e.Source] = struct{}{}
		targetSet[e.Target] = struct{}{}
	}

	// start node: appears in sourceSet but not in targetSet
	for s := range sourceSet {
		if _, ok := targetSet[s]; !ok {
			starts = append(starts, s)
		}
	}

	// end node: appears in targetSet but not in sourceSet
	for t := range targetSet {
		if _, ok := sourceSet[t]; !ok {
			ends = append(ends, t)
		}
	}

	return
}

func getGeminiModel(ctx context.Context) (*gemini.ChatModel, error) {
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
