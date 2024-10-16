package platform_provider

import (
	"context"
	"fmt"

	"github.com/futugyou/circleci"
)

type CircleClient struct {
	client *circleci.CircleciClient
}

func NewCircleClient(token string) (*CircleClient, error) {
	client := circleci.NewCircleciClient(token)
	return &CircleClient{
		client,
	}, nil
}

func (g *CircleClient) CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error) {
	resultChan := make(chan *Project, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		org_slug := ""
		if org, ok := request.Parameters["org_slug"]; ok {
			org_slug = org
		} else {
			errorChan <- fmt.Errorf("create project request need 'org_slug' in parameters")
			return
		}

		if _, err := g.client.Project.CreateProject(org_slug, request.Name); err != nil {
			errorChan <- err
			return
		}
		project, err := g.client.Project.GetProject(org_slug, request.Name)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- &Project{
			ID:   project.ID,
			Name: project.Name,
			Url:  "", // url may set in applcation layer
		}
	}()

	return resultChan, errorChan
}

func (g *CircleClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	return nil, nil
}

func (g *CircleClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *CircleClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}
