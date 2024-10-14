package platform_provider

import (
	"context"

	"github.com/futugyou/circleci"
)

type CircleClient struct {
	client *circleci.CircleciClient
}

func NewCircleClient(token string) (*CircleClient, error) {
	client := circleci.NewCircleciClient(token)
	return &CircleClient{
		client,
	}, nil
}

func (g *CircleClient) CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *CircleClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	return nil, nil
}

func (g *CircleClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *CircleClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}
