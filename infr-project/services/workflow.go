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

func (s *WorkflowService) Workflows(owner string, repo string) {
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

func (s *WorkflowService) WorkflowRuns(owner string, repo string, workflowID int64) {
	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100},
	}

	wfrs, _, err := s.client.Actions.ListWorkflowRunsByID(context.Background(), owner, repo, workflowID, opts)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, wf := range wfrs.WorkflowRuns {
		log.Println(wf.GetName())
	}
	log.Println(wfrs.GetTotalCount())
}

func (s *WorkflowService) WorkflowJobs(owner string, repo string, runID int64) {
	opts := &github.ListWorkflowJobsOptions{
		Filter: "latest",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100},
	}

	wfjs, _, err := s.client.Actions.ListWorkflowJobs(context.Background(), owner, repo, runID, opts)
	if err != nil {
		log.Println(err.Error())
		return
	}
	for _, wf := range wfjs.Jobs {
		log.Println(wf.GetName())
	}
	log.Println(wfjs.GetTotalCount())
}
