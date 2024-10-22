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

type CreateWebhookRequest struct {
	Events           []string `json:"events,omitempty"`
	Url              string   `json:"url,omitempty"`
	ProjectIds       []string `json:"projectIds,omitempty"`
	BaseUrlParameter `json:"-"`
}

func (v *WebhookService) CreateWebhook(ctx context.Context, request CreateWebhookRequest) (*WebhookInfo, error) {
	u := &url.URL{
		Path: "/v1/webhooks",
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &WebhookInfo{}
	if err := v.client.http.Post(ctx, path, request, response); err != nil {
		return nil, err
	}
	return response, nil
}

type DeleteWebhookRequest struct {
	WebhookId        string
	BaseUrlParameter `json:"-"`
}

func (v *WebhookService) DeleteWebhook(ctx context.Context, request DeleteWebhookRequest) (*string, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/webhooks/%s", request.WebhookId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := ""
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

type GetWebhookParameter struct {
	WebhookId        string
	BaseUrlParameter `json:"-"`
}

func (v *WebhookService) GetWebhook(ctx context.Context, request GetWebhookParameter) (*WebhookInfo, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/v1/webhooks/%s", request.WebhookId),
	}
	params := request.GetUrlValues()
	u.RawQuery = params.Encode()
	path := u.String()

	response := &WebhookInfo{}
	if err := v.client.http.Delete(ctx, path, response); err != nil {
		return nil, err
	}
	return response, nil
}

type ListWebhookParameter struct {
	ProjectId        *string
	BaseUrlParameter `json:"-"`
}

func (v *WebhookService) ListWebhook(ctx context.Context, request ListWebhookParameter) ([]WebhookInfo, error) {
	u := &url.URL{
		Path: "/v1/webhooks",
	}
	params := request.GetUrlValues()
	if request.ProjectId != nil {
		params.Add("projectId", *request.ProjectId)
	}
	u.RawQuery = params.Encode()
	path := u.String()

	response := []WebhookInfo{}
	if err := v.client.http.Delete(ctx, path, &response); err != nil {
		return nil, err
	}
	return response, nil
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
