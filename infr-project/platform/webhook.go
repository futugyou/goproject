package platform

// WebhookState is the interface for webhook states.
type WebhookState interface {
	privateWebhookState() // Prevents external implementation
	String() string
}

// webhookState is the underlying implementation for WebhookState.
type webhookState string

// privateWebhookState makes webhookState implement WebhookState.
func (c webhookState) privateWebhookState() {}

// String makes webhookState implement WebhookState.
func (c webhookState) String() string {
	return string(c)
}

// Constants for the different webhook states.
const (
	WebhookInit     webhookState = "Init"
	WebhookCreating webhookState = "Creating"
	WebhookReady    webhookState = "Ready"
)

type Webhook struct {
	Name     string            `json:"name" validate:"required,min=3,max=50"`
	Url      string            `json:"url" validate:"required,min=3,max=50"`
	Activate bool              `json:"activate" validate:"required"`
	State    WebhookState      `json:"state" validate:"oneof=Init Creating Ready"`
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
