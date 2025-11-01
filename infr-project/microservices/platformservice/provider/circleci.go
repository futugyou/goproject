package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/futugyou/circleci"
)

type circleClient struct {
	client *circleci.CircleciClient
	token  string
}

func newCircleClient(_ context.Context, token string) (*circleClient, error) {
	client := circleci.NewClientV2(token)
	return &circleClient{
		client,
		token,
	}, nil
}

func newCircleClientV1(token string) (*circleClient, error) {
	client := circleci.NewClientV1(token)
	return &circleClient{
		client,
		token,
	}, nil
}

const circleciProjectUrl = "https://app.circleci.com/pipelines/%s/%s/%s"
const circleciProjectBadge = "https://dl.circleci.com/status-badge/img/%s/%s/%s/tree/%s.svg?style=svg"
const circleciProjectBadgeFull = "[![CircleCI](https://dl.circleci.com/status-badge/img/%s/%s/%s/tree/%s.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/%s/%s/%s/tree/%s)"
const circleciWorkflowBadge = "https://dl.circleci.com/insights-snapshot/%s/%s/%s/%s/%s/badge.svg?window=30d"
const circleciWorkflowBadgeFull = "[![CircleCI](https://dl.circleci.com/insights-snapshot/%s/%s/%s/%s/%s/badge.svg?window=30d)](https://app.circleci.com/insights/%s/%s/%s/workflows/%s/overview?branch=%s&reporting-window=last-30-days&insights-snapshot=true)"

var circleciVCSMapping = map[string]string{"gh": "github", "bb": "bitbucket"}
var circleciBadgeVCSMapping = map[string]string{"github": "gh", "bitbucket": "bb"}

// Parameters MUST include org_slug. eg. gh/demo
// Parameters can set circleci_project_url. eg. https://app.circleci.com/pipelines/%s/%s/%s
func (g *circleClient) CreateProject(ctx context.Context, request CreateProjectRequest) (*Project, error) {
	org_slug := ""
	if org, ok := request.Parameters["org_slug"]; ok {
		org_slug = org
	} else {
		return nil, fmt.Errorf("create project request need 'org_slug' in parameters")
	}

	if _, err := g.client.Project.CreateProject(ctx, org_slug, request.Name); err != nil {
		return nil, err
	}

	project, err := g.client.Project.GetProject(ctx, org_slug, request.Name)
	if err != nil {
		return nil, err
	}

	orgSplit := strings.Split(org_slug, "/")
	url := ""
	if len(orgSplit) == 2 {
		_, vcs_full := g.getCircleciVCS(orgSplit[0])
		circleciUrl := circleciProjectUrl
		if url, ok := request.Parameters["circleci_project_url"]; ok {
			circleciUrl = url
		}
		url = fmt.Sprintf(circleciUrl, vcs_full, orgSplit[1], project.Name)
	}

	return &Project{
		ID:   project.ID,
		Name: project.Name,
		Url:  url,
	}, nil
}

// For now circleci api v2 do not include list project api,
// So, we need use api v1.1
func (g *circleClient) ListProject(ctx context.Context, filter ProjectFilter) ([]Project, error) {
	client, err := newCircleClientV1(g.token)
	if err != nil {
		return nil, err
	}

	circleciProjects, err := client.client.Project.ListProject(ctx)
	if err != nil {
		return nil, err
	}

	url := circleciProjectUrl
	if u, ok := filter.Parameters["circleci_project_url"]; ok {
		url = u
	}

	projects := []Project{}
	for _, pro := range circleciProjects {
		projects = append(projects, g.buildProject(pro, url))
	}

	return projects, nil
}

// Need org_slug in ProjectFilter 'Parameters'
// circleci webhook will set some other information in hook Parameters, it include 'Scope' 'SigningSecret' 'VerifyTLS'
func (g *circleClient) GetProject(ctx context.Context, filter ProjectFilter) (*Project, error) {
	org_slug := ""
	if org, ok := filter.Parameters["org_slug"]; ok {
		org_slug = org
	} else {
		return nil, fmt.Errorf("get project request need 'org_slug' in parameters")
	}

	circleciProject, err := g.client.Project.GetProject(ctx, org_slug, filter.Name)
	if err != nil {
		return nil, err
	}

	vcs, vcs_full := g.getCircleciVCS(circleciProject.VcsInfo.Provider)

	url_format := circleciProjectUrl
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

	return project, nil
}

// if need webhook secret, set it in WebHook.Parameters with key 'SigningSecret'
func (g *circleClient) CreateWebHook(ctx context.Context, request CreateWebHookRequest) (*WebHook, error) {
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
		return nil, err
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

	return webHook, nil
}

func (g *circleClient) DeleteWebHook(ctx context.Context, request DeleteWebHookRequest) error {
	_, err := g.client.Webhook.DeleteWebhook(ctx, request.WebHookId)
	return err
}

func (g *circleClient) GetUser(ctx context.Context) (*User, error) {
	circleUser, err := g.client.User.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:   circleUser.ID,
		Name: circleUser.Name,
	}

	return user, nil
}

func (g *circleClient) getCircleciVCS(vcsString string) (vcs string, vcs_full string) {
	vcsString = strings.ToLower(vcsString)
	if v, ok := circleciVCSMapping[vcsString]; ok {
		vcs = vcsString
		vcs_full = v
	}

	if v, ok := circleciBadgeVCSMapping[vcsString]; ok {
		vcs = v
		vcs_full = vcsString
	}

	return vcs, vcs_full
}

func (*circleClient) buildProjectBadge(vcs string, user string, repo string, branch string) (badgeUrl string, badgeMarkDown string) {
	badgeUrl = fmt.Sprintf(circleciProjectBadge, vcs, user, repo, branch)
	badgeMarkDown = fmt.Sprintf(circleciProjectBadgeFull, vcs, user, repo, branch, vcs, user, repo, branch)
	return
}

func (*circleClient) buildWorkflowBadge(vcs string, vcs_full string, org_name string, project_name string, workflow_name string, branch string) (badgeUrl string, badgeMarkDown string) {
	badgeUrl = fmt.Sprintf(circleciWorkflowBadge,
		vcs,
		org_name,
		project_name,
		branch,
		workflow_name,
	)
	badgeMarkDown = fmt.Sprintf(circleciWorkflowBadgeFull, vcs,
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

func (*circleClient) buildProjectUrl(org_slug string, url_format string, vcs_full string, name string) string {
	url := ""
	orgSplit := strings.Split(org_slug, "/")
	if len(orgSplit) == 2 {
		url = fmt.Sprintf(url_format, vcs_full, orgSplit[1], name)
	}

	return url
}

func (g *circleClient) buildProject(pro circleci.ProjectListItem, url string) Project {
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

func (g *circleClient) buildWrokflow(vcs string, vcs_full string, pro circleci.ProjectListItem) map[string]WorkflowRun {
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
