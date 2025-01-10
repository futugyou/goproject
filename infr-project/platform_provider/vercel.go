package platform_provider

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tool "github.com/futugyou/extensions"

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

		team_slug, team_id, _ := g.getTeamSlugAndId(ctx, request.Parameters)
		req := vercel.UpsertProjectRequest{
			Name: request.Name,
			BaseUrlParameter: vercel.BaseUrlParameter{
				TeamSlug: &team_slug,
				TeamId:   &team_id,
			},
		}

		vercelProject, err := g.client.Projects.CreateProject(ctx, req)
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
		request := vercel.ListProjectParameter{
			BaseUrlParameter: vercel.BaseUrlParameter{
				TeamSlug: &team_slug,
				TeamId:   &team_id,
			},
		}
		vercelProjects, err := g.client.Projects.ListProject(ctx, request)
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

			properties := map[string]string{}
			for key, v := range project.Targets {
				k := strings.ToUpper(fmt.Sprintf("%s_Alias", key))
				properties[k] = strings.Join(v.Alias, ",")
			}

			readState := ""
			if target, ok := project.Targets["production"]; ok {
				readState = target.ReadyState
			}

			badgeURL, badgeMarkdown := g.buildVercelBadge("Deployment", url, readState)
			projects = append(projects, Project{
				ID:            project.Id,
				Name:          project.Name,
				Url:           url,
				Properties:    properties,
				Envs:          g.buildVercelEnv(project.Env),
				Deployments:   g.buildVercelDeployment(project.LatestDeployments, url),
				BadgeURL:      badgeURL,
				BadgeMarkDown: badgeMarkdown,
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
		request := vercel.GetProjectParameter{
			IdOrName: filter.Name,
			BaseUrlParameter: vercel.BaseUrlParameter{
				TeamSlug: &team_slug,
				TeamId:   &team_id,
			},
		}
		vercelProject, err := g.client.Projects.GetProject(ctx, request)
		if err != nil {
			errorChan <- err
			return
		}

		hooks := []WebHook{}
		req := vercel.ListWebhookParameter{
			ProjectId: &vercelProject.Id,
			BaseUrlParameter: vercel.BaseUrlParameter{
				TeamSlug: &team_slug,
				TeamId:   &team_id,
			},
		}
		if vercelHooks, err := g.client.Webhooks.ListWebhook(ctx, req); err != nil {
			log.Println(err.Error())
		} else {
			for _, hook := range vercelHooks {
				paras := map[string]string{}
				paras["SigningSecret"] = hook.Secret
				hooks = append(hooks, WebHook{
					ID:         hook.Id,
					Name:       hook.Id,
					Url:        hook.Url,
					Events:     hook.Events,
					Activate:   true,
					Parameters: paras,
				})
			}
		}

		properties := map[string]string{}
		for key, v := range vercelProject.Targets {
			k := strings.ToUpper(fmt.Sprintf("%s_Alias", key))
			properties[k] = strings.Join(v.Alias, ",")
		}

		readState := ""
		if target, ok := vercelProject.Targets["production"]; ok {
			readState = target.ReadyState
		}

		url := ""
		if len(team_slug) > 0 {
			url = fmt.Sprintf(VercelProjectUrl, team_slug, vercelProject.Name)
		}

		badgeURL, badgeMarkdown := g.buildVercelBadge("Deployment", url, readState)
		project := &Project{
			ID:            vercelProject.Id,
			Name:          vercelProject.Name,
			Url:           url,
			WebHooks:      hooks,
			Properties:    properties,
			Envs:          g.buildVercelEnv(vercelProject.Env),
			Deployments:   g.buildVercelDeployment(vercelProject.LatestDeployments, url),
			BadgeURL:      badgeURL,
			BadgeMarkDown: badgeMarkdown,
		}

		resultChan <- project
	}()

	return resultChan, errorChan
}

func (*VercelClient) buildVercelBadge(lable string, url string, readyState string) (badgeUrl string, badgeMarkDown string) {
	lable = strings.ReplaceAll(lable, "-", "%20")
	badgeUrl = fmt.Sprintf(CommonProjectBadge, lable, "Norecord", "red", "vercel", url)
	if readyState == "READY" {
		badgeUrl = fmt.Sprintf(CommonProjectBadge, lable, readyState, "brightgreen", "vercel", url)
	} else if len(readyState) > 0 {
		badgeUrl = fmt.Sprintf(CommonProjectBadge, lable, readyState, "red", "vercel", url)
	}

	badgeMarkDown = fmt.Sprintf("![%s](%s)", lable, badgeUrl)
	return
}

func (*VercelClient) buildVercelEnv(vercelEnvs []vercel.Env) map[string]Env {
	envs := map[string]Env{}
	for _, v := range vercelEnvs {
		envs[v.ID] = Env{
			ID:        v.ID,
			Key:       v.Key,
			CreatedAt: tool.Int64ToTime(v.CreatedAt).Format(time.RFC3339Nano),
			UpdatedAt: tool.Int64ToTime(v.UpdatedAt).Format(time.RFC3339Nano),
			Type:      v.Type,
			Value:     v.Value,
		}
	}
	return envs
}

func (g *VercelClient) buildVercelDeployment(vercelDeployments []vercel.LatestDeployment, url string) map[string]Deployment {
	deployments := map[string]Deployment{}
	for _, v := range vercelDeployments {
		badgeURL, badgeMarkdown := g.buildVercelBadge(v.Name, url, v.ReadyState)
		deployments[v.ID] = Deployment{
			ID:            v.ID,
			Name:          v.Name,
			Plan:          v.Plan,
			ReadyState:    v.ReadyState,
			ReadySubstate: v.ReadySubstate,
			CreatedAt:     tool.Int64ToTime(v.CreatedAt).Format(time.RFC3339Nano),
			BadgeURL:      badgeURL,
			BadgeMarkdown: badgeMarkdown,
		}
	}
	return deployments
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
			BaseUrlParameter: vercel.BaseUrlParameter{
				TeamSlug: &team_slug,
				TeamId:   &team_id,
			},
		}
		vercelHook, err := g.client.Webhooks.CreateWebhook(ctx, req)
		if err != nil {
			errorChan <- err
			return
		}

		paras := map[string]string{}
		paras["SigningSecret"] = vercelHook.Secret
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
		request := vercel.ListTeamParameter{}
		teams, err := g.client.Teams.ListTeam(ctx, request)
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

func (g *VercelClient) DeleteWebHookAsync(ctx context.Context, request DeleteWebHookRequest) <-chan error {
	errorChan := make(chan error, 1)
	req := vercel.DeleteWebhookRequest{
		WebhookId: request.WebHookId,
	}

	go func() {
		defer close(errorChan)
		_, err := g.client.Webhooks.DeleteWebhook(ctx, req)
		errorChan <- err
	}()

	return errorChan
}

func (g *VercelClient) GetUserAsync(ctx context.Context) (<-chan *User, <-chan error) {
	resultChan := make(chan *User, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		vercelUser, err := g.client.User.GetUser(ctx)
		if err != nil {
			errorChan <- err
			return
		}

		user := &User{
			ID:   vercelUser.Id,
			Name: vercelUser.Name,
		}

		resultChan <- user
	}()

	return resultChan, errorChan
}
