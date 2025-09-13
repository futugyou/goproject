package graph

import (
	"context"
	"fmt"
	"os"

	chroma_go "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/types"
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
			client, err := getChromaClient(ctx)
			if err != nil {
				return nil, err
			}
			return chroma.NewIndexer(ctx, &chroma.IndexerConfig{
				Embedding: embedding,
				Client:    client,
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

func getChromaClient(ctx context.Context) (*chroma_go.Client, error) {
	opt := []chroma_go.ClientOption{
		chroma_go.WithDatabase(os.Getenv("chroma_database")),
		chroma_go.WithBasePath(os.Getenv("chroma_base_path")),
		chroma_go.WithAuth(types.NewBasicAuthCredentialsProvider(os.Getenv("chroma_user"), os.Getenv("chroma_password"))),
	}

	return chroma_go.NewClient(opt...)
}

func getMongodbClient(ctx context.Context) (*mongo.Client, error) {
	uri := os.Getenv("mongodb_url")
	return mongo.Connect(ctx, options.Client().ApplyURI(uri))
}
