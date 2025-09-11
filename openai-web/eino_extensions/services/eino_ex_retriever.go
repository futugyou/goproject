package services

import (
	"context"

	"github.com/cloudwego/eino/components/retriever"
	"github.com/futugyousuzu/go-openai-web/models"
)

func (e *EinoService) getRetrieverNode(ctx context.Context, node models.Node) (retriever.Retriever, error) {
	panic("unimplemented")
}
