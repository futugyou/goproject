package chroma

import (
	"context"
	"fmt"

	chroma_go "github.com/amikos-tech/chroma-go/pkg/api/v2"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

type IndexerConfig struct {
	client    chroma_go.Client
	Embedding embedding.Embedder
}

type Indexer struct {
	config     *IndexerConfig
	collection chroma_go.Collection
}

func NewIndexer(ctx context.Context, config *IndexerConfig) (*Indexer, error) {
	if config.Embedding == nil {
		return nil, fmt.Errorf("[NewIndexer] embedding not provided for chroma indexer")
	}

	collectionName := "eino-collection"
	opts := []chroma_go.ClientOption{}
	client, err := chroma_go.NewHTTPClient(opts...)
	if err != nil {
		return nil, err
	}
	config.client = client

	emb := &EinoEmbedderAdapter{embedder: config.Embedding}
	getopts := []chroma_go.CreateCollectionOption{
		chroma_go.WithEmbeddingFunctionCreate(emb),
		chroma_go.WithIfNotExistsCreate(),
	}
	collection, err := config.client.GetOrCreateCollection(ctx, collectionName, getopts...)
	if err != nil {
		return nil, fmt.Errorf("failed to get or create collection: %w", err)
	}

	return &Indexer{
		config:     config,
		collection: collection,
	}, nil
}

func (i *Indexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	documents := make([]string, len(docs))
	metadatas := make([]chroma_go.DocumentMetadata, len(docs))
	ids := make([]chroma_go.DocumentID, len(docs))
	idss := make([]string, len(docs))

	contentsToEmbed := make([]string, len(docs))

	for idx, doc := range docs {
		if doc.ID == "" {
			doc.ID = uuid.New().String()
		}
		documents[idx] = doc.Content
		metadatas[idx] = newDocumentMetadata(doc.MetaData)
		ids[idx] = chroma_go.DocumentID(doc.ID)
		idss[idx] = doc.ID
		contentsToEmbed[idx] = doc.Content
	}

	vectors64, err := i.config.Embedding.EmbedStrings(ctx, contentsToEmbed)
	if err != nil {
		return nil, fmt.Errorf("failed to embed documents: %w", err)
	}

	if len(vectors64) != len(docs) {
		return nil, fmt.Errorf("mismatch between number of documents and generated vectors")
	}

	chromaEmbeddings, err := embeddings.NewEmbeddingsFromFloat32(float64ToFloat32(vectors64))
	if err != nil {
		return nil, err
	}
	addopts := []chroma_go.CollectionAddOption{
		chroma_go.WithEmbeddings(chromaEmbeddings...),
		chroma_go.WithTexts(documents...),
		chroma_go.WithIDs(ids...),
		chroma_go.WithMetadatas(metadatas...),
	}
	err = i.collection.Add(ctx, addopts...)
	if err != nil {
		return nil, fmt.Errorf("failed to add documents to chroma: %w", err)
	}

	return idss, nil
}

// GetType returns the indexer's type.
func (i *Indexer) GetType() string {
	return "chroma"
}

// IsCallbacksEnabled checks if callbacks are enabled.
func (i *Indexer) IsCallbacksEnabled() bool {
	return true
}

func newDocumentMetadata(datas map[string]any) chroma_go.DocumentMetadata {
	s := []*chroma_go.MetaAttribute{}
	for k, v := range datas {
		s = append(s, newAttributeFromAny(k, v))
	}
	return chroma_go.NewDocumentMetadata(s...)
}

func newAttributeFromAny(key string, value any) *chroma_go.MetaAttribute {
	switch v := value.(type) {
	case string:
		return chroma_go.NewStringAttribute(key, v)
	case int:
		return chroma_go.NewIntAttribute(key, int64(v))
	case int64:
		return chroma_go.NewIntAttribute(key, v)
	case float32:
		return chroma_go.NewFloatAttribute(key, float64(v))
	case float64:
		return chroma_go.NewFloatAttribute(key, v)
	case bool:
		return chroma_go.NewBoolAttribute(key, v)
	default:
		return chroma_go.NewStringAttribute(key, fmt.Sprintf("%v", value))
	}
}
