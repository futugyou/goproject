package mongo

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexerConfig struct {
	Client           *mongo.Client
	Database         string
	Collection       string
	IndexName        *string
	VectorDim        int
	BatchSize        int `json:"batch_size"`
	Embedding        embedding.Embedder
	DocumentToFields func(ctx context.Context, doc *schema.Document) (*Fields, error)
}

type Fields struct {
	ID       string                 `bson:"_id" json:"_id" `
	Content  string                 `bson:"content" json:"content" `
	Vector   []float64              `bson:"plot_embedding" json:"plot_embedding" `
	Metadata map[string]interface{} `bson:"metadata" json:"metadata" `
}

type Indexer struct {
	config *IndexerConfig
}

func NewIndexer(ctx context.Context, config *IndexerConfig) (*Indexer, error) {
	if config.Embedding == nil {
		return nil, fmt.Errorf("[NewIndexer] embedding not provided for mongo indexer")
	}
	if config.Client == nil {
		return nil, fmt.Errorf("[NewIndexer] mongo client not provided")
	}

	if config.Database == "" {
		config.Database = "enio_database"
	}
	if config.Collection == "" {
		config.Collection = "enio_documents"
	}
	if config.DocumentToFields == nil {
		config.DocumentToFields = defaultDocumentToFields
	}
	if config.BatchSize == 0 {
		config.BatchSize = 10
	}
	return &Indexer{config: config}, nil
}

func (i *Indexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) (ids []string, err error) {
	options := indexer.GetCommonOptions(&indexer.Options{
		Embedding: i.config.Embedding,
	}, opts...)

	ctx = callbacks.EnsureRunInfo(ctx, i.GetType(), components.ComponentOfIndexer)
	ctx = callbacks.OnStart(ctx, &indexer.CallbackInput{Docs: docs})
	defer func() {
		if err != nil {
			callbacks.OnError(ctx, err)
		}
	}()

	if err = i.batchUpsert(ctx, docs, options); err != nil {
		return nil, err
	}

	ids = make([]string, 0, len(docs))
	for _, doc := range docs {
		ids = append(ids, doc.ID)
	}
	callbacks.OnEnd(ctx, &indexer.CallbackOutput{IDs: ids})
	return ids, nil
}

func (i *Indexer) batchUpsert(ctx context.Context, docs []*schema.Document, options *indexer.Options) error {
	emb := options.Embedding
	batchSize := i.config.BatchSize

	coll := i.config.Client.Database(i.config.Database).Collection(i.config.Collection)

	if err := i.ensureIndex(ctx, coll); err != nil {
		if isIndexExistsError(err) {
			return err
		}
	}

	for start := 0; start < len(docs); start += batchSize {
		end := start + batchSize
		if end > len(docs) {
			end = len(docs)
		}
		batch := docs[start:end]
		var (
			docs  []Fields
			texts []string
		)
		for _, doc := range batch {
			fields, err := i.config.DocumentToFields(ctx, doc)
			if err != nil {
				return err
			}
			docs = append(docs, *fields)
			texts = append(texts, fields.Content)
		}
		vectors, err := emb.EmbedStrings(ctx, texts)
		if err != nil {
			return fmt.Errorf("[batchUpsert] embedding failed, %w", err)
		}
		if len(vectors) != len(batch) {
			return fmt.Errorf("[batchUpsert] invalid vector length, expected=%d, got=%d", len(batch), len(vectors))
		}
		for i := 0; i < len(docs); i++ {
			docs[i].Vector = vectors[i]
		}
		var writeModels []mongo.WriteModel
		for _, doc := range docs {
			filter := bson.M{"_id": doc.ID}
			update := bson.M{"$set": doc}
			model := mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(update).
				SetUpsert(true)

			writeModels = append(writeModels, model)
		}
		_, err = coll.BulkWrite(ctx, writeModels)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Indexer) ensureIndex(ctx context.Context, coll *mongo.Collection) error {
	const indexName = "vector_search_index"
	opts := options.SearchIndexes().SetName(indexName).SetType("vectorSearch")
	// Defines the index definition
	vectorSearchIndexModel := mongo.SearchIndexModel{
		Definition: vectorDefinition{
			Fields: []vectorDefinitionField{{
				Type:          "vector",
				Path:          "plot_embedding",
				NumDimensions: 1536,
				Similarity:    "dotProduct",
				Quantization:  "scalar",
			}},
		},
		Options: opts,
	}
	// Creates the index
	_, err := coll.SearchIndexes().CreateOne(ctx, vectorSearchIndexModel)
	return err
}

func (i *Indexer) GetType() string {
	return "mongo"
}

func (i *Indexer) IsCallbacksEnabled() bool {
	return true
}

func defaultDocumentToFields(ctx context.Context, doc *schema.Document) (*Fields, error) {
	if doc.ID == "" {
		return nil, fmt.Errorf("[defaultDocumentToFields] doc id not set")
	}
	return &Fields{
		ID:       doc.ID,
		Content:  doc.Content,
		Metadata: doc.MetaData,
	}, nil
}

func float64SliceToFloat32(v []float64) []float32 {
	f := make([]float32, len(v))
	for i, x := range v {
		f[i] = float32(x)
	}
	return f
}

func isIndexExistsError(err error) bool {
	if cmdErr, ok := err.(mongo.CommandError); ok {
		if cmdErr.Code == 48 || cmdErr.Code == 85 {
			return true
		}
	}
	return strings.Contains(err.Error(), "Index with name") || strings.Contains(err.Error(), "already exists")
}

type vectorDefinitionField struct {
	Type          string `bson:"type"`
	Path          string `bson:"path"`
	NumDimensions int    `bson:"numDimensions"`
	Similarity    string `bson:"similarity"`
	Quantization  string `bson:"quantization"`
}
type vectorDefinition struct {
	Fields []vectorDefinitionField `bson:"fields"`
}
