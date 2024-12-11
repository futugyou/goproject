package platform_provider

import (
	"context"
	"fmt"
	"strconv"

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
			Url:  repository.GetURL(),
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
			result = append(result, Project{
				ID:   repo.GetName(),
				Name: repo.GetName(),
				Url:  repo.GetURL(),
				Properties: map[string]string{
					"GITHUB_REPO":           repo.GetName(),
					"GITHUB_DETAULT_BRANCH": repo.GetDefaultBranch(),
				},
			})
		}
		resultChan <- result
	}()
	return resultChan, errorChan
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

		opts := &github.ListOptions{PerPage: 1000}
		githooks, _, err := g.client.Repositories.ListHooks(ctx, GITHUB_OWNER, filter.Name, opts)
		if err != nil {
			errorChan <- err
			return
		}
		hooks := []WebHook{}
		for _, hook := range githooks {
			if hook.GetActive() {
				paras := map[string]string{}
				githookconfig := hook.Config
				if githookconfig != nil {
					paras["ContentType"] = githookconfig.GetContentType()
					paras["InsecureSSL"] = githookconfig.GetInsecureSSL()
					paras["Secret"] = githookconfig.GetSecret()
					paras["URL"] = githookconfig.GetURL()
				}
				hooks = append(hooks, WebHook{
					ID:         fmt.Sprintf("%d", hook.GetID()),
					Name:       hook.GetName(),
					Url:        hook.GetURL(),
					Events:     hook.Events,
					Parameters: paras,
				})
			}
		}
		resultChan <- &Project{
			ID:   repository.GetName(),
			Name: repository.GetName(),
			Url:  repository.GetURL(),
			Properties: map[string]string{
				"GITHUB_REPO":           repository.GetName(),
				"GITHUB_DETAULT_BRANCH": repository.GetDefaultBranch(),
			},
			Hooks: hooks,
		}
	}()
	return resultChan, errorChan
}

// if need webhook secret, set it in WebHook.Parameters with key 'WEBHOOK_SECRET'
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

		if s, ok := request.WebHook.Parameters["WEBHOOK_SECRET"]; ok && len(s) > 0 {
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
			paras["Secret"] = githookconfig.GetSecret()
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
