package platform_provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/futugyou/circleci"
)

type CircleClient struct {
	client *circleci.CircleciClient
	token  string
}

func NewCircleClient(token string) (*CircleClient, error) {
	client := circleci.NewClientV2(token)
	return &CircleClient{
		client,
		token,
	}, nil
}

func NewCircleClientV1(token string) (*CircleClient, error) {
	client := circleci.NewClientV1(token)
	return &CircleClient{
		client,
		token,
	}, nil
}

const CircleciProjectUrl = "https://app.circleci.com/pipelines/%s/%s/%s"
const CircleciProjectBadge = "https://dl.circleci.com/status-badge/img/%s/%s/%s/tree/%s.svg?style=svg"
const CircleciProjectBadgeFull = "[![CircleCI](https://dl.circleci.com/status-badge/img/%s/%s/%s/tree/%s.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/%s/%s/%s/tree/%s)"
const CircleciWorkflowBadge = "https://dl.circleci.com/insights-snapshot/%s/%s/%s/%s/%s/badge.svg?window=30d"
const CircleciWorkflowBadgeFull = "[![CircleCI](https://dl.circleci.com/insights-snapshot/%s/%s/%s/%s/%s/badge.svg?window=30d)](https://app.circleci.com/insights/%s/%s/%s/workflows/%s/overview?branch=%s&reporting-window=last-30-days&insights-snapshot=true)"

var CircleciVCSMapping = map[string]string{"gh": "github", "bb": "bitbucket"}
var CircleciBadgeVCSMapping = map[string]string{"github": "gh", "bitbucket": "bb"}

// Parameters MUST include org_slug. eg. gh/demo
// Parameters can set circleci_project_url. eg. https://app.circleci.com/pipelines/%s/%s/%s
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

		if _, err := g.client.Project.CreateProject(ctx, org_slug, request.Name); err != nil {
			errorChan <- err
			return
		}

		project, err := g.client.Project.GetProject(ctx, org_slug, request.Name)
		if err != nil {
			errorChan <- err
			return
		}

		orgSplit := strings.Split(org_slug, "/")
		url := ""
		if len(orgSplit) == 2 {
			_, vcs_full := g.getCircleciVCS(orgSplit[0])
			circleciProjectUrl := CircleciProjectUrl
			if url, ok := request.Parameters["circleci_project_url"]; ok {
				circleciProjectUrl = url
			}
			url = fmt.Sprintf(circleciProjectUrl, vcs_full, orgSplit[1], project.Name)
		}

		resultChan <- &Project{
			ID:   project.ID,
			Name: project.Name,
			Url:  url,
		}
	}()

	return resultChan, errorChan
}

// For now circleci api v2 do not include list project api,
// So, we need use api v1.1
func (g *CircleClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	resultChan := make(chan []Project, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		client, err := NewCircleClientV1(g.token)
		if err != nil {
			errorChan <- err
			return
		}

		circleciProjects, err := client.client.Project.ListProject(ctx)
		if err != nil {
			errorChan <- err
			return
		}

		url := CircleciProjectUrl
		if u, ok := filter.Parameters["circleci_project_url"]; ok {
			url = u
		}

		projects := []Project{}
		for _, pro := range circleciProjects {
			projects = append(projects, g.buildProject(pro, url))
		}

		resultChan <- projects
	}()

	return resultChan, errorChan
}

