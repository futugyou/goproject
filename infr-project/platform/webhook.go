package platform

type Webhook struct {
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Activate bool              `json:"activate"`
	State    WebhookState      `json:"state"`
	Property map[string]string `json:"property"`
}

func NewWebhook(name string, url string, property map[string]string) *Webhook {
	return &Webhook{
		Name:     name,
		Url:      url,
		Activate: true,
		State:    WebhookInit,
		Property: property,
	}
}
