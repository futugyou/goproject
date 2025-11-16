package provider

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/futugyou/extensions"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

type githubClient struct {
	client  *github.Client
	branch  string
	private bool
	owner   string
}

func newGithubClient(ctx context.Context, token string) (*githubClient, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &githubClient{
		client:  client,
		branch:  "master",
		private: true,
		owner:   "",
	}, nil
}

// Need GITHUB_BRANCH, GITHUB_PRIVATE, GITHUB_OWNER in CreateProjectRequest 'Parameters',
// default value 'master', 'true', ”.
// repository name use CreateProjectRequest 'Name'.
func (g *githubClient) CreateProject(ctx context.Context, request CreateProjectRequest) (*Project, error) {
	if value, ok := request.Parameters["GITHUB_BRANCH"]; ok && len(value) > 0 {
		g.branch = value
	}
	if privateString, ok := request.Parameters["GITHUB_PRIVATE"]; ok {
		if value, err := strconv.ParseBool(privateString); err != nil {
			g.private = value
		}
	}
	if value, ok := request.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
		g.owner = value
	}
	repo := &github.Repository{
		Name:          github.String(request.Name),
		DefaultBranch: github.String(g.branch),
		MasterBranch:  github.String(g.branch),
		Private:       github.Bool(g.private),
		AutoInit:      github.Bool(true),
	}

	repository, _, err := g.client.Repositories.Create(ctx, g.owner, repo)
	if err != nil {
		return nil, err
	}
	return &Project{
		ID:   fmt.Sprintf("%d", repository.GetID()),
		Name: repository.GetName(),
		Url:  repository.GetHTMLURL(),
	}, nil
}

func (g *githubClient) ListProject(ctx context.Context, filter ProjectFilter) ([]Project, error) {
	if value, ok := filter.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
		g.owner = value
	}

	var repos []*github.Repository
	var err error
	if len(g.owner) > 0 {
		opts := &github.RepositoryListByUserOptions{
			Type:      "owner",
			Sort:      "pushed",
			Direction: "desc",
			ListOptions: github.ListOptions{
				PerPage: 1000,
			},
		}
		repos, _, err = g.client.Repositories.ListByUser(ctx, g.owner, opts)
	} else {
		opts := &github.RepositoryListByAuthenticatedUserOptions{
			Type:      "all",
			Sort:      "pushed",
			Direction: "desc",
			ListOptions: github.ListOptions{
				PerPage: 1000,
			},
		}

		repos, _, err = g.client.Repositories.ListByAuthenticatedUser(ctx, opts)
	}

	if err != nil {
		return nil, err
	}

	result := []Project{}
	for _, repo := range repos {
		project := g.buildGithubProject(repo)
		result = append(result, *project)
	}
	return result, nil
}

func (g *githubClient) buildGithubProject(repo *github.Repository) *Project {
	badgeURL, badgeMarkdown := g.buildGithubProjectBadge(repo.GetArchived(), repo.GetHTMLURL())
	paras := map[string]string{}
	paras["GITHUB_REPO"] = repo.GetName()
	paras["GITHUB_DETAULT_BRANCH"] = repo.GetDefaultBranch()
	paras["ISSUES"] = fmt.Sprintf("%d", repo.GetOpenIssuesCount())
	paras["FORKS"] = fmt.Sprintf("%d", repo.GetForksCount())
	paras["WATCHS"] = fmt.Sprintf("%d", repo.GetStargazersCount())

	return &Project{
		ID:            fmt.Sprintf("%d", repo.GetID()),
		Name:          repo.GetName(),
		Url:           repo.GetHTMLURL(),
		Description:   repo.GetDescription(),
		Properties:    paras,
		BadgeURL:      badgeURL,
		BadgeMarkDown: badgeMarkdown,
		Tags:          repo.Topics,
		Environments:  []string{},
	}
}

