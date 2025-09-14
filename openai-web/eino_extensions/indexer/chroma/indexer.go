package chroma

import (
	"context"
	"fmt"

	chroma_go "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

type IndexerConfig struct {
	Client    *chroma_go.Client
	Embedding embedding.Embedder
}

type Indexer struct {
	config     *IndexerConfig
	collection *chroma_go.Collection
}

func NewIndexer(ctx context.Context, config *IndexerConfig) (*Indexer, error) {
	if config.Client == nil {
		return nil, fmt.Errorf("[NewIndexer] chroma client not provided")
	}
	if config.Embedding == nil {
		return nil, fmt.Errorf("[NewIndexer] embedding not provided for chroma indexer")
	}

	collectionName := "eino-collection"
	// emb := &EinoEmbedderAdapter{embedder: config.Embedding}

	collection, err := config.Client.GetCollection(ctx, collectionName, nil)
	if err != nil {
		// If the collection doesn't exist, create it.
		// NOTE: The `embeddings` parameter is for the collection's embedding function.
		// Here, we can directly pass our embedder if it wraps a Chroma-compatible one.
		// Alternatively, if the Chroma server handles embedding, we'd pass nil or a name.
		collection, err = config.Client.CreateCollection(ctx, collectionName, nil, false, nil, "")
		if err != nil {
			return nil, fmt.Errorf("failed to get or create collection: %w", err)
		}
	}

	return &Indexer{
		config:     config,
		collection: collection,
	}, nil
}

func (i *Indexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	documents := make([]string, len(docs))
	metadatas := make([]map[string]interface{}, len(docs))
	ids := make([]string, len(docs))

	contentsToEmbed := make([]string, len(docs))

	for idx, doc := range docs {
		if doc.ID == "" {
			doc.ID = uuid.New().String()
		}
		documents[idx] = doc.Content
		metadatas[idx] = doc.MetaData
		ids[idx] = doc.ID
		contentsToEmbed[idx] = doc.Content
	}

	vectors64, err := i.config.Embedding.EmbedStrings(ctx, contentsToEmbed)
	if err != nil {
		return nil, fmt.Errorf("failed to embed documents: %w", err)
	}

	if len(vectors64) != len(docs) {
		return nil, fmt.Errorf("mismatch between number of documents and generated vectors")
	}

	chromaEmbeddings := make([]*types.Embedding, len(vectors64))
	for i, vec64 := range vectors64 {
		vec32 := make([]float32, len(vec64))
		for j, val := range vec64 {
			vec32[j] = float32(val)
		}

		chromaEmbeddings[i] = &types.Embedding{
			ArrayOfFloat32: &vec32,
		}
	}

	_, err = i.collection.Add(ctx, chromaEmbeddings, metadatas, documents, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to add documents to chroma: %w", err)
	}

	return ids, nil
}

// GetType returns the indexer's type.
func (i *Indexer) GetType() string {
	return "chroma"
}

// IsCallbacksEnabled checks if callbacks are enabled.
func (i *Indexer) IsCallbacksEnabled() bool {
	return true
}
