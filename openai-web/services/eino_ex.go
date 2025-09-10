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

	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino-ext/components/document/parser/docx"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino-ext/components/document/parser/xlsx"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/semantic"
	embedding "github.com/cloudwego/eino-ext/components/embedding/gemini"
	"github.com/cloudwego/eino-ext/components/model/gemini"
	"github.com/cloudwego/eino-ext/components/tool/googlesearch"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/components/document/parser"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/futugyousuzu/go-openai-web/models"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
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
	// TODO: Need some built-in Functional Node
	for _, node := range model.Nodes {
		switch node.Type {
		case "branch":
			n, err := e.getGraphBranch(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddBranch(node.ID, n)
		case "model":
			g.AddChatModelNode(node.ID, e.chain)
		case "template":
			n, err := e.getChatTemplateNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddChatTemplateNode(node.ID, n)
		case "doc":
			n, err := e.getDocumentTransformerNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddDocumentTransformerNode(node.ID, n)
		case "embed":
			g.AddEmbeddingNode(node.ID, e.embed)
		case "graph":
			n, err := e.getGraphNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddGraphNode(node.ID, n)
		case "indexer":
			n, err := e.getIndexerNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddIndexerNode(node.ID, n)
		case "lambda":
			n, err := e.getLambdaNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddLambdaNode(node.ID, n)
		case "loader":
			n, err := e.getLoaderNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddLoaderNode(node.ID, n)
		case "passthrough":
			g.AddPassthroughNode(node.ID)
		case "retriever":
			n, err := e.getRetrieverNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddRetrieverNode(node.ID, n)
		case "tools":
			n, err := e.getToolsNode(ctx, node)
			if err != nil {
				return nil, err
			}
			g.AddToolsNode(node.ID, n)
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

func (e *EinoService) getToolsNode(ctx context.Context, node models.Node) (*compose.ToolsNode, error) {
	tools := []tool.BaseTool{}
	if googletool, ok := node.Data["googletool"].(string); ok && len(googletool) > 0 {
		googleTool, err := googlesearch.NewTool(ctx, &googlesearch.Config{
			APIKey:         os.Getenv("GOOGLE_API_KEY"),
			SearchEngineID: os.Getenv("GOOGLE_SEARCH_ENGINE_ID"),
		})
		if err != nil {
			return nil, err
		}
		tools = append(tools, googleTool)
	}

	if mcptoolurl, ok := node.Data["mcptoolurl"].(string); ok && len(mcptoolurl) > 0 {
		cli, err := client.NewSSEMCPClient(mcptoolurl)
		if err != nil {
			return nil, err
		}
		err = cli.Start(ctx)
		if err != nil {
			return nil, err
		}

		initRequest := mcp.InitializeRequest{}
		initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initRequest.Params.ClientInfo = mcp.Implementation{
			Name:    "enio-client",
			Version: "1.0.0",
		}

		_, err = cli.Initialize(ctx, initRequest)
		if err != nil {
			return nil, err
		}

		mcpTools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
		tools = append(tools, mcpTools...)
	}

	if len(tools) > 0 {
		return compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
			Tools: tools,
		})
	}

	return nil, fmt.Errorf("invalid tool node: %s", node.ID)
}

func (e *EinoService) getRetrieverNode(ctx context.Context, node models.Node) (retriever.Retriever, error) {
	panic("unimplemented")
}

func (e *EinoService) getLoaderNode(ctx context.Context, node models.Node) (document.Loader, error) {
	if loader, ok := node.Data["loader"].(string); ok && len(loader) > 0 {
		var loaderParser parser.Parser = nil
		if p, ok := node.Data["parser"].(string); ok && len(p) > 0 {
			switch p {
			case "docx":
				loaderParser, _ = docx.NewDocxParser(ctx, &docx.Config{
					ToSections:      true,
					IncludeComments: true,
					IncludeHeaders:  true,
					IncludeFooters:  true,
					IncludeTables:   true,
				})
			case "pdf":
				loaderParser, _ = pdf.NewPDFParser(ctx, &pdf.Config{ToPages: true})
			case "xlsx":
				loaderParser, _ = xlsx.NewXlsxParser(ctx, &xlsx.Config{})
			}
		}

		switch loader {
		case "url":
			loaderConfig := &url.LoaderConfig{
				Parser: loaderParser,
			}

			return url.NewLoader(ctx, loaderConfig)
		}
	}

	return nil, fmt.Errorf("invalid document loader node: %s", node.ID)
}

func (e *EinoService) getLambdaNode(ctx context.Context, node models.Node) (*compose.Lambda, error) {
	panic("unimplemented")
}

func (e *EinoService) getGraphNode(ctx context.Context, node models.Node) (compose.AnyGraph, error) {
	panic("unimplemented")
}

func (e *EinoService) getDocumentTransformerNode(ctx context.Context, node models.Node) (document.Transformer, error) {
	if transformer, ok := node.Data["transformer"].(string); ok && len(transformer) > 0 {
		switch transformer {
		case "markdown":
			headers := map[string]string{
				"#":   "h1",
				"##":  "h2",
				"###": "h3",
			}
			if h, ok := node.Data["transformer_header"].(map[string]string); ok {
				headers = h
			}
			return markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{Headers: headers, TrimHeaders: false})
		case "semantic":
			return semantic.NewSplitter(ctx, &semantic.Config{
				Embedding: e.embed,
			})
		case "recursive":
			return recursive.NewSplitter(ctx, &recursive.Config{ChunkSize: 1000, OverlapSize: 200})
		}

	}

	return nil, fmt.Errorf("invalid document transformer node: %s", node.ID)
}

func (e *EinoService) getChatTemplateNode(ctx context.Context, node models.Node) (prompt.ChatTemplate, error) {
	role := schema.System
	if r, ok := node.Data["role"].(schema.RoleType); ok && len(r) > 0 {
		role = r
	}
	if content, ok := node.Data["content"].(string); ok && len(content) > 0 {
		return prompt.FromMessages(schema.FString, &schema.Message{
			Role:    role,
			Content: content,
		}), nil
	}

	return nil, fmt.Errorf("invalid chat template node: %s", node.ID)
}

func (e *EinoService) getIndexerNode(ctx context.Context, node models.Node) (indexer.Indexer, error) {
	panic("unimplemented")
}

func (e *EinoService) getGraphBranch(ctx context.Context, node models.Node) (*compose.GraphBranch, error) {
	return compose.NewGraphBranch(func(ctx context.Context, in map[string]any) (string, error) {
		return "", nil
	}, map[string]bool{}), nil
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
