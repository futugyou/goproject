package provider

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tool "github.com/futugyou/extensions"

	"github.com/futugyou/vercel"
)

type vercelClient struct {
	client *vercel.VercelClient
}

func newVercelClient(_ context.Context, token string) (*vercelClient, error) {
	client := vercel.NewClient(token)
	return &vercelClient{
		client,
	}, nil
}

const vercelProjectUrl = "https://vercel.com/%s/%s"

// Optional. TEAM_SLUG TEAM_ID in CreateProjectRequest 'Parameters',
// default value ” ”.
func (g *vercelClient) CreateProject(ctx context.Context, request CreateProjectRequest) (*Project, error) {
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
		return nil, err
	}

	url := ""
	if len(team_slug) > 0 {
		url = fmt.Sprintf(vercelProjectUrl, team_slug, vercelProject.Name)
	}

	return &Project{
		ID:   vercelProject.Id,
		Name: vercelProject.Name,
		Url:  url,
	}, nil
}

// Optional. TEAM_SLUG TEAM_ID in CreateProjectRequest 'Parameters',
// default value ” ”.
func (g *vercelClient) ListProject(ctx context.Context, filter ProjectFilter) ([]Project, error) {
	team_slug, team_id, _ := g.getTeamSlugAndId(ctx, filter.Parameters)
	request := vercel.ListProjectParameter{
		BaseUrlParameter: vercel.BaseUrlParameter{
			TeamSlug: &team_slug,
			TeamId:   &team_id,
		},
	}
	vercelProjects, err := g.client.Projects.ListProject(ctx, request)
	if err != nil {
		return nil, err
	}

	projects := []Project{}
	for _, project := range vercelProjects.Projects {
		url := ""
		if len(team_slug) > 0 {
			url = fmt.Sprintf(vercelProjectUrl, team_slug, project.Name)
		}

		properties := map[string]string{}
		environments := []string{}
		for key, v := range project.Targets {
			environments = append(environments, key)
			k := strings.ToUpper(fmt.Sprintf("%s_Alias", key))
			properties[k] = strings.Join(v.Alias, ",")
		}

		readState := ""
		if target, ok := project.Targets["production"]; ok {
			readState = target.ReadyState
		}

		badgeURL, badgeMarkdown := g.buildVercelBadge("Deployment", url, readState)
		projects = append(projects, Project{
			ID:                   project.Id,
			Name:                 project.Name,
			Url:                  url,
			Properties:           properties,
			EnvironmentVariables: g.buildVercelEnv(project.Env),
			Environments:         environments,
			Deployments:          g.buildVercelDeploymentWithLatest(project.LatestDeployments, url),
			BadgeURL:             badgeURL,
			BadgeMarkDown:        badgeMarkdown,
		})
	}

	return projects, nil
}

func (g *vercelClient) GetProject(ctx context.Context, filter ProjectFilter) (*Project, error) {
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
		return nil, err
	}

	webHook := g.getWebHook(ctx, filter, team_slug, team_id)

	properties := map[string]string{}
	environments := []string{}
	for key, v := range vercelProject.Targets {
		environments = append(environments, key)
		k := strings.ToUpper(fmt.Sprintf("%s_Alias", key))
		properties[k] = strings.Join(v.Alias, ",")
	}

	readState := ""
	if target, ok := vercelProject.Targets["production"]; ok {
		readState = target.ReadyState
	}

	url := ""
	if len(team_slug) > 0 {
		url = fmt.Sprintf(vercelProjectUrl, team_slug, vercelProject.Name)
	}

	deploymentRequestLimit := "20"
	deploymentRequest := vercel.ListDeploymentParameter{
		Limit:     &deploymentRequestLimit,
		ProjectId: &vercelProject.Id,
		BaseUrlParameter: vercel.BaseUrlParameter{
			TeamSlug: &team_slug,
			TeamId:   &team_id,
		},
	}

	deployments := map[string]Deployment{}
	if depls, err := g.client.Deployments.ListDeployment(ctx, deploymentRequest); err != nil {
		log.Println(err.Error())
	} else {
		deployments = g.buildVercelDeployment(depls.Deployments, url)
	}

	badgeURL, badgeMarkdown := g.buildVercelBadge("Deployment", url, readState)
	project := &Project{
		ID:                   vercelProject.Id,
		Name:                 vercelProject.Name,
		Url:                  url,
		WebHook:              webHook,
		Properties:           properties,
		EnvironmentVariables: g.buildVercelEnv(vercelProject.Env),
		Deployments:          deployments, //g.buildVercelDeployment(vercelProject.LatestDeployments, url),
		BadgeURL:             badgeURL,
		Environments:         environments,
		BadgeMarkDown:        badgeMarkdown,
	}

	return project, nil
}

