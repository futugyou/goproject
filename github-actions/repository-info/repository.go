package main

import (
	"context"
	"log"

	"github.com/google/go-github/v53/github"
)

type RepositoryService struct {
	Client *github.Client
}

func NewRepositoryService(client *github.Client) *RepositoryService {
	return &RepositoryService{
		Client: client,
	}
}

type RepoModel struct {
	Owner  string `json:"owner"`
	Name   string `json:"name"`
	Branch string `json:"branch"`
}

func (r *RepositoryService) GetRepository(reponame string) []RepoModel {
	var result = make([]RepoModel, 0)
	var repos = make([]*github.Repository, 0)
	var err error
	if len(reponame) == 0 {
		opts := &github.RepositoryListOptions{
			Sort: "pushed",
			ListOptions: github.ListOptions{
				PerPage: 1000,
			},
		}

		repos, _, err = r.Client.Repositories.List(context.Background(), "", opts)
		if err != nil {
			log.Println(err)
		}
	} else {
		repo, _, err := r.Client.Repositories.Get(context.Background(), "", reponame)
		if err != nil {
			log.Println(err)
			return nil
		}

		repos = append(repos, repo)
	}

	for _, repo := range repos {
		result = append(result, RepoModel{
			Owner:  *repo.Owner.Login,
			Name:   *repo.Name,
			Branch: *repo.DefaultBranch,
		})
	}
	return result
}
