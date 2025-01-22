package platform_provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/futugyou/extensions"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

var GITHUB_BRANCH = "master"
var GITHUB_PRIVATE = true
var GITHUB_OWNER = ""

type GithubClient struct {
	client *github.Client
}

func NewGithubClient(token string) (*GithubClient, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &GithubClient{
		client,
	}, nil
}

// Need GITHUB_BRANCH, GITHUB_PRIVATE, GITHUB_OWNER in CreateProjectRequest 'Parameters',
// default value 'master', 'true', ”.
// repository name use CreateProjectRequest 'Name'.
func (g *GithubClient) CreateProjectAsync(ctx context.Context, request CreateProjectRequest) (<-chan *Project, <-chan error) {
	resultChan := make(chan *Project, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)

		if value, ok := request.Parameters["GITHUB_BRANCH"]; ok && len(value) > 0 {
			GITHUB_BRANCH = value
		}
		if privateString, ok := request.Parameters["GITHUB_PRIVATE"]; ok {
			if value, err := strconv.ParseBool(privateString); err != nil {
				GITHUB_PRIVATE = value
			}
		}
		if value, ok := request.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
			GITHUB_OWNER = value
		}
		repo := &github.Repository{
			Name:          github.String(request.Name),
			DefaultBranch: github.String(GITHUB_BRANCH),
			MasterBranch:  github.String(GITHUB_BRANCH),
			Private:       github.Bool(GITHUB_PRIVATE),
		}

		repository, _, err := g.client.Repositories.Create(ctx, GITHUB_OWNER, repo)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- &Project{
			ID:   repository.GetName(), // gitHub repository uses name more often than id
			Name: repository.GetName(),
			Url:  repository.GetHTMLURL(),
		}
	}()
	return resultChan, errorChan
}

