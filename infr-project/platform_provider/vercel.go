package platform_provider

import (
	"context"
	"fmt"

	"github.com/futugyou/vercel"
)

type VercelClient struct {
	client *vercel.VercelClient
}

func NewVercelClient(token string) (*VercelClient, error) {
	client := vercel.NewClient(token)
	return &VercelClient{
		client,
	}, nil
}

const VercelProjectUrl = "https://vercel.com/%s/%s"

// Optional. TEAM_SLUG TEAM_ID in CreateProjectRequest 'Parameters',
// default value ” ”.
func (g *VercelClient) CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error) {
	resultChan := make(chan *Project, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		req := vercel.CreateProjectRequest{
			Name: request.Name,
		}

		team_slug := ""
		team_id := ""
		if value, ok := request.Parameters["TEAM_SLUG"]; ok {
			team_slug = value
		}
		if value, ok := request.Parameters["TEAM_ID"]; ok {
			team_id = value
		}

		if len(team_slug) == 0 && len(team_id) == 0 {
			teamSlugCh, errCh := g.getDefaultTeamtSlug(ctx)
			select {
			case teamSlug := <-teamSlugCh:
				team_slug = *teamSlug
			case <-errCh:
			case <-ctx.Done():
			}
		}

		vercelProject, err := g.client.Project.CreateProject(team_slug, team_id, req)
		if err != nil {
			errorChan <- err
			return
		}

		url := ""
		if len(team_slug) > 0 {
			url = fmt.Sprintf(VercelProjectUrl, team_slug, vercelProject.Name)
		}

		resultChan <- &Project{
			ID:   vercelProject.Id,
			Name: vercelProject.Name,
			Url:  url,
		}
	}()

	return resultChan, errorChan
}

func (g *VercelClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	return nil, nil
}

func (g *VercelClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *VercelClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}

// TODO: pass context to all circelci/vercel sdk
func (g *VercelClient) getDefaultTeamtSlug(_ context.Context) (<-chan *string, <-chan error) {
	resultChan := make(chan *string, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		team_slug := ""
		teams, err := g.client.Team.ListTeam("", "", "")
		if err != nil {
			errorChan <- err
			return
		}

		if len(teams.Teams) > 0 {
			team_slug = teams.Teams[0].Slug
		}

		resultChan <- &team_slug
	}()

	return resultChan, errorChan
}
