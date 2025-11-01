package domain

import (
	"maps"
	"fmt"
	"os"
)

func GetWebhookUrl(platform, project string) string {
	return fmt.Sprintf(os.Getenv("PROJECT_WEBHOOK_URL"), platform, project)
}

type Webhook struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Url        string              `json:"url"`
	Events     []string            `json:"events"`
	Activate   bool                `json:"activate"`
	State      WebhookState        `json:"state"`
	Properties map[string]Property `json:"properties"`
	Secrets    map[string]Secret   `json:"secrets"`
}

type WebhookOption func(*Webhook)

func WithWebhookActivate(activate bool) WebhookOption {
	return func(w *Webhook) {
		w.Activate = activate
	}
}

func WithWebhookState(state WebhookState) WebhookOption {
	return func(w *Webhook) {
		w.State = state
	}
}

func WithWebhookId(id string) WebhookOption {
	return func(w *Webhook) {
		w.ID = id
	}
}

func WithWebhookEvents(evens []string) WebhookOption {
	return func(w *Webhook) {
		w.Events = evens
	}
}

func WithWebhookProperties(properties map[string]Property) WebhookOption {
	return func(w *Webhook) {
		w.Properties = properties
	}
}

func WithWebhookSecrets(secrets map[string]Secret) WebhookOption {
	return func(w *Webhook) {
		w.Secrets = secrets
	}
}

func NewWebhook(name string, url string, opts ...WebhookOption) *Webhook {
	webhook := &Webhook{
		ID:         "",
		Name:       name,
		Url:        url,
		Events:     []string{},
		Activate:   true,
		State:      WebhookInit,
		Properties: make(map[string]Property),
		Secrets:    make(map[string]Secret),
	}

	for _, opt := range opts {
		opt(webhook)
	}

	return webhook
}

func (w *Webhook) UpdateProviderHookId(id string) *Webhook {
	w.ID = id
	return w
}

func (w *Webhook) UpdateProperties(properties map[string]Property) *Webhook {
	maps.Copy(w.Properties, properties)

	return w
}

func (w *Webhook) UpdateSecrets(secrets map[string]Secret) *Webhook {
	for k, v := range secrets {
		w.Secrets[k] = v
	}

	return w
}
