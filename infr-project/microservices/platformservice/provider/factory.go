package provider

import (
	"context"
	"fmt"
)

func PlatformProviderFactory(ctx context.Context, provider string, token string) (PlatformProvider, error) {
	if provider == "circleci" {
		client, err := newCircleClient(ctx, token)
		return client, err
	}
	if provider == "github" {
		client, err := newGithubClient(ctx, token)
		return client, err
	}
	if provider == "vercel" {
		client, err := newVercelClient(ctx, token)
		return client, err
	}
	return nil, fmt.Errorf("provider type %s is not support", provider)
}