// Need GITHUB_OWNER in ProjectFilter 'Parameters',
// default value ”.
// repository name use ProjectFilter 'Name'.
// github webhook config will set in hook Parameters, it include 'ContentType' 'InsecureSSL' 'Secret' 'URL'
func (g *githubClient) GetProject(ctx context.Context, filter ProjectFilter) (*Project, error) {
	if value, ok := filter.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
		g.owner = value
	}

	repository, _, err := g.client.Repositories.Get(ctx, g.owner, filter.Name)
	if err != nil {
		return nil, err
	}

	g.branch = repository.GetDefaultBranch()
	if g.branch == "" {
		g.branch = "master"
	}

	wfs := map[string]Workflow{}
	opts := &github.ListOptions{Page: 1, PerPage: 1000}
	if workflows, _, err := g.client.Actions.ListWorkflows(ctx, g.owner, filter.Name, opts); err != nil {
		log.Println(err.Error())
	} else {
		for _, v := range workflows.Workflows {
			workflow := Workflow{
				ID:        fmt.Sprintf("%d", v.GetID()),
				Name:      v.GetName(),
				Status:    v.GetState(),
				CreatedAt: v.GetCreatedAt().Format(time.RFC3339Nano),
				BadgeURL:  v.GetBadgeURL(),
			}

			paths := strings.Split(v.GetHTMLURL(), "/")
			path := ""
			if len(paths) > 0 {
				path = paths[len(paths)-1]
			}

			if len(path) > 0 {
				workflow.BadgeMarkdown = fmt.Sprintf("[![%s](%s)](https://github.com/%s/%s/actions/workflows/%s)",
					v.GetName(),
					v.GetBadgeURL(),
					g.owner,
					filter.Name,
					path,
				)
			}
			wfs[fmt.Sprintf("%d", v.GetID())] = workflow
		}
	}

	webHook := g.getWebHook(ctx, filter)

	envs := map[string]EnvironmentVariable{}
	if gitSecrets, _, err := g.client.Actions.ListRepoSecrets(ctx, g.owner, filter.Name, opts); err != nil {
		log.Println(err.Error())
	} else {
		for _, v := range gitSecrets.Secrets {
			if v == nil {
				continue
			}
			envs[v.Name] = EnvironmentVariable{
				ID:        v.Name,
				Key:       v.Name,
				CreatedAt: v.CreatedAt.Format(time.RFC3339Nano),
				UpdatedAt: v.UpdatedAt.Format(time.RFC3339Nano),
				Type:      v.Visibility,
				Value:     "",
			}
		}
	}

	deployments := map[string]Deployment{}
	// this api can not get deployment status, may be need graphql.
	if gitDeployments, _, err := g.client.Repositories.ListDeployments(ctx, g.owner, filter.Name, &github.DeploymentsListOptions{
		ListOptions: github.ListOptions{Page: 1, PerPage: 20},
	}); err != nil {
		log.Println(err.Error())
	} else {
		for _, v := range gitDeployments {
			badgeUrl, badgeMarkdown := g.buildGithubCommonBadge(v.GetTask(), "Unknown", v.GetStatusesURL())
			deployments[fmt.Sprintf("%d", v.GetID())] = Deployment{
				ID:            fmt.Sprintf("%d", v.GetID()),
				Name:          v.GetTask(),
				Environment:   v.GetEnvironment(),
				ReadyState:    "Unknown",
				ReadySubstate: "",
				CreatedAt:     v.GetCreatedAt().Format(time.RFC3339Nano),
				BadgeURL:      badgeUrl,
				BadgeMarkdown: badgeMarkdown,
				Description:   v.GetDescription(),
			}
		}
	}

	environments := []string{}
	if gitRuns, _, err := g.client.Repositories.ListEnvironments(ctx, g.owner, filter.Name, &github.EnvironmentListOptions{
		ListOptions: github.ListOptions{Page: 1, PerPage: 20},
	}); err != nil {
		log.Println(err.Error())
	} else {
		for _, v := range gitRuns.Environments {
			environments = append(environments, v.GetName())
		}
	}

	runs := map[string]WorkflowRun{}
	if gitRuns, _, err := g.client.Actions.ListRepositoryWorkflowRuns(ctx, g.owner, filter.Name, &github.ListWorkflowRunsOptions{
		Branch:      g.branch,
		ListOptions: github.ListOptions{Page: 1, PerPage: 20},
	}); err != nil {
		log.Println(err.Error())
	} else {
		for _, v := range gitRuns.WorkflowRuns {
			badgeUrl, badgeMarkdown := g.buildGithubCommonBadge(v.GetName(), v.GetStatus(), v.GetHTMLURL())
			runs[fmt.Sprintf("%d", v.GetID())] = WorkflowRun{
				ID:            fmt.Sprintf("%d", v.GetID()),
				Name:          v.GetName(),
				Description:   v.GetDisplayTitle(),
				Status:        v.GetStatus(),
				CreatedAt:     v.GetCreatedAt().Format(time.RFC3339Nano),
				BadgeURL:      badgeUrl,
				BadgeMarkdown: badgeMarkdown,
			}
		}
	}

	readme, _ := g.getGithubReadmeMarkdown(ctx, g.owner, filter.Name, g.branch)

	project := g.buildGithubProject(repository)
	project.WebHook = webHook
	project.EnvironmentVariables = envs
	project.Workflows = wfs
	project.WorkflowRuns = runs
	project.Deployments = deployments
	project.Environments = environments
	project.Readme = readme

	return project, nil
}