func (g *GithubClient) ListProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan []Project, <-chan error) {
	resultChan := make(chan []Project, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)

		if value, ok := filter.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
			GITHUB_OWNER = value
		}

		var repos []*github.Repository
		var err error
		if len(GITHUB_OWNER) > 0 {
			opts := &github.RepositoryListByUserOptions{
				Type:      "owner",
				Sort:      "pushed",
				Direction: "desc",
				ListOptions: github.ListOptions{
					PerPage: 1000,
				},
			}
			repos, _, err = g.client.Repositories.ListByUser(ctx, GITHUB_OWNER, opts)
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
			errorChan <- err
			return
		}

		result := []Project{}
		for _, repo := range repos {
			project := g.buildGithubProject(repo)
			result = append(result, project)
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (g *GithubClient) buildGithubProject(repo *github.Repository) Project {
	badgeURL, badgeMarkdown := g.buildGithubProjectBadge(repo.GetArchived(), repo.GetHTMLURL())
	paras := map[string]string{}
	paras["GITHUB_REPO"] = repo.GetName()
	paras["GITHUB_DETAULT_BRANCH"] = repo.GetDefaultBranch()
	paras["ISSUES"] = fmt.Sprintf("%d", repo.GetOpenIssuesCount())
	paras["FORKS"] = fmt.Sprintf("%d", repo.GetForksCount())
	paras["WATCHS"] = fmt.Sprintf("%d", repo.GetStargazersCount())
	return Project{
		ID:            repo.GetName(),
		Name:          repo.GetName(),
		Url:           repo.GetHTMLURL(),
		Description:   repo.GetDescription(),
		Properties:    paras,
		BadgeURL:      badgeURL,
		BadgeMarkDown: badgeMarkdown,
		Environments:  []string{},
	}
}

// Need GITHUB_OWNER in ProjectFilter 'Parameters',
// default value ”.
// repository name use ProjectFilter 'Name'.
// github webhook config will set in hook Parameters, it include 'ContentType' 'InsecureSSL' 'Secret' 'URL'
func (g *GithubClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	resultChan := make(chan *Project, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)

		if value, ok := filter.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
			GITHUB_OWNER = value
		}

		repository, _, err := g.client.Repositories.Get(ctx, GITHUB_OWNER, filter.Name)
		if err != nil {
			errorChan <- err
			return
		}

		wfs := map[string]Workflow{}
		opts := &github.ListOptions{Page: 1, PerPage: 1000}
		if workflows, _, err := g.client.Actions.ListWorkflows(ctx, GITHUB_OWNER, filter.Name, opts); err != nil {
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
						GITHUB_OWNER,
						filter.Name,
						path,
					)
				}
				wfs[fmt.Sprintf("%d", v.GetID())] = workflow
			}
		}

		hooks := []WebHook{}
		if githooks, _, err := g.client.Repositories.ListHooks(ctx, GITHUB_OWNER, filter.Name, opts); err != nil {
			log.Println(err.Error())
		} else {
			for _, hook := range githooks {
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

				hooks = append(hooks, WebHook{
					ID:         fmt.Sprintf("%d", hook.GetID()),
					Name:       hook.GetName(),
					Url:        hook.GetURL(),
					Events:     hook.Events,
					Activate:   hook.GetActive(),
					Parameters: paras,
				})
			}
		}

		envs := map[string]EnvironmentVariable{}
		if gitSecrets, _, err := g.client.Actions.ListRepoSecrets(ctx, GITHUB_OWNER, filter.Name, opts); err != nil {
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
		if gitDeployments, _, err := g.client.Repositories.ListDeployments(ctx, GITHUB_OWNER, filter.Name, &github.DeploymentsListOptions{
			ListOptions: github.ListOptions{Page: 1, PerPage: 20},
		}); err != nil {
			log.Println(err.Error())
		} else {
			for _, v := range gitDeployments {
				deployments[fmt.Sprintf("%d", v.GetID())] = Deployment{
					ID:          fmt.Sprintf("%d", v.GetID()),
					Name:        v.GetTask(),
					Environment: v.GetEnvironment(),
					CreatedAt:   v.GetCreatedAt().Format(time.RFC3339Nano),
				}
			}
		}

		environments := []string{}
		if gitRuns, _, err := g.client.Repositories.ListEnvironments(ctx, GITHUB_OWNER, filter.Name, &github.EnvironmentListOptions{
			ListOptions: github.ListOptions{Page: 1, PerPage: 20},
		}); err != nil {
			log.Println(err.Error())
		} else {
			for _, v := range gitRuns.Environments {
				environments = append(environments, v.GetEnvironmentName())
			}
		}

		runs := map[string]WorkflowRun{}
		if gitRuns, _, err := g.client.Actions.ListRepositoryWorkflowRuns(ctx, GITHUB_OWNER, filter.Name, &github.ListWorkflowRunsOptions{
			Branch:      repository.GetDefaultBranch(),
			ListOptions: github.ListOptions{Page: 1, PerPage: 20},
		}); err != nil {
			log.Println(err.Error())
		} else {
			for _, v := range gitRuns.WorkflowRuns {
				badgeUrl, badgeMarkdown := g.buildGithubWorkflowBadge(v.GetName(), v.GetStatus(), v.GetHTMLURL())
				runs[fmt.Sprintf("%d", v.GetID())] = WorkflowRun{
					ID:            fmt.Sprintf("%d", v.GetID()),
					Name:          v.GetName(),
					Status:        v.GetStatus(),
					CreatedAt:     v.GetCreatedAt().Format(time.RFC3339Nano),
					BadgeURL:      badgeUrl,
					BadgeMarkdown: badgeMarkdown,
				}
			}
		}

		project := g.buildGithubProject(repository)
		project.WebHooks = hooks
		project.EnvironmentVariables = envs
		project.Workflows = wfs
		project.WorkflowRuns = runs
		project.Deployments = deployments
		project.Environments = environments
		resultChan <- &project
	}()
	return resultChan, errorChan
}

