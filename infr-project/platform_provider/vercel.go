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
			teamCh, errCh := g.getDefaultTeam(ctx)
			select {
			case team := <-teamCh:
				if team != nil {
					team_slug = team.Slug
					team_id = team.Id
				}
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
func (g *VercelClient) getDefaultTeam(_ context.Context) (<-chan *VercelTeam, <-chan error) {
	resultChan := make(chan *VercelTeam, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		team := new(VercelTeam)
		teams, err := g.client.Team.ListTeam("", "", "")
		if err != nil {
			errorChan <- err
			return
		}

		if len(teams.Teams) > 0 {
			team.Slug = teams.Teams[0].Slug
			team.Id = teams.Teams[0].Id
			team.Name = teams.Teams[0].Name
		}

		resultChan <- team
	}()

	return resultChan, errorChan
}

type VercelTeam struct {
	Id   string
	Slug string
	Name string
}