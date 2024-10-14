package platform_provider

import (
	"context"

	"github.com/futugyou/vercel"
)

type VercelClient struct {
	client *vercel.VercelClient
}

func NewVercelClient(token string) (*VercelClient, error) {
	client := vercel.NewVercelClient(token)
	return &VercelClient{
		client,
	}, nil
}

func (g *VercelClient) CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *VercelClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	return nil, nil
}

func (g *VercelClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *VercelClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}
