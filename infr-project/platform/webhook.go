package platform

type Webhook struct {
	Name           string              `json:"name"`
	Url            string              `json:"url"`
	Activate       bool                `json:"activate"`
	State          WebhookState        `json:"state"`
	Properties     map[string]Property `json:"properties"`
	Secrets        map[string]Secret   `json:"secrets"`
	ProviderHookId string              `json:"provider_hook_id"`
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
		w.ProviderHookId = id
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
		Name:       name,
		Activate:   true,
		Url:        url,
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
	w.ProviderHookId = id
	return w
}