func (g *githubClient) GetSimpleProjectInfo(ctx context.Context, filter ProjectFilter) (*Project, error) {
	if value, ok := filter.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
		g.owner = value
	}

	repository, _, err := g.client.Repositories.Get(ctx, g.owner, filter.Name)
	if err != nil {
		return nil, err
	}

	return g.buildGithubProject(repository), nil
}

func (g *githubClient) getWebHook(ctx context.Context, filter ProjectFilter) *WebHook {
	if filter.WebHookID == nil {
		return nil
	}

	webHookID, err := strconv.ParseInt(*filter.WebHookID, 10, 64)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	hook, _, err := g.client.Repositories.GetHook(ctx, g.owner, filter.Name, webHookID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	paras := map[string]string{}
	paras["TestURL"] = hook.GetTestURL()
	paras["PingURL"] = hook.GetPingURL()
	var tls = "true"
	if hook.GetConfig().GetInsecureSSL() == "1" {
		tls = "false"
	}
	paras["ContentType"] = hook.GetConfig().GetContentType()
	paras["VerifyTLS"] = tls
	paras["SigningSecret"] = hook.GetConfig().GetSecret()
	paras["URL"] = hook.GetConfig().GetURL()

	return &WebHook{
		ID:         fmt.Sprintf("%d", hook.GetID()),
		Name:       hook.GetName(),
		Url:        hook.GetURL(),
		Events:     hook.Events,
		Activate:   hook.GetActive(),
		Parameters: paras,
	}
}

func (g *githubClient) getGithubReadmeMarkdown(ctx context.Context, owner, repo, branch string) (string, error) {
	readme, _, err := g.client.Repositories.GetReadme(ctx, owner, repo, &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		return "", fmt.Errorf("failed to get README: %w", err)
	}

	content, err := readme.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to decode README content: %w", err)
	}

	content = fixRelativeLinks(content, owner, repo, branch)

	return content, nil
}

func fixRelativeLinks(markdown, owner, repo, branch string) string {
	baseRawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/", owner, repo, branch)
	baseBlobURL := fmt.Sprintf("https://github.com/%s/%s/blob/%s/", owner, repo, branch)

	imgRe := regexp.MustCompile(`!\[([^\]]*)]\(([^)]+)\)`)
	markdown = imgRe.ReplaceAllStringFunc(markdown, func(match string) string {
		sub := imgRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		url := sub[2]
		if regexp.MustCompile(`(?i)^https?://`).MatchString(url) {
			return match
		}
		return fmt.Sprintf("![%s](%s%s)", sub[1], baseRawURL, url)
	})

	linkRe := regexp.MustCompile(`\[([^\]]+)]\(([^)]+)\)`)
	markdown = linkRe.ReplaceAllStringFunc(markdown, func(match string) string {
		sub := linkRe.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		url := sub[2]
		if regexp.MustCompile(`(?i)^https?://`).MatchString(url) {
			return match
		}
		return fmt.Sprintf("[%s](%s%s)", sub[1], baseBlobURL, url)
	})

	return markdown
}

func (g *githubClient) buildGithubCommonBadge(name string, status string, url string) (badgeUrl string, badgeMarkDown string) {
	name = extensions.Sanitize2String(name, " ")
	color := "red"
	switch status {
	case "in_progress", "queued", "neutral", "skipped", "waiting", "pending", "action_required", "Unknown":
		color = "yellow"
	case "completed", "success":
		color = "brightgreen"
	}

	status = strings.ReplaceAll(status, "-", "%20")
	badgeUrl = fmt.Sprintf(CommonProjectBadge, name, status, color, "github", url)
	badgeMarkDown = fmt.Sprintf("![%s](%s)", "status", badgeUrl)
	return
}

