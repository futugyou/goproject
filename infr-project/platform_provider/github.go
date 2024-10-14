package platform_provider

import (
	"context"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(token string) (*GithubClient, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GithubClient{
		client,
	}, nil
}

func (g *GithubClient) CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *GithubClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	return nil, nil
}

func (g *GithubClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *GithubClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}
