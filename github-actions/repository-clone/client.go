package main

import (
	"context"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Client *github.Client
}

func NewGithubClient(token string) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(context.Background(), ts)

	client := github.NewClient(tc)
	return &GithubClient{
		Client: client,
	}
}
