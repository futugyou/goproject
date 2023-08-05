package main

import (
	"log"

	"github.com/google/go-github/v53/github"
)

var repos = make([]*github.Repository, 0)

func GetRepository() {
	var err error
	if reponame == nil {
		opts := &github.RepositoryListOptions{
			Sort: "pushed",
			ListOptions: github.ListOptions{
				PerPage: 1000,
			},
		}

		repos, _, err = client.Repositories.List(ctx, "", opts)
		if err != nil {
			log.Println(err)
		}
	} else {
		repo, _, err := client.Repositories.Get(ctx, "", *reponame)
		if err != nil {
			log.Println(err)
			return
		}

		repos = append(repos, repo)
	}
}
