package graph

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/futugyousuzu/go-openai-web/eino_extensions/retriever/chroma"
	mongo_retriever "github.com/futugyousuzu/go-openai-web/eino_extensions/retriever/mongo"
	"github.com/futugyousuzu/go-openai-web/models"
)

func getRetrieverNode(ctx context.Context, node models.Node, embedding embedding.Embedder) (retriever.Retriever, error) {
	if idx, ok := node.Data["retriever"].(string); ok && len(idx) > 0 {
		switch idx {
		case "chroma":
			client, err := getChromaClient(ctx)
			if err != nil {
				return nil, err
			}
			return chroma.NewRetriever(ctx, &chroma.RetrieverConfig{
				Embedding:      embedding,
				Client:         client,
				CollectionName: os.Getenv("retriever_collection"),
			})
		case "mongo":
			client, err := getMongodbClient(ctx)
			if err != nil {
				return nil, err
			}
			return mongo_retriever.NewRetriever(ctx, &mongo_retriever.RetrieverConfig{
				Client:     client,
				Embedding:  embedding,
				Collection: os.Getenv("retriever_collection"),
				Database:   os.Getenv("retriever_database"),
			})
		}
	}

	return nil, fmt.Errorf("invalid retriever node: %s", node.ID)
}
