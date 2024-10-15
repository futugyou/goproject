package platform_provider

import (
	"context"
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

		opts := &github.RepositoryListByAuthenticatedUserOptions{
			Sort: "pushed",
			ListOptions: github.ListOptions{
				PerPage: 1000,
			},
		}

		repos, _, err := g.client.Repositories.ListByAuthenticatedUser(ctx, opts)
		if err != nil {
			errorChan <- err
			return
		}
		result := []Project{}
		for _, repo := range repos {
			result = append(result, Project{
				Name: repo.GetName(),
				Url:  repo.GetURL(),
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
					Name:       hook.GetName(),
					Url:        hook.GetURL(),
					Parameters: paras,
				})
			}
		}
		resultChan <- &Project{
			Name:  repository.GetName(),
			Url:   repository.GetURL(),
			Hooks: hooks,
		}
	}()
	return resultChan, errorChan
}

func (g *GithubClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}
