package services

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/futugyousuzu/go-openai-web/models"
)

func (e *EinoService) getGraphBranch(ctx context.Context, node models.Node) (*compose.GraphBranch, error) {
	return compose.NewGraphBranch(func(ctx context.Context, in map[string]any) (string, error) {
		return "", nil
	}, map[string]bool{}), nil
}

func (e *EinoService) getGraphNode(ctx context.Context, node models.Node) (compose.AnyGraph, error) {
	panic("unimplemented")
}
