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

// Need GITHUB_BRANCH, GITHUB_PRIVATE, GITHUB_OWNER in request 'Parameters',
// default value 'master', 'true', ”.
// repository name use request 'Name'.
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
			Name: *repository.Name,
			Url:  *repository.URL,
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
				Name: *repo.Name,
				Url:  *repo.URL,
			})
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (g *GithubClient) GetProjectAsync(ctx context.Context, filter ProjectFilter) (<-chan *Project, <-chan error) {
	return nil, nil
}

func (g *GithubClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}
