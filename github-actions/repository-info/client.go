package main

import (
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Client *github.Client
}

func NewGithubClient() *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GithubClient{
		Client: client,
	}
}
