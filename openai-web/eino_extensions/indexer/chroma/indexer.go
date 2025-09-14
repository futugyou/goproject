package chroma

import (
	"context"
	"fmt"
	"os"

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

	collection, err := config.client.GetOrCreateCollection(
		ctx,
		collectionName,
		chroma_go.WithEmbeddingFunctionCreate(emb),
		chroma_go.WithIfNotExistsCreate(),
	)
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
	err = i.collection.Add(
		ctx,
		chroma_go.WithEmbeddings(chromaEmbeddings...),
		chroma_go.WithTexts(documents...),
		chroma_go.WithIDs(ids...),
		chroma_go.WithMetadatas(metadatas...),
	)
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
