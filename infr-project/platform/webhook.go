package platform

type Webhook struct {
	Name       string            `json:"name"`
	Url        string            `json:"url"`
	Activate   bool              `json:"activate"`
	State      WebhookState      `json:"state"`
	Properties map[string]string `json:"properties"`
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

func WithWebhookProperty(properties map[string]string) WebhookOption {
	return func(w *Webhook) {
		w.Properties = properties
	}
}

func NewWebhook(name string, url string, opts ...WebhookOption) *Webhook {
	webhook := &Webhook{
		Name:       name,
		Activate:   true,
		Url:        url,
		State:      WebhookInit,
		Properties: make(map[string]string),
	}

	for _, opt := range opts {
		opt(webhook)
	}

	return webhook
}
