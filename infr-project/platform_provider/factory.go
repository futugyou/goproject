package platform_provider

import "fmt"

func PlatformProviderFactory(provider string, token string) (IPlatformProviderAsync, error) {
	if provider == "circleci" {
		client, err := NewCircleClient(token)
		return client, err
	}
	if provider == "github" {
		client, err := NewGithubClient(token)
		return client, err
	}
	if provider == "vercel" {
		client, err := NewVercelClient(token)
		return client, err
	}
	return nil, fmt.Errorf("provider type %s is not support", provider)
}
