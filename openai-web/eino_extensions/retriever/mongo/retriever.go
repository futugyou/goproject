package mongodb

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	// defaultVectorSearchIndex is the default name for the vector search index.
	// You need to create this index in MongoDB.
	defaultVectorSearchIndex = "vector_index"
	// defaultVectorField is the default field name for the vector content.
	defaultVectorField = "vector_content"
	// defaultContentField is the default field name for the document content.
	defaultContentField = "content"
	// defaultTopK is the default number of results to return.
	defaultTopK = 5
)

// RetrieverConfig holds the configuration for the MongoDB retriever.
type RetrieverConfig struct {
	// Client is the MongoDB client.
	Client *mongo.Client
	// Database is the name of the MongoDB database.
	Database string
	// Collection is the name of the MongoDB collection.
	Collection string
	// Index is the name of the vector search index. Default "vector_index".
	Index string
	// VectorField is the field name for the vector data. Default "vector_content".
	VectorField string
	// TopK is the number of results to return. Default 5.
	TopK int
	// Embedding is the method to vectorize the query.
	Embedding embedding.Embedder
}

// Retriever implements the retriever.Retriever interface for MongoDB.
type Retriever struct {
	config *RetrieverConfig
}

// NewRetriever creates a new MongoDB-based Retriever.
func NewRetriever(ctx context.Context, config *RetrieverConfig) (*Retriever, error) {
	if config.Embedding == nil {
		return nil, fmt.Errorf("[NewRetriever] embedding not provided for mongodb retriever")
	}

	if config.Client == nil {
		return nil, fmt.Errorf("[NewRetriever] mongo client not provided")
	}

	if config.Database == "" || config.Collection == "" {
		return nil, fmt.Errorf("[NewRetriever] database or collection name not provided")
	}

	if config.Index == "" {
		config.Index = defaultVectorSearchIndex
	}

	if config.VectorField == "" {
		config.VectorField = defaultVectorField
	}

	if config.TopK == 0 {
		config.TopK = defaultTopK
	}

	return &Retriever{
		config: config,
	}, nil
}

// Retrieve finds documents similar to the given query using MongoDB's vector search.
func (r *Retriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) (docs []*schema.Document, err error) {
	co := retriever.GetCommonOptions(&retriever.Options{
		Index:     &r.config.Index,
		TopK:      &r.config.TopK,
		Embedding: r.config.Embedding,
	}, opts...)
	io := retriever.GetImplSpecificOptions(&implOptions{}, opts...)

	ctx = callbacks.EnsureRunInfo(ctx, r.GetType(), components.ComponentOfRetriever)
	ctx = callbacks.OnStart(ctx, &retriever.CallbackInput{
		Query: query,
		TopK:  *co.TopK,
	})
	defer func() {
		if err != nil {
			callbacks.OnError(ctx, err)
		}
	}()

	emb := co.Embedding
	if emb == nil {
		return nil, fmt.Errorf("[mongodb retriever] embedding not provided")
	}

	vectors, err := emb.EmbedStrings(r.makeEmbeddingCtx(ctx, emb), []string{query})
	if err != nil {
		return nil, err
	}

	if len(vectors) != 1 {
		return nil, fmt.Errorf("[mongodb retriever] invalid return length of vector, got=%d, expected=1", len(vectors))
	}

	collection := r.config.Client.Database(r.config.Database).Collection(r.config.Collection)

	// MongoDB's vector search is done via an aggregation pipeline.
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$vectorSearch", Value: bson.D{
				{Key: "index", Value: *co.Index},
				{Key: "path", Value: r.config.VectorField},
				{Key: "queryVector", Value: vectors[0]},
				{Key: "numCandidates", Value: *co.TopK * 10},
				{Key: "limit", Value: *co.TopK},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: defaultContentField, Value: 1},
				{Key: r.config.VectorField, Value: 1},
				{Key: "score", Value: bson.D{{Key: "$meta", Value: "vectorSearchScore"}}},
			}},
		},
	}

	if io.FilterQuery != "" {
		// If a filter is provided, add a $match stage to the pipeline.
		// The filter query must be a valid BSON document.
		var filterDoc bson.D
		if err := bson.UnmarshalExtJSON([]byte(io.FilterQuery), false, &filterDoc); err != nil {
			return nil, fmt.Errorf("invalid filter query format: %w", err)
		}

		pipeline = append(mongo.Pipeline{
			bson.D{
				{Key: "$match", Value: filterDoc},
			},
		}, pipeline...)
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result map[string]any
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		doc := &schema.Document{
			ID:       fmt.Sprintf("%v", result["_id"]),
			MetaData: make(map[string]any),
		}
		// Extract content and vector from the result map.
		if content, ok := result[defaultContentField].(string); ok {
			doc.Content = content
		}
		if vector, ok := result[r.config.VectorField].(bson.A); ok {
			var floatVector []float64
			for _, v := range vector {
				if f, ok := v.(float64); ok { // BSON floats are float64.
					floatVector = append(floatVector, f)
				}
			}
			doc.WithDenseVector(floatVector)
		}

		// Copy any other fields to metadata.
		for k, v := range result {
			if k != "_id" && k != defaultContentField && k != r.config.VectorField {
				doc.MetaData[k] = v
			}
		}

		docs = append(docs, doc)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
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

// GetType returns the type of the retriever.
func (r *Retriever) GetType() string {
	return "mongo"
}

// IsCallbacksEnabled checks if callbacks are enabled.
func (r *Retriever) IsCallbacksEnabled() bool {
	return true
}

// implOptions contains implementation-specific options for the MongoDB retriever.
type implOptions struct {
	FilterQuery string
}

// WithFilterQuery redis filter query.
// see: https://redis.io/docs/latest/develop/interact/search-and-query/advanced-concepts/vectors/#filters
func WithFilterQuery(filter string) retriever.Option {
	return retriever.WrapImplSpecificOptFn(func(o *implOptions) {
		o.FilterQuery = filter
	})
}
