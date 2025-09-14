package graph

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/futugyousuzu/go-openai-web/eino_extensions/indexer/chroma"
	mongo_indexer "github.com/futugyousuzu/go-openai-web/eino_extensions/indexer/mongo"
	"github.com/futugyousuzu/go-openai-web/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getIndexerNode(ctx context.Context, node models.Node, embedding embedding.Embedder) (indexer.Indexer, error) {
	if idx, ok := node.Data["indexer"].(string); ok && len(idx) > 0 {
		switch idx {
		case "chroma":
			return chroma.NewIndexer(ctx, &chroma.IndexerConfig{
				Embedding: embedding,
			})
		case "mongo":
			client, err := getMongodbClient(ctx)
			if err != nil {
				return nil, err
			}
			return mongo_indexer.NewIndexer(ctx, &mongo_indexer.IndexerConfig{
				Client:    client,
				Embedding: embedding,
			})
		}
	}

	return nil, fmt.Errorf("invalid indexer node: %s", node.ID)
}

func getMongodbClient(ctx context.Context) (*mongo.Client, error) {
	uri := os.Getenv("mongodb_url")
	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}
