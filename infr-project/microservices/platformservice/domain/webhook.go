package domain

import (
	"fmt"
)

func GetWebhookUrl(platform, project, url string) string {
	return fmt.Sprintf(url, platform, project)
}

// objectvaule, `ID` is just a regular field, the same as `Url`.
type Webhook struct {
	ID            string       `json:"id"`
	Url           string       `json:"url"`
	Events        []string     `json:"events"`
	State         WebhookState `json:"state"`
	SigningSecret string       `json:"signing_secret"`
}

type WebhookOption func(*Webhook)

func WithWebhookState(state WebhookState) WebhookOption {
	return func(w *Webhook) {
		w.State = state
	}
}

func WithWebhookID(id string) WebhookOption {
	return func(w *Webhook) {
		w.ID = id
	}
}

func WithWebhookEvents(evens []string) WebhookOption {
	return func(w *Webhook) {
		w.Events = evens
	}
}

func WithWebhookSigningSecret(signingSecret string) WebhookOption {
	return func(w *Webhook) {
		w.SigningSecret = signingSecret
	}
}

func NewWebhook(url string, opts ...WebhookOption) *Webhook {
	webhook := &Webhook{
		Url:    url,
		Events: []string{},
		State:  WebhookCreating,
	}

	for _, opt := range opts {
		opt(webhook)
	}

	return webhook
}