// Need org_slug in ProjectFilter 'Parameters'
// circleci webhook will set some other information in hook Parameters, it include 'Scope' 'SigningSecret' 'VerifyTLS'
func (g *CircleClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	resultChan := make(chan *Project, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		org_slug := ""
		if org, ok := filter.Parameters["org_slug"]; ok {
			org_slug = org
		} else {
			errorChan <- fmt.Errorf("get project request need 'org_slug' in parameters")
			return
		}

		circleciProject, err := g.client.Project.GetProject(ctx, org_slug, filter.Name)
		if err != nil {
			errorChan <- err
			return
		}

		vcs, vcs_full := g.getCircleciVCS(circleciProject.VcsInfo.Provider)

		url_format := CircleciProjectUrl
		if url, ok := filter.Parameters["circleci_project_url"]; ok {
			url_format = url
		}

		url := g.buildProjectUrl(org_slug, url_format, vcs_full, circleciProject.Name)

		webHooks := []WebHook{}
		if circleciWebhooks, err := g.client.Webhook.ListWebhook(ctx, circleciProject.ID); err != nil {
			log.Println(err.Error())
		} else {
			for _, hook := range circleciWebhooks.Items {
				paras := map[string]string{}
				paras["Scope"] = hook.Scope.Type
				paras["ScopeId"] = hook.Scope.Id
				paras["SigningSecret"] = hook.SigningSecret
				paras["VerifyTLS"] = strconv.FormatBool(hook.VerifyTLS)
				webHooks = append(webHooks, WebHook{
					ID:         hook.Id,
					Name:       hook.Name,
					Url:        hook.Url,
					Events:     hook.Events,
					Activate:   true,
					Parameters: paras,
				})
			}
		}

		envs := map[string]EnvironmentVariable{}
		if circleciEnvs, err := g.client.Project.GetEnvironmentVariables(ctx, circleciProject.Slug); err != nil {
			log.Println(err.Error())
		} else {
			for _, e := range circleciEnvs.Items {
				envs[e.Name] = EnvironmentVariable{
					ID:        e.Name,
					Key:       e.Name,
					CreatedAt: e.CreatedAt,
					UpdatedAt: e.CreatedAt,
					Type:      "Masked",
					Value:     e.Value,
				}
			}
		}

		runs := map[string]WorkflowRun{}
		if circleciRuns, err := g.client.Pipeline.GetYourPipelines(ctx, circleciProject.Slug); err != nil {
			log.Println(err.Error())
		} else {
			for _, e := range circleciRuns.Items {
				name := fmt.Sprintf("%d", e.Number)
				if len(e.Vcs.Commit.Subject) > 0 {
					name = e.Vcs.Commit.Subject
				}

				description := ""
				if len(e.Vcs.Commit.Body) > 0 {
					description = e.Vcs.Commit.Body
				}

				badgeURL, badgeMarkdown := buildCommonBadge(name, e.State, "created", "circleci", nil)
				runs[e.ID] = WorkflowRun{
					ID:            e.ID,
					Name:          name,
					Description:   description,
					Status:        e.State,
					CreatedAt:     e.CreatedAt,
					BadgeURL:      badgeURL,
					BadgeMarkdown: badgeMarkdown,
				}
			}
		}

		badgeURL, badgeMarkdown := "", ""
		vcsurlinfos := strings.Split(circleciProject.VcsInfo.VcsURL, "/")
		if len(vcsurlinfos) >= 2 {
			badgeURL, badgeMarkdown = g.buildProjectBadge(vcs, vcsurlinfos[len(vcsurlinfos)-2], vcsurlinfos[len(vcsurlinfos)-1], circleciProject.VcsInfo.DefaultBranch)
		}

		project := &Project{
			ID:                   circleciProject.Name,
			Name:                 circleciProject.Name,
			Url:                  url,
			Description:          circleciProject.GetMessage(),
			WebHooks:             webHooks,
			Properties:           map[string]string{"VCS_TYPE": vcs_full, "VCS_URL": circleciProject.VcsInfo.VcsURL},
			EnvironmentVariables: envs,
			Environments:         []string{},
			WorkflowRuns:         runs,
			BadgeURL:             badgeURL,
			BadgeMarkDown:        badgeMarkdown,
		}

		resultChan <- project
	}()

	return resultChan, errorChan
}

// if need webhook secret, set it in WebHook.Parameters with key 'SigningSecret'
func (g *CircleClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	resultChan := make(chan *WebHook, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)
		secret := ""
		if s, ok := request.WebHook.Parameters["SigningSecret"]; ok && len(s) > 0 {
			secret = s
		}
		verifyTLS := false
		if tls, ok := request.WebHook.Parameters["VerifyTLS"]; ok {
			if value, err := strconv.ParseBool(tls); err != nil {
				verifyTLS = value
			}
		}

		events := request.WebHook.Events
		if len(events) == 0 {
			events = circleci.WebhookEvents
		} else {
			events = Intersect(events, circleci.WebhookEvents)
		}

		req := circleci.CreateWebhookRequest{
			Name:          request.WebHook.Name,
			Events:        events,
			Url:           request.WebHook.Url,
			VerifyTLS:     verifyTLS,
			SigningSecret: secret,
			ScopeId:       request.ProjectId,
			ScopeType:     "project",
		}
		hook, err := g.client.Webhook.CreateWebhook(ctx, req)
		if err != nil {
			errorChan <- err
			return
		}

		paras := map[string]string{}
		paras["Scope"] = hook.Scope.Type
		paras["SigningSecret"] = hook.SigningSecret
		paras["VerifyTLS"] = strconv.FormatBool(hook.VerifyTLS)
		webHook := &WebHook{
			ID:         hook.Id,
			Name:       hook.Name,
			Url:        hook.Url,
			Events:     hook.Events,
			Parameters: paras,
		}

		resultChan <- webHook
	}()

	return resultChan, errorChan
}

