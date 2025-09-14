package chroma

import (
	"context"
	"fmt"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
)

const (
	// defaultTopK is the default number of results to return.
	defaultTopK = 5
)

// RetrieverConfig holds the configuration for the ChromaDB retriever.
type RetrieverConfig struct {
	// Client is the ChromaDB client.
	client chroma_go.Client
	// CollectionName is the name of the ChromaDB collection to search.
	CollectionName string
	// TopK is the number of results to return. Default 5.
	TopK int
	// Embedding is the method to vectorize the query.
	Embedding embedding.Embedder
}

// Retriever implements the retriever.Retriever interface for ChromaDB.
type Retriever struct {
	config *RetrieverConfig
	// col is the ChromaDB collection object, initialized in NewRetriever.
	col chroma_go.Collection
}

// NewRetriever creates a new ChromaDB-based Retriever.
func NewRetriever(ctx context.Context, config *RetrieverConfig) (*Retriever, error) {
	if config.Embedding == nil {
		return nil, fmt.Errorf("[NewRetriever] embedding not provided for chroma retriever")
	}
	if config.CollectionName == "" {
		return nil, fmt.Errorf("[NewRetriever] collection name not provided")
	}

	if config.TopK == 0 {
		config.TopK = defaultTopK
	}

	client, err := chroma_go.NewHTTPClient(
		chroma_go.WithBaseURL(os.Getenv("chroma_base_url")),
		chroma_go.WithDatabaseAndTenant(os.Getenv("chroma_database"), os.Getenv("chroma_tenant")),
		chroma_go.WithAuth(chroma_go.NewBasicAuthCredentialsProvider(os.Getenv("chroma_user"), os.Getenv("chroma_password"))),
	)
	if err != nil {
		return nil, err
	}

	config.client = client
	emb := &EinoEmbedderAdapter{embedder: config.Embedding}
	col, err := config.client.GetOrCreateCollection(
		ctx,
		config.CollectionName,
		chroma_go.WithEmbeddingFunctionCreate(emb),
		chroma_go.WithIfNotExistsCreate(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create collection: %w", err)
	}

	return &Retriever{
		config: config,
		col:    col,
	}, nil
}

// Retrieve finds documents similar to the given query using ChromaDB.
func (r *Retriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) (docs []*schema.Document, err error) {
	co := retriever.GetCommonOptions(&retriever.Options{
		TopK: &r.config.TopK,
	}, opts...)
	io := retriever.GetImplSpecificOptions(&implOptions{}, opts...)

	ctx = callbacks.EnsureRunInfo(ctx, r.GetType(), components.ComponentOfRetriever)
	ctx = callbacks.OnStart(ctx, &retriever.CallbackInput{
		Query:          query,
		TopK:           *co.TopK,
		ScoreThreshold: co.ScoreThreshold,
	})
	defer func() {
		if err != nil {
			callbacks.OnError(ctx, err)
		}
	}()

	result, err := r.col.Query(ctx,
		chroma_go.WithQueryTexts(query),
		chroma_go.WithNResults(int(*co.TopK)),
		chroma_go.WithIncludeQuery(chroma_go.IncludeMetadatas, chroma_go.IncludeDocuments, chroma_go.IncludeEmbeddings),
		chroma_go.WithWhereQuery(io.FilterQuery))
	if err != nil {
		return nil, err
	}

	// The Query method returns a list of results for each query. We only have one.
	if len(result.GetIDGroups()) == 0 {
		callbacks.OnEnd(ctx, &retriever.CallbackOutput{Docs: docs})
		return docs, nil
	}

	// Parse the results into eino schema.
	for i := range result.GetIDGroups()[0] {
		doc := &schema.Document{
			ID:       (string)(result.GetIDGroups()[0][i]),
			Content:  "",
			MetaData: make(map[string]any),
		}
		if impl, ok := result.GetMetadatasGroups()[0][i].(*chroma_go.DocumentMetadataImpl); ok {
			keys := impl.Keys()
			for _, key := range keys {
				if v, ok := impl.GetRaw(key); ok {
					doc.MetaData[key] = v
				}
			}
		}
		// The documents field is a list of strings.
		if len(result.GetDocumentsGroups()[0]) > i {
			doc.Content = result.GetDocumentsGroups()[0][i].ContentString()
		}

		docs = append(docs, doc)
	}

	callbacks.OnEnd(ctx, &retriever.CallbackOutput{Docs: docs})

	return docs, nil
}

// makeEmbeddingCtx creates a new context for embedding.
func (r *Retriever) makeEmbeddingCtx(ctx context.Context, emb embedding.Embedder) context.Context {
	runInfo := &callbacks.RunInfo{
		Component: components.ComponentOfEmbedding,
	}
	if embType, ok := components.GetType(emb); ok {
		runInfo.Type = embType
	}
	runInfo.Name = runInfo.Type + string(runInfo.Component)
	return callbacks.ReuseHandlers(ctx, runInfo)
}

const typ = "chroma"

// GetType returns the type of the retriever.
func (r *Retriever) GetType() string {
	return typ
}

// IsCallbacksEnabled checks if callbacks are enabled.
func (r *Retriever) IsCallbacksEnabled() bool {
	return true
}

type implOptions struct {
	FilterQuery chroma_go.WhereFilter
}

// WithFilter provides a filter query for ChromaDB.
// The filter must be a valid map[string]any, e.g., map[string]any{"$eq": map[string]any{"category": "tech"}}.
func WithFilter(filter chroma_go.WhereFilter) retriever.Option {
	return retriever.WrapImplSpecificOptFn(func(o *implOptions) {
		o.FilterQuery = filter
	})
}
