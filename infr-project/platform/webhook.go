package platform

type Webhook struct {
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Activate bool              `json:"activate"`
	State    WebhookState      `json:"state"`
	Property map[string]string `json:"property"`
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

func WithWebhookProperty(property map[string]string) WebhookOption {
	return func(w *Webhook) {
		w.Property = property
	}
}

func NewWebhook(name string, url string, opts ...WebhookOption) *Webhook {
	webhook := &Webhook{
		Name:     name,
		Activate: true,
		Url:      url,
		State:    WebhookInit,
		Property: make(map[string]string),
	}

	for _, opt := range opts {
		opt(webhook)
	}

	return webhook
}
