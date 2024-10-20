package circleci

import (
	"context"
	"fmt"
)

type WebhookService service

var WebhookEvents []string = []string{
	"workflow-completed",
	"job-completed",
}

type CreateWebhookRequest struct {
	Name          string
	Events        []string
	Url           string
	VerifyTLS     bool
	SigningSecret string
	ScopeId       string
	ScopeType     string
}

func (s *WebhookService) CreateWebhook(ctx context.Context, request CreateWebhookRequest) (*Webhook, error) {
	path := "/webhook"
	req := Webhook{
		Name:          request.Name,
		Events:        request.Events,
		Url:           request.Url,
		VerifyTLS:     request.VerifyTLS,
		SigningSecret: request.SigningSecret,
		Scope: WebhookScope{
			Id:   request.ScopeId,
			Type: request.ScopeType,
		},
	}

	result := &Webhook{}
	if err := s.client.http.Post(ctx, path, req, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WebhookService) ListWebhook(ctx context.Context, projectId string) (*ListWebhookResponse, error) {
	path := fmt.Sprintf("/webhook?scope-id=%s&scope-type=project", projectId)
	result := &ListWebhookResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WebhookService) GetWebhook(ctx context.Context, webhookId string) (*Webhook, error) {
	path := fmt.Sprintf("/webhook/%s", webhookId)
	result := &Webhook{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

type UpdateWebhookRequest struct {
	WebhookId     string
	Name          string
	Events        []string
	Url           string
	VerifyTLS     bool
	SigningSecret string
}

func (s *WebhookService) UpdateWebhook(ctx context.Context, request UpdateWebhookRequest) (*Webhook, error) {
	path := fmt.Sprintf("/webhook/%s", request.WebhookId)
	req := Webhook{
		Name:          request.Name,
		Events:        request.Events,
		Url:           request.Url,
		VerifyTLS:     request.VerifyTLS,
		SigningSecret: request.SigningSecret,
	}

	result := &Webhook{}
	if err := s.client.http.Put(ctx, path, req, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *WebhookService) DeleteWebhook(ctx context.Context, webhookId string) (*BaseResponse, error) {
	path := fmt.Sprintf("/webhook/%s", webhookId)
	result := &BaseResponse{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

type ListWebhookResponse struct {
	Items         []Webhook `json:"items"`
	NextPageToken string    `json:"next_page_token"`
	Message       *string   `json:"message,omitempty"`
}

type Webhook struct {
	Name          string       `json:"name,omitempty"`
	Events        []string     `json:"events,omitempty"` // "workflow-completed" "job-completed"
	Url           string       `json:"url,omitempty"`
	VerifyTLS     bool         `json:"verify-tls,omitempty"`
	SigningSecret string       `json:"signing-secret,omitempty"`
	Scope         WebhookScope `json:"scope,omitempty"`
	UpdatedAt     string       `json:"updated-at,omitempty"`
	Id            string       `json:"id,omitempty"`
	CreatedAt     string       `json:"created-at,omitempty"`
	Message       *string      `json:"message,omitempty"`
}

type WebhookScope struct {
	Id   string `json:"id"`
	Type string `json:"type"` //"project"
}
