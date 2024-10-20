package vercel

import (
	"context"
	"fmt"
	"net/url"
)

type WebhookService service

var WebHookEvent []string = []string{
	"deployment.created",
	"deployment.succeeded",
	"deployment.ready",
	"deployment.promoted",
	"deployment.canceled",
	"deployment.error",
	"deployment.check-rerequested",
	"project.created",
	"project.removed",
	"integration-configuration.scope-change-confirmed",
	"integration-configuration.removed",
	"integration-configuration.permission-upgraded",
	"domain.created",
}

func (v *WebhookService) CreateWebhook(ctx context.Context, slug string, teamId string, req CreateWebhookRequest) (*WebhookInfo, error) {
	path := "/v1/webhooks"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &WebhookInfo{}
	err := v.client.http.Post(ctx, path, req, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *WebhookService) DeleteWebhook(ctx context.Context, id string, slug string, teamId string) (*string, error) {
	path := fmt.Sprintf("/v1/webhooks/%s", id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := ""
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (v *WebhookService) GetWebhook(ctx context.Context, id string, slug string, teamId string) (*WebhookInfo, error) {
	path := fmt.Sprintf("/v1/webhooks/%s", id)
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := &WebhookInfo{}
	err := v.client.http.Delete(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (v *WebhookService) ListWebhook(ctx context.Context, projectId string, slug string, teamId string) ([]WebhookInfo, error) {
	path := "/v1/webhooks"
	queryParams := url.Values{}
	if len(slug) > 0 {
		queryParams.Add("slug", slug)
	}
	if len(teamId) > 0 {
		queryParams.Add("teamId", teamId)
	}
	if len(projectId) > 0 {
		queryParams.Add("projectId", projectId)
	}
	if len(queryParams) > 0 {
		path += "?" + queryParams.Encode()
	}
	result := []WebhookInfo{}
	err := v.client.http.Delete(ctx, path, &result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CreateWebhookRequest struct {
	Events     []string `json:"events,omitempty"`
	Url        string   `json:"url,omitempty"`
	ProjectIds []string `json:"projectIds,omitempty"`
}

type WebhookInfo struct {
	CreatedAt        int          `json:"createdAt,omitempty"`
	Events           []string     `json:"events,omitempty"`
	Url              string       `json:"url,omitempty"`
	ProjectIds       []string     `json:"projectIds,omitempty"`
	Id               string       `json:"id,omitempty"`
	OwnerId          string       `json:"OwnerId,omitempty"`
	Secret           string       `json:"secret,omitempty"`
	UpdatedAt        int          `json:"updatedAt,omitempty"`
	ProjectsMetadata interface{}  `json:"projectsMetadata,omitempty"`
	Error            *VercelError `json:"error,omitempty"`
}