func (g *githubClient) buildGithubProjectBadge(archived bool, url string) (badgeUrl string, badgeMarkDown string) {
	badgeUrl = fmt.Sprintf(CommonProjectBadge, "status", "Unarchived", "brightgreen", "github", url)
	if archived {
		badgeUrl = fmt.Sprintf(CommonProjectBadge, "status", "Archived", "red", "github", url)
	}

	badgeMarkDown = fmt.Sprintf("![%s](%s)", "status", badgeUrl)
	return
}

// if need webhook secret, set it in WebHook.Parameters with key 'SigningSecret'
func (g *githubClient) CreateWebHook(ctx context.Context, request CreateWebHookRequest) (*WebHook, error) {
	config := &github.HookConfig{
		ContentType: github.String("json"),
		InsecureSSL: github.String("1"),
		URL:         github.String(request.Url),
	}

	if len(request.SigningSecret) > 0 {
		config.Secret = github.String(request.SigningSecret)
	}

	events := request.Events
	if len(events) == 0 {
		events = []string{"push", "pull_request"}
	} else {
		events = Intersect(events, []string{"push", "pull_request"})
	}
	hookParam := &github.Hook{
		Name:   github.String(request.Name),
		Config: config,
		Events: events,
	}

	githook, _, err := g.client.Repositories.CreateHook(ctx, request.PlatformID, request.ProjectID, hookParam)
	if err != nil {
		return nil, err
	}

	paras := map[string]string{}
	githookconfig := githook.Config
	if githookconfig != nil {
		paras["ContentType"] = githookconfig.GetContentType()
		paras["InsecureSSL"] = githookconfig.GetInsecureSSL()
		paras["URL"] = githookconfig.GetURL()
	}
	hook := &WebHook{
		ID:            fmt.Sprintf("%d", githook.GetID()),
		Name:          githook.GetName(),
		Url:           githook.GetURL(),
		Events:        githook.Events,
		Parameters:    paras,
		SigningSecret: githookconfig.GetSecret(),
	}

	return hook, nil
}

func (g *githubClient) GetWebHookByUrl(ctx context.Context, request GetWebHookRequest) (*WebHook, error) {
	if value, ok := request.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
		g.owner = value
	}
	opts := &github.ListOptions{Page: 1, PerPage: 20}
	hooks, _, err := g.client.Repositories.ListHooks(ctx, g.owner, request.ProjectID, opts)
	if err != nil {
		return nil, err
	}

	for _, githook := range hooks {
		if githook.GetURL() == request.Url {
			paras := map[string]string{}
			githookconfig := githook.Config
			if githookconfig != nil {
				paras["ContentType"] = githookconfig.GetContentType()
				paras["InsecureSSL"] = githookconfig.GetInsecureSSL()
				paras["SigningSecret"] = githookconfig.GetSecret()
				paras["URL"] = githookconfig.GetURL()
			}
			return &WebHook{
				ID:         fmt.Sprintf("%d", githook.GetID()),
				Name:       githook.GetName(),
				Url:        githook.GetURL(),
				Events:     githook.Events,
				Parameters: paras,
			}, nil
		}
	}

	return nil, fmt.Errorf("webhook not found")
}

// Need GITHUB_OWNER and GITHUB_REPO in DeleteWebHookRequest 'Parameters'
func (g *githubClient) DeleteWebHook(ctx context.Context, request DeleteWebHookRequest) error {
	webHookId, err := strconv.ParseInt(request.WebHookID, 10, 64)
	if err != nil {
		return err
	}

	owner := ""
	if value, ok := request.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
		owner = value
	} else {
		return fmt.Errorf("github DeleteHook need GITHUB_OWNER")
	}

	repo := ""
	if value, ok := request.Parameters["GITHUB_REPO"]; ok && len(value) > 0 {
		repo = value
	} else {
		return fmt.Errorf("github DeleteHook need GITHUB_REPO")
	}

	_, err = g.client.Repositories.DeleteHook(ctx, owner, repo, webHookId)

	return err
}

func (g *githubClient) GetUser(ctx context.Context) (*User, error) {
	githubUser, _, err := g.client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:   fmt.Sprintf("%d", githubUser.GetID()),
		Name: githubUser.GetName(),
	}

	return user, nil
}
