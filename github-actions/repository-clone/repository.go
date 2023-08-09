package main

import (
	"context"
	"fmt"
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

func (r *RepositoryService) CreateRepository(info *CloneInfo) error {
	repoinfo := &github.Repository{
		Name:          github.String(info.DestName),
		DefaultBranch: github.String(info.DestBranch),
		MasterBranch:  github.String(info.DestBranch),
		Private:       github.Bool(true),
		// AutoInit:      github.Bool(true),
	}
	_, _, err := r.Client.Repositories.Create(context.Background(), info.DestOwner, repoinfo)
	if err != nil {
		fmt.Println(err)
		return err
	}

	git := NewGitCommand(info)
	git.SetConfig()
	git.InitRepository()
	actionsPermissionsRepository := github.ActionsPermissionsRepository{
		Enabled: github.Bool(false),
	}
	_, _, err = r.Client.Repositories.EditActionsPermissions(context.Background(), info.DestOwner, info.DestName, actionsPermissionsRepository)
	if err != nil {
		fmt.Println(err)
	}

	return err
}
