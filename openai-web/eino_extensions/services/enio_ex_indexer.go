package services

import (
	"context"

	"github.com/cloudwego/eino/components/indexer"
	"github.com/futugyousuzu/go-openai-web/models"
)

func (e *EinoService) getIndexerNode(ctx context.Context, node models.Node) (indexer.Indexer, error) {
	panic("unimplemented")
}
