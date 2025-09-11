package services

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/recursive"
	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/semantic"
	"github.com/cloudwego/eino/components/document"
	"github.com/futugyousuzu/go-openai-web/models"
)

func (e *EinoService) getDocumentTransformerNode(ctx context.Context, node models.Node) (document.Transformer, error) {
	if transformer, ok := node.Data["transformer"].(string); ok && len(transformer) > 0 {
		switch transformer {
		case "markdown":
			headers := map[string]string{
				"#":   "h1",
				"##":  "h2",
				"###": "h3",
			}
			if h, ok := node.Data["transformer_header"].(map[string]string); ok {
				headers = h
			}
			return markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{Headers: headers, TrimHeaders: false})
		case "semantic":
			return semantic.NewSplitter(ctx, &semantic.Config{
				Embedding: e.embed,
			})
		case "recursive":
			return recursive.NewSplitter(ctx, &recursive.Config{ChunkSize: 1000, OverlapSize: 200})
		}

	}

	return nil, fmt.Errorf("invalid document transformer node: %s", node.ID)
}
