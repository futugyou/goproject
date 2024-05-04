package services

import (
	"context"
	"log"

	"github.com/google/go-github/v61/github"
)

type WorkflowService struct {
	client *github.Client
}

func NewWorkflowService(token string) *WorkflowService {
	client := github.NewClient(nil).WithAuthToken(token)
	return &WorkflowService{
		client: client,
	}
}

func (s *WorkflowService) Workflow(owner string, repo string) {
	opts := &github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	wfs, _, err := s.client.Actions.ListWorkflows(context.Background(), owner, repo, opts)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, wf := range wfs.Workflows {
		log.Println(wf.GetName())
	}
	log.Println(wfs.GetTotalCount())
}
