package chroma

import (
	"context"
	"fmt"

	chroma_go "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
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
	Client *chroma_go.Client
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
	col *chroma_go.Collection
}

// NewRetriever creates a new ChromaDB-based Retriever.
func NewRetriever(ctx context.Context, config *RetrieverConfig) (*Retriever, error) {
	if config.Embedding == nil {
		return nil, fmt.Errorf("[NewRetriever] embedding not provided for chroma retriever")
	}

	if config.Client == nil {
		return nil, fmt.Errorf("[NewRetriever] chroma client not provided")
	}

	if config.CollectionName == "" {
		return nil, fmt.Errorf("[NewRetriever] collection name not provided")
	}

	if config.TopK == 0 {
		config.TopK = defaultTopK
	}

	// Get the collection from the client.
	col, err := config.Client.GetCollection(ctx, config.CollectionName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection %s: %w", config.CollectionName, err)
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
		Query:  query,
		TopK:   *co.TopK,
		ScoreThreshold: co.ScoreThreshold,
	})
	defer func() {
		if err != nil {
			callbacks.OnError(ctx, err)
		}
	}()

	// The query text to search for, passed as a slice.
	queryTexts := []string{query}

	// The number of results to return.
	nResults := int32(*co.TopK)

	// The metadata filter, which is an optional map. Pass nil if no filter is provided.
	var whereFilter map[string]interface{}
	if io.FilterQuery != nil {
		whereFilter = io.FilterQuery
	}

	// Document content filter, which is not used in this retriever.
	whereDocuments := make(map[string]interface{})

	// What data to include in the response.
	include := []types.QueryEnum{types.IDocuments, types.IMetadatas, types.IDistances}

	// Call the Query method with the correct parameters.
	result, err := r.col.Query(
		ctx,
		queryTexts,
		nResults,
		whereFilter,
		whereDocuments,
		include,
	)
	if err != nil {
		return nil, err
	}

	// The Query method returns a list of results for each query. We only have one.
	if len(result.Ids) == 0 || len(result.Ids[0]) == 0 {
		callbacks.OnEnd(ctx, &retriever.CallbackOutput{Docs: docs})
		return docs, nil
	}

	// Parse the results into eino schema.
	for i := range result.Ids[0] {
		doc := &schema.Document{
			ID:       result.Ids[0][i],
			Content:  "",
			MetaData: result.Metadatas[0][i],
		}

		// The documents field is a list of strings.
		if len(result.Documents[0]) > i {
			doc.Content = result.Documents[0][i]
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
	FilterQuery map[string]any
}
 
// WithFilter provides a filter query for ChromaDB.
// The filter must be a valid map[string]any, e.g., map[string]any{"$eq": map[string]any{"category": "tech"}}.
func WithFilter(filter map[string]any) retriever.Option {
	return retriever.WrapImplSpecificOptFn(func(o *implOptions) {
		o.FilterQuery = filter
	})
}