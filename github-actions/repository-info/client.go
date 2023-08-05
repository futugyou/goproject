package main

import (
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

var (
	client *github.Client
)

func GetGitHubClient() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc)
}
