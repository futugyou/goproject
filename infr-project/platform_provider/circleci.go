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

		circleciProjects, err := client.client.Project.ListProject(ctx)
		if err != nil {
			errorChan <- err
			return
		}

		circleciProjectUrl := CircleciProjectUrl
		if url, ok := filter.Parameters["circleci_project_url"]; ok {
			circleciProjectUrl = url
		}

		projects := []Project{}
		for _, pro := range circleciProjects {
			// if pro.Followed {
			// ListProject API can only get followed projects, and the field is 'following' not 'followed'
			projects = append(projects, Project{
				ID:   pro.Reponame,
				Name: pro.Reponame,
				Url:  fmt.Sprintf(circleciProjectUrl, pro.VcsType, pro.Username, pro.Reponame),
				Properties: map[string]string{
					"VCS_TYPE": pro.VcsType,
					"VCS_URL":  pro.VcsURL,
				},
			})
			// }
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

		circleciProject, err := g.client.Project.GetProject(ctx, org_slug, filter.Name)
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

		circleciWebhooks, err := g.client.Webhook.ListWebhook(ctx, circleciProject.ID)
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
				Events:     hook.Events,
				Parameters: paras,
			})
		}

		project := &Project{
			ID:    circleciProject.ID,
			Name:  circleciProject.Name,
			Url:   url,
			Hooks: webHooks,
			Properties: map[string]string{
				"VCS_TYPE": circleciProject.VcsInfo.Provider,
				"VCS_URL":  circleciProject.VcsInfo.VcsURL,
			},
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
