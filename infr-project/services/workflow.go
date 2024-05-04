package services

import "github.com/google/go-github/v61/github"

type WorkflowService struct {
	client *github.Client
}

func NewWorkflowService(token string) *WorkflowService {
	client := github.NewClient(nil).WithAuthToken(token)
	return &WorkflowService{
		client: client,
	}
}