func (g *vercelClient) GetSimpleProjectInfo(ctx context.Context, filter ProjectFilter) (*Project, error) {
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
		return nil, err
	}

	properties := map[string]string{}
	environments := []string{}
	for key, v := range vercelProject.Targets {
		environments = append(environments, key)
		k := strings.ToUpper(fmt.Sprintf("%s_Alias", key))
		properties[k] = strings.Join(v.Alias, ",")
	}

	readState := ""
	if target, ok := vercelProject.Targets["production"]; ok {
		readState = target.ReadyState
	}

	url := ""
	if len(team_slug) > 0 {
		url = fmt.Sprintf(vercelProjectUrl, team_slug, vercelProject.Name)
	}

	badgeURL, badgeMarkdown := g.buildVercelBadge("Deployment", url, readState)

	project := &Project{
		ID:                   vercelProject.Id,
		Name:                 vercelProject.Name,
		Url:                  url,
		Properties:           properties,
		EnvironmentVariables: g.buildVercelEnv(vercelProject.Env),
		BadgeURL:             badgeURL,
		Environments:         environments,
		BadgeMarkDown:        badgeMarkdown,
	}

	return project, nil
}

func (g *vercelClient) getWebHook(ctx context.Context, filter ProjectFilter, team_slug string, team_id string) *WebHook {
	req := vercel.GetWebhookParameter{
		WebhookId: *filter.WebHookID,
		BaseUrlParameter: vercel.BaseUrlParameter{
			TeamSlug: &team_slug,
			TeamId:   &team_id,
		},
	}

	hook, err := g.client.Webhooks.GetWebhook(ctx, req)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	paras := map[string]string{}
	paras["SigningSecret"] = hook.Secret
	return &WebHook{
		ID:         hook.Id,
		Name:       hook.Id,
		Url:        hook.Url,
		Events:     hook.Events,
		Activate:   true,
		Parameters: paras,
	}
}