func (g *GithubClient) buildGithubWorkflowBadge(name string, status string, url string) (badgeUrl string, badgeMarkDown string) {
	name = extensions.Sanitize2String(name, " ")
	color := "red"
	switch status {
	case "in_progress", "queued", "neutral", "skipped", "waiting", "pending", "action_required":
		color = "yellow"
	case "completed", "success":
		color = "brightgreen"
	}

	status = strings.ReplaceAll(status, "-", "%20")
	badgeUrl = fmt.Sprintf(CommonProjectBadge, name, status, color, "github", url)
	badgeMarkDown = fmt.Sprintf("![%s](%s)", "status", badgeUrl)
	return
}

func (g *GithubClient) buildGithubProjectBadge(archived bool, url string) (badgeUrl string, badgeMarkDown string) {
	badgeUrl = fmt.Sprintf(CommonProjectBadge, "status", "Unarchived", "brightgreen", "github", url)
	if archived {
		badgeUrl = fmt.Sprintf(CommonProjectBadge, "status", "Archived", "red", "github", url)
	}

	badgeMarkDown = fmt.Sprintf("![%s](%s)", "status", badgeUrl)
	return
}

// if need webhook secret, set it in WebHook.Parameters with key 'SigningSecret'
func (g *GithubClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	resultChan := make(chan *WebHook, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		config := &github.HookConfig{
			ContentType: github.String("json"),
			InsecureSSL: github.String("1"),
			URL:         github.String(request.WebHook.Url),
		}

		if s, ok := request.WebHook.Parameters["SigningSecret"]; ok && len(s) > 0 {
			config.Secret = github.String(s)
		}

		events := request.WebHook.Events
		if len(events) == 0 {
			events = []string{"push", "pull_request"}
		} else {
			events = Intersect(events, []string{"push", "pull_request"})
		}
		hookParam := &github.Hook{
			Name:   github.String(request.WebHook.Name),
			Config: config,
			Events: events,
		}

		githook, _, err := g.client.Repositories.CreateHook(ctx, request.PlatformId, request.ProjectId, hookParam)
		if err != nil {
			errorChan <- err
			return
		}

		paras := map[string]string{}
		githookconfig := githook.Config
		if githookconfig != nil {
			paras["ContentType"] = githookconfig.GetContentType()
			paras["InsecureSSL"] = githookconfig.GetInsecureSSL()
			paras["SigningSecret"] = githookconfig.GetSecret()
			paras["URL"] = githookconfig.GetURL()
		}
		hook := &WebHook{
			ID:         fmt.Sprintf("%d", githook.GetID()),
			Name:       githook.GetName(),
			Url:        githook.GetURL(),
			Events:     githook.Events,
			Parameters: paras,
		}
		resultChan <- hook
	}()
	return resultChan, errorChan
}

// Need GITHUB_OWNER in DeleteWebHookRequest 'Parameters'
// Need GITHUB_REPO in DeleteWebHookRequest 'Parameters',
func (g *GithubClient) DeleteWebHookAsync(ctx context.Context, request DeleteWebHookRequest) <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		defer close(errorChan)
		webHookId, err := strconv.ParseInt(request.WebHookId, 10, 64)
		if err != nil {
			errorChan <- err
			return
		}

		owner := ""
		if value, ok := request.Parameters["GITHUB_OWNER"]; ok && len(value) > 0 {
			owner = value
		} else {
			errorChan <- fmt.Errorf("github DeleteHook need GITHUB_OWNER")
			return
		}

		repo := ""
		if value, ok := request.Parameters["GITHUB_REPO"]; ok && len(value) > 0 {
			repo = value
		} else {
			errorChan <- fmt.Errorf("github DeleteHook need GITHUB_REPO")
			return
		}

		_, err = g.client.Repositories.DeleteHook(ctx, owner, repo, webHookId)
		errorChan <- err
	}()

	return errorChan
}

func (g *GithubClient) GetUserAsync(ctx context.Context) (<-chan *User, <-chan error) {
	resultChan := make(chan *User, 1)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		githubUser, _, err := g.client.Users.Get(ctx, "")
		if err != nil {
			errorChan <- err
			return
		}

		user := &User{
			ID:   fmt.Sprintf("%d", githubUser.GetID()),
			Name: githubUser.GetName(),
		}

		resultChan <- user
	}()

	return resultChan, errorChan
}
