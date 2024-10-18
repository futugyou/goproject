package platform_provider

import (
	"context"
	"fmt"
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

// TODO: add all vcs mapping
var CircleciVCSMapping = map[string]string{"gh": "github"}

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

		if _, err := g.client.Project.CreateProject(org_slug, request.Name); err != nil {
			errorChan <- err
			return
		}

		project, err := g.client.Project.GetProject(org_slug, request.Name)
		if err != nil {
			errorChan <- err
			return
		}

		orgSplit := strings.Split(org_slug, "/")
		url := ""
		if len(orgSplit) == 2 {
			vsc := "github"
			if vcs_slug, ok := CircleciVCSMapping[orgSplit[0]]; ok {
				vsc = vcs_slug
			}
			circleciProjectUrl := CircleciProjectUrl
			if url, ok := request.Parameters["circleci_project_url"]; ok {
				circleciProjectUrl = url
			}
			url = fmt.Sprintf(circleciProjectUrl, vsc, orgSplit[1], project.Name)
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

		circleciProjects, err := client.client.Project.ListProject()
		if err != nil {
			errorChan <- err
			return
		}

		projects := []Project{}
		for _, pro := range circleciProjects {
			if pro.Followed {
				projects = append(projects, Project{
					ID:   pro.Reponame,
					Name: pro.Reponame,
					Url:  pro.VcsURL,
				})
			}
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
			errorChan <- fmt.Errorf("create project request need 'org_slug' in parameters")
			return
		}

		circleciProject, err := g.client.Project.GetProject(org_slug, filter.Name)
		if err != nil {
			errorChan <- err
			return
		}

		orgSplit := strings.Split(org_slug, "/")
		url := ""
		if len(orgSplit) == 2 {
			vsc := "github"
			if vcs_slug, ok := CircleciVCSMapping[orgSplit[0]]; ok {
				vsc = vcs_slug
			}
			circleciProjectUrl := CircleciProjectUrl
			if url, ok := filter.Parameters["circleci_project_url"]; ok {
				circleciProjectUrl = url
			}
			url = fmt.Sprintf(circleciProjectUrl, vsc, orgSplit[1], circleciProject.Name)
		}

		circleciWebhooks, err := g.client.Webhook.ListWebhook(circleciProject.ID)
		if err != nil {
			errorChan <- err
			return
		}
		webHooks := []WebHook{}
		for _, hook := range circleciWebhooks.Items {
			paras := map[string]string{}
			paras["Scope"] = hook.Scope.Type
			paras["SigningSecret"] = hook.SigningSecret
			paras["VerifyTLS"] = strconv.FormatBool(hook.VerifyTLS)
			webHooks = append(webHooks, WebHook{
				ID:         hook.Id,
				Name:       hook.Name,
				Url:        hook.Url,
				Parameters: paras,
			})
		}
		project := &Project{
			ID:    circleciProject.ID,
			Name:  circleciProject.Name,
			Url:   url,
			Hooks: webHooks,
		}

		resultChan <- project
	}()

	return resultChan, errorChan
}

func (g *CircleClient) CreateWebHookAsync(ctx context.Context, request CreateWebHookRequest) (<-chan *WebHook, <-chan error) {
	return nil, nil
}