func (*vercelClient) buildVercelBadge(lable string, url string, readyState string) (badgeUrl string, badgeMarkDown string) {
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

func (*vercelClient) buildVercelEnv(vercelEnvs []vercel.Env) map[string]EnvironmentVariable {
	envs := map[string]EnvironmentVariable{}
	for _, v := range vercelEnvs {
		envs[v.ID] = EnvironmentVariable{
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

func (g *vercelClient) buildVercelDeployment(vercelDeployments []vercel.Deployment, url string) map[string]Deployment {
	deployments := map[string]Deployment{}
	for _, v := range vercelDeployments {
		if v.Target == nil || len(v.Name) == 0 || len(v.ReadyState) == 0 {
			continue
		}

		badgeURL, badgeMarkdown := g.buildVercelBadge(v.Name, url, v.ReadyState)
		deployments[v.Uid] = Deployment{
			ID:            v.Uid,
			Name:          v.Name,
			Environment:   *v.Target,
			ReadyState:    v.ReadyState,
			ReadySubstate: v.ReadySubstate,
			CreatedAt:     tool.Int64ToTime(v.CreatedAt).Format(time.RFC3339Nano),
			BadgeURL:      badgeURL,
			BadgeMarkdown: badgeMarkdown,
			Description:   v.Meta.GitCommitMessage,
		}
	}
	return deployments
}

func (g *vercelClient) buildVercelDeploymentWithLatest(vercelDeployments []vercel.LatestDeployment, url string) map[string]Deployment {
	deployments := map[string]Deployment{}
	for _, v := range vercelDeployments {
		if len(v.Name) == 0 || len(v.Target) == 0 || len(v.ReadyState) == 0 {
			continue
		}

		badgeURL, badgeMarkdown := g.buildVercelBadge(v.Name, url, v.ReadyState)
		deployments[v.ID] = Deployment{
			ID:            v.ID,
			Name:          v.Name,
			Environment:   v.Target,
			ReadyState:    v.ReadyState,
			ReadySubstate: v.ReadySubstate,
			CreatedAt:     tool.Int64ToTime(v.CreatedAt).Format(time.RFC3339Nano),
			BadgeURL:      badgeURL,
			BadgeMarkdown: badgeMarkdown,
		}
	}
	return deployments
}

func (g *vercelClient) CreateWebHook(ctx context.Context, request CreateWebHookRequest) (*WebHook, error) {
	team_slug, team_id, _ := g.getTeamSlugAndId(ctx, request.Parameters)
	events := request.Events
	if len(events) == 0 {
		events = []string{"deployment.succeeded"}
	} else {
		events = Intersect(events, vercel.WebHookEvent)
	}

	req := vercel.CreateWebhookRequest{
		Events:     events,
		Url:        request.Url,
		ProjectIds: []string{request.ProjectID},
		BaseUrlParameter: vercel.BaseUrlParameter{
			TeamSlug: &team_slug,
			TeamId:   &team_id,
		},
	}
	vercelHook, err := g.client.Webhooks.CreateWebhook(ctx, req)
	if err != nil {
		return nil, err
	}

	paras := map[string]string{}
	paras["SigningSecret"] = vercelHook.Secret
	hook := &WebHook{
		ID:         vercelHook.Id,
		Name:       request.Name,
		Url:        vercelHook.Url,
		Events:     vercelHook.Events,
		Parameters: paras,
	}

	return hook, nil
}

func (g *vercelClient) GetWebHookByUrl(ctx context.Context, request GetWebHookRequest) (*WebHook, error) {
	team_slug, team_id, _ := g.getTeamSlugAndId(ctx, request.Parameters)

	req := vercel.ListWebhookParameter{
		ProjectId: &request.ProjectID,
		BaseUrlParameter: vercel.BaseUrlParameter{
			TeamSlug: &team_slug,
			TeamId:   &team_id,
		},
	}

	hooks, err := g.client.Webhooks.ListWebhook(ctx, req)
	if err != nil {
		return nil, err
	}

	for _, vercelHook := range hooks {
		if vercelHook.Url == request.Url {
			paras := map[string]string{}
			paras["SigningSecret"] = vercelHook.Secret
			return &WebHook{
				ID:         vercelHook.Id,
				Name:       "",
				Url:        vercelHook.Url,
				Events:     vercelHook.Events,
				Parameters: paras,
			}, nil
		}
	}

	return nil, fmt.Errorf("webhook not found")
}

func (g *vercelClient) getTeamSlugAndId(ctx context.Context, parameters map[string]string) (team_slug string, team_id string, err error) {
	if value, ok := parameters["TEAM_SLUG"]; ok {
		team_slug = value
	}
	if value, ok := parameters["TEAM_ID"]; ok {
		team_id = value
	}

	if len(team_slug) == 0 && len(team_id) == 0 {
		team, er := g.getDefaultTeam(ctx)
		if er != nil {
			err = er
			return
		}

		team_slug = team.Slug
		team_id = team.Id
	}

	return
}

func (g *vercelClient) getDefaultTeam(ctx context.Context) (*VercelTeam, error) {
	team := new(VercelTeam)
	request := vercel.ListTeamParameter{}
	teams, err := g.client.Teams.ListTeam(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(teams.Teams) > 0 {
		team.Slug = teams.Teams[0].Slug
		team.Id = teams.Teams[0].Id
		team.Name = teams.Teams[0].Name
	}

	return team, nil
}

type VercelTeam struct {
	Id   string
	Slug string
	Name string
}

func (g *vercelClient) DeleteWebHook(ctx context.Context, request DeleteWebHookRequest) error {
	req := vercel.DeleteWebhookRequest{
		WebhookId: request.WebHookID,
	}

	_, err := g.client.Webhooks.DeleteWebhook(ctx, req)

	return err
}

func (g *vercelClient) GetUser(ctx context.Context) (*User, error) {
	vercelUser, err := g.client.User.GetUser(ctx)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:   vercelUser.Id,
		Name: vercelUser.Name,
	}

	return user, nil
}
