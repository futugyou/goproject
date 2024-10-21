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

		team_slug, team_id, _ := g.getTeamSlugAndId(ctx, request.Parameters)

		vercelProject, err := g.client.Projects.CreateProject(ctx, team_slug, team_id, req)
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

// Optional. TEAM_SLUG TEAM_ID in CreateProjectRequest 'Parameters',
// default value ” ”.
func (g *VercelClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	resultChan := make(chan []Project, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		team_slug, team_id, _ := g.getTeamSlugAndId(ctx, filter.Parameters)

		vercelProjects, err := g.client.Projects.ListProject(ctx, team_slug, team_id)
		if err != nil {
			errorChan <- err
			return
		}

		projects := []Project{}
		for _, project := range vercelProjects.Projects {
			url := ""
			if len(team_slug) > 0 {
				url = fmt.Sprintf(VercelProjectUrl, team_slug, project.Name)
			}

			projects = append(projects, Project{
				ID:   project.Id,
				Name: project.Name,
				Url:  url,
			})
		}

		resultChan <- projects
	}()

	return resultChan, errorChan
}

func (g *VercelClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	resultChan := make(chan *Project, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		team_slug, team_id, _ := g.getTeamSlugAndId(ctx, filter.Parameters)

		vercelProject, err := g.client.Projects.GetProject(ctx, filter.Name, team_slug, team_id)
		if err != nil {
			errorChan <- err
			return
		}

		url := ""
		if len(team_slug) > 0 {
			url = fmt.Sprintf(VercelProjectUrl, team_slug, vercelProject.Name)
		}

		vercelHooks, err := g.client.Webhooks.ListWebhook(ctx, vercelProject.Id, team_slug, team_id)
		if err != nil {
			errorChan <- err
			return
		}

		hooks := []WebHook{}
		for _, hook := range vercelHooks {
			paras := map[string]string{}
			paras["Secret"] = hook.Secret
			hooks = append(hooks, WebHook{
				ID:         hook.Id,
				Name:       hook.Id,
				Url:        hook.Url,
				Events:     hook.Events,
				Parameters: paras,
			})
		}
		project := &Project{
			ID:    vercelProject.Id,
			Name:  vercelProject.Name,
			Url:   url,
			Hooks: hooks,
		}

		resultChan <- project
	}()

	return resultChan, errorChan
}

func (g *VercelClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	resultChan := make(chan *WebHook, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		team_slug, team_id := "", ""
		events := request.WebHook.Events
		if len(events) == 0 {
			events = []string{"deployment.succeeded"}
		} else {
			events = Intersect(events, vercel.WebHookEvent)
		}

		req := vercel.CreateWebhookRequest{
			Events:     events,
			Url:        request.WebHook.Url,
			ProjectIds: []string{request.ProjectId},
		}
		vercelHook, err := g.client.Webhooks.CreateWebhook(ctx, team_slug, team_id, req)
		if err != nil {
			errorChan <- err
			return
		}

		paras := map[string]string{}
		paras["Secret"] = vercelHook.Secret
		hook := &WebHook{
			ID:         vercelHook.Id,
			Name:       vercelHook.Id,
			Url:        vercelHook.Url,
			Events:     vercelHook.Events,
			Parameters: paras,
		}

		resultChan <- hook
	}()

	return resultChan, errorChan
}

func (g *VercelClient) getTeamSlugAndId(ctx context.Context, parameters map[string]string) (team_slug string, team_id string, err error) {
	if value, ok := parameters["TEAM_SLUG"]; ok {
		team_slug = value
	}
	if value, ok := parameters["TEAM_ID"]; ok {
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
		case err = <-errCh:
		case <-ctx.Done():
			err = fmt.Errorf("context timeout")
		}
	}
	return
}

func (g *VercelClient) getDefaultTeam(ctx context.Context) (<-chan *VercelTeam, <-chan error) {
	resultChan := make(chan *VercelTeam, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		team := new(VercelTeam)
		teams, err := g.client.Teams.ListTeam(ctx, "", "", "")
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