func (g *CircleClient) DeleteWebHookAsync(ctx context.Context, request DeleteWebHookRequest) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)
		_, err := g.client.Webhook.DeleteWebhook(ctx, request.WebHookId)
		errorChan <- err
	}()

	return errorChan
}

func (g *CircleClient) GetUserAsync(ctx context.Context) (<-chan *User, <-chan error) {
	resultChan := make(chan *User, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		circleUser, err := g.client.User.GetUserInfo(ctx)
		if err != nil {
			errorChan <- err
			return
		}

		user := &User{
			ID:   circleUser.ID,
			Name: circleUser.Name,
		}

		resultChan <- user
	}()

	return resultChan, errorChan
}

func (g *CircleClient) getCircleciVCS(vcsString string) (vcs string, vcs_full string) {
	vcsString = strings.ToLower(vcsString)
	if v, ok := CircleciVCSMapping[vcsString]; ok {
		vcs = vcsString
		vcs_full = v
	}

	if v, ok := CircleciBadgeVCSMapping[vcsString]; ok {
		vcs = v
		vcs_full = vcsString
	}

	return vcs, vcs_full
}

func (*CircleClient) buildProjectBadge(vcs string, user string, repo string, branch string) (badgeUrl string, badgeMarkDown string) {
	badgeUrl = fmt.Sprintf(CircleciProjectBadge, vcs, user, repo, branch)
	badgeMarkDown = fmt.Sprintf(CircleciProjectBadgeFull, vcs, user, repo, branch, vcs, user, repo, branch)
	return
}

func (*CircleClient) buildWorkflowBadge(vcs string, vcs_full string, org_name string, project_name string, workflow_name string, branch string) (badgeUrl string, badgeMarkDown string) {
	badgeUrl = fmt.Sprintf(CircleciWorkflowBadge,
		vcs,
		org_name,
		project_name,
		branch,
		workflow_name,
	)
	badgeMarkDown = fmt.Sprintf(CircleciWorkflowBadgeFull, vcs,
		org_name,
		project_name,
		branch,
		workflow_name,
		vcs_full,
		org_name,
		project_name,
		workflow_name,
		branch,
	)
	return
}

func (*CircleClient) buildProjectUrl(org_slug string, url_format string, vcs_full string, name string) string {
	url := ""
	orgSplit := strings.Split(org_slug, "/")
	if len(orgSplit) == 2 {
		url = fmt.Sprintf(url_format, vcs_full, orgSplit[1], name)
	}

	return url
}

func (g *CircleClient) buildProject(pro circleci.ProjectListItem, url string) Project {
	vcs, vcs_full := g.getCircleciVCS(pro.VcsType)
	badgeURL, badgeMarkdown := g.buildProjectBadge(vcs, pro.Username, pro.Reponame, pro.DefaultBranch)

	return Project{
		ID:            pro.Reponame,
		Name:          pro.Reponame,
		Url:           fmt.Sprintf(url, vcs_full, pro.Username, pro.Reponame),
		Properties:    map[string]string{"VCS_TYPE": vcs_full, "VCS_URL": pro.VcsURL},
		WorkflowRuns:  g.buildWrokflow(vcs, vcs_full, pro),
		Environments:  []string{},
		BadgeURL:      badgeURL,
		BadgeMarkDown: badgeMarkdown,
	}
}

func (g *CircleClient) buildWrokflow(vcs string, vcs_full string, pro circleci.ProjectListItem) map[string]WorkflowRun {
	default_branch := pro.DefaultBranch
	branchInfo := pro.Branches[default_branch]

	workflows := map[string]WorkflowRun{}
	for k, v := range branchInfo.LatestWorkflows {
		badgeURL, badgeMarkdown := g.buildWorkflowBadge(vcs, vcs_full, pro.Username, pro.Reponame, k, default_branch)
		workflows[k] = WorkflowRun{
			ID:            v.ID,
			Name:          k,
			Status:        v.Status,
			CreatedAt:     v.CreatedAt,
			BadgeURL:      badgeURL,
			BadgeMarkdown: badgeMarkdown,
		}
	}

	return workflows
}
